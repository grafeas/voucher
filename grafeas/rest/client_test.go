package rest

import (
	"context"
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/golang/mock/gomock"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/attestation"
	vgrafeas "github.com/Shopify/voucher/grafeas"
	mocks "github.com/Shopify/voucher/grafeas/rest/mocks"
	"github.com/Shopify/voucher/grafeas/rest/objects"
	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/signer"
	"github.com/Shopify/voucher/signer/kms"
	"github.com/Shopify/voucher/signer/pgp"
)

var basePath string

func TestMain(m *testing.M) {
	flag.StringVar(&basePath, "grafeas", "", "the base path to the grafeas instance to use for testing")
	flag.Parse()
	os.Exit(m.Run())
}

func getCanonicalRef(t *testing.T) reference.Canonical {
	named, err := reference.ParseNamed("us.gcr.io/grafeas/grafeas-server@sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733")
	require.NoError(t, err, "named")
	canonicalRef, err := reference.WithDigest(named, digest.FromString("sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733"))
	require.NoError(t, err, "canonicalRef")
	return canonicalRef
}

func TestCanAttest(t *testing.T) {
	project := "project"
	keyringKms, _ := kms.NewSigner(make(map[string]kms.Key))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grafeas := mocks.NewMockGrafeasAPIService(ctrl)
	tcs := map[string]struct {
		keyring        signer.AttestationSigner
		expectedResult bool
	}{
		"can attest pgp": {
			keyring:        pgp.NewKeyRing(),
			expectedResult: true,
		},
		"can attest kms": {
			keyring:        keyringKms,
			expectedResult: true,
		},
		"cannot attest": {
			expectedResult: false,
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			client, err := NewClient(context.Background(), project, project, test.keyring, grafeas)
			canAttest := client.CanAttest()
			assert.Equal(t, test.expectedResult, canAttest)
			require.NoError(t, err)
		})
	}
}

func TestNewPayloadBody(t *testing.T) {
	project := "project"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grafeas := mocks.NewMockGrafeasAPIService(ctrl)
	client, _ := NewClient(context.Background(), project, project, pgp.NewKeyRing(), grafeas)
	ref := getCanonicalRef(t)
	tcs := map[string]struct {
		reference       reference.Canonical
		expectedPayload attestation.Payload
	}{
		"valid reference": {
			reference: ref,
			expectedPayload: attestation.Payload{
				Critical: attestation.PayloadCritical{
					Identity: attestation.PayloadIdentity{
						DockerReference: ref.Name(),
					},
					Image: attestation.PayloadImage{
						DockerManifestDigest: ref.Digest(),
					},
					Type: "Google cloud binauthz container signature",
				},
			},
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			payload, err := client.NewPayloadBody(test.reference)
			expectedPayloadStr, _ := test.expectedPayload.ToString()
			assert.Equal(t, expectedPayloadStr, payload)
			require.NoError(t, err)
		})
	}
}

func TestAddAttestationToImage(t *testing.T) {
	ctx := context.Background()
	project := "project"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grafeas := mocks.NewMockGrafeasAPIService(ctrl)
	keyringKms, _ := kms.NewSigner(map[string]kms.Key{
		"key": {
			Algo: kms.AlgoSHA512,
		},
	})
	tcs := map[string]struct {
		reference     reference.Canonical
		keyring       signer.AttestationSigner
		payload       voucher.Attestation
		expectedError error
	}{
		"cannot attest": {
			expectedError: errCannotAttest,
		},
		"sign error": {
			reference:     getCanonicalRef(t),
			keyring:       keyringKms,
			payload:       voucher.Attestation{},
			expectedError: errors.New("no signing entity exists for check"),
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			client, _ := NewClient(context.Background(), project, project, test.keyring, grafeas)
			_, err := client.AddAttestationToImage(ctx, test.reference, test.payload)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetAttestations(t *testing.T) {
	ctx := context.Background()
	project := "project"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grafeasMock := mocks.NewMockGrafeasAPIService(ctrl)
	client, _ := NewClient(context.Background(), project, project, pgp.NewKeyRing(), grafeasMock)
	ref := getCanonicalRef(t)
	occs := createAllOccurrences()
	tcs := map[string]struct {
		returnOccs     objects.ListOccurrencesResponse
		expectedResult []voucher.SignedAttestation
		expectedError  error
	}{
		"valid input": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: occs,
			},
			expectedResult: []voucher.SignedAttestation{{
				Attestation: voucher.Attestation{
					CheckName: "notename",
					Body:      string(objects.AttestationUnspecified),
				}},
			},
		},
		"no data": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: []objects.Occurrence{},
			},
			expectedError: &voucher.NoMetadataError{
				Type: voucher.AttestationType,
				Err:  vgrafeas.ErrNoOccurrences,
			},
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			grafeasMock.EXPECT().ListOccurrences(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.returnOccs, nil)
			attestations, err := client.GetAttestations(ctx, ref)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, attestations)
		})
	}
}

func TestGetVulnerabilities(t *testing.T) {
	ctx := context.Background()
	project := "project"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ref := getCanonicalRef(t)
	occs := createAllOccurrences()
	activeContAnalysis := objects.DiscoveredContinuousAnalysisActive
	successStatus := objects.DiscoveredAnalysisStatusFinishedSuccess
	noteKindD := objects.NoteKindDiscovery
	tcs := map[string]struct {
		returnOccs     objects.ListOccurrencesResponse
		expectedResult []voucher.Vulnerability
		expectedError  error
	}{
		"valid input": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: occs,
			},
			expectedResult: []voucher.Vulnerability{{
				Name:     "notename",
				Severity: voucher.NegligibleSeverity,
			}},
		},
		"no data": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: []objects.Occurrence{{Name: "name4",
					Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
					NoteName: "notename_invalid", Kind: &noteKindD,
					Discovered: &objects.DiscoveryDetails{
						Discovered: &objects.DiscoveryDiscovered{ContinuousAnalysis: &activeContAnalysis,
							AnalysisStatus: &successStatus}}}},
			},
			expectedError: &voucher.NoMetadataError{
				Type: voucher.VulnerabilityType,
				Err:  vgrafeas.ErrNoOccurrences,
			},
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			grafeasMock := mocks.NewMockGrafeasAPIService(ctrl)
			client, _ := NewClient(context.Background(), project, project, pgp.NewKeyRing(), grafeasMock)
			grafeasMock.EXPECT().ListOccurrences(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.returnOccs, nil).AnyTimes()
			attestations, err := client.GetVulnerabilities(ctx, ref)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, attestations)
		})
	}
}

func TestGetBuildDetail(t *testing.T) {
	ctx := context.Background()
	project := "project"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grafeasMock := mocks.NewMockGrafeasAPIService(ctrl)
	client, _ := NewClient(context.Background(), project, project, pgp.NewKeyRing(), grafeasMock)
	ref := getCanonicalRef(t)
	occs := createAllOccurrences()
	tcs := map[string]struct {
		returnOccs     objects.ListOccurrencesResponse
		expectedResult repository.BuildDetail
		expectedError  error
	}{
		"valid input": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: occs,
			},
			expectedResult: repository.BuildDetail{
				RepositoryURL: "https://github.com/Shopify/voucher",
				Commit:        "2",
				BuildCreator:  "Shopify",
				BuildURL:      "some_url2",
				ProjectID:     "ShopifyId2",
				Artifacts: []repository.BuildArtifact{{
					Checksum: "sha256:71e3e78693c011e59b3fc84940f7672aeeb0a55427b6f5157bd08ab9e9ac746c",
					ID:       "some_id2",
				}},
			},
		},
		"no data": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: []objects.Occurrence{},
			},
			expectedError: &voucher.NoMetadataError{
				Type: voucher.BuildDetailsType,
				Err:  vgrafeas.ErrNoOccurrences,
			},
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			grafeasMock.EXPECT().ListOccurrences(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.returnOccs, nil)
			attestations, err := client.GetBuildDetail(ctx, ref)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, attestations)
		})
	}
}

func createAllOccurrences() []objects.Occurrence {
	noteKindVuln := objects.NoteKindVulnerability
	noteKindAtt := objects.NoteKindAttestation
	noteKindB := objects.NoteKindBuild
	noteKindD := objects.NoteKindDiscovery
	contentType := objects.AttestationUnspecified
	vulnSeverity := objects.SeverityMinimal
	packageKind := objects.VersionKindNormal
	activeContAnalysis := objects.DiscoveredContinuousAnalysisActive
	successStatus := objects.DiscoveredAnalysisStatusFinishedSuccess
	occs := []objects.Occurrence{
		{Name: "name1", Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
			NoteName: "notename", Kind: &noteKindVuln,
			Vulnerability: &objects.VulnerabilityDetails{Type: "vulnmedium", Severity: &vulnSeverity,
				PackageIssue: []objects.VulnerabilityPackageIssue{{
					AffectedLocation: &objects.VulnerabilityLocation{CpeURI: "uri", Package: "package_test",
						Version: &objects.PackageVersion{Name: "v0.0.0", Kind: &packageKind, Revision: "r"}}}}}},

		{Name: "name2", Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
			NoteName: "notename", Kind: &noteKindAtt,
			Attestation: &objects.AttestationDetails{Attestation: &objects.Attestation{
				GenericSignedAttestation: &objects.AttestationGenericSigned{ContentType: &contentType}}}},

		{Name: "name3", Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
			NoteName: "notename", Kind: &noteKindB,
			Build: &objects.BuildDetails{Provenance: &objects.ProvenanceBuild{ID: "provenceid", Creator: "Shopify",
				ProjectID: "ShopifyId2", LogsURI: "some_url2", BuiltArtifacts: []objects.ProvenanceArtifact{{
					Checksum: "sha256:71e3e78693c011e59b3fc84940f7672aeeb0a55427b6f5157bd08ab9e9ac746c", ID: "some_id2"}},
				SourceProvenance: &objects.ProvenanceSource{Context: &objects.SourceContext{
					Git: &objects.GitSourceContext{URL: "https://github.com/Shopify/voucher",
						RevisionID: "2"}}}}}},

		{Name: "name4", Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
			NoteName: "notename", Kind: &noteKindD,
			Discovered: &objects.DiscoveryDetails{
				Discovered: &objects.DiscoveryDiscovered{ContinuousAnalysis: &activeContAnalysis,
					AnalysisStatus: &successStatus}}},
	}
	return occs
}
