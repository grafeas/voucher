package grafeas

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/golang/mock/gomock"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	voucher "github.com/grafeas/voucher/v2"
	"github.com/grafeas/voucher/v2/attestation"
	mocks "github.com/grafeas/voucher/v2/grafeas/mocks"
	"github.com/grafeas/voucher/v2/grafeas/objects"
	"github.com/grafeas/voucher/v2/repository"
	"github.com/grafeas/voucher/v2/signer"
	"github.com/grafeas/voucher/v2/signer/kms"
	"github.com/grafeas/voucher/v2/signer/pgp"
	vtesting "github.com/grafeas/voucher/v2/testing"
)

var basePath string

const imgPath = "gcr.io/grafeas/grafeas-server@sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733"

func TestMain(m *testing.M) {
	flag.StringVar(&basePath, "grafeas", "", "the base path to the grafeas instance to use for testing")
	flag.Parse()
	os.Exit(m.Run())
}

func getCanonicalRef(t *testing.T, img string) reference.Canonical {
	named, err := reference.ParseNamed(img)
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
	ref := getCanonicalRef(t, imgPath)
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
	testSigner := vtesting.NewPGPSigner(t)

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
			reference:     getCanonicalRef(t, imgPath),
			keyring:       testSigner,
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
	ref := getCanonicalRef(t, imgPath)
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
				Err:  errNoOccurrences,
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
	occs := createAllOccurrences()
	errRef := getCanonicalRef(t, "us.gcr.io/grafeas/grafeas-server@sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733")
	validRef := getCanonicalRef(t, imgPath)
	successStatus := objects.DiscoveredAnalysisStatusFinishedSuccess
	noteKindD := objects.NoteKindDiscovery
	setPollOptions(1, 0)
	tcs := map[string]struct {
		returnOccs       objects.ListOccurrencesResponse
		expectedResult   []voucher.Vulnerability
		expectedError    error
		expectedErrorStr string
		ref              reference.Canonical
	}{
		"valid input": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: occs,
			},
			ref: validRef,
			expectedResult: []voucher.Vulnerability{{
				Name:     "notename",
				Severity: voucher.NegligibleSeverity,
			}},
		},
		"no vulnerability data": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: []objects.Occurrence{{Name: "name4",
					Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
					NoteName: "notename_invalid", Kind: &noteKindD,
					Discovered: &objects.DiscoveryDetails{
						Discovered: &objects.DiscoveryDiscovered{AnalysisStatus: &successStatus}}}},
			},
			expectedError: &voucher.NoMetadataError{
				Type: voucher.VulnerabilityType,
				Err:  errNoOccurrences,
			},
			ref: validRef,
		},
		"no data": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: []objects.Occurrence{},
			},
			expectedError:  errDiscoveriesUnfinished,
			expectedResult: []voucher.Vulnerability{},
			ref:            validRef,
		},
		"image ref error": {
			expectedErrorStr: fmt.Sprintf("could not find project path in reference \"%s\"", errRef),
			expectedResult:   []voucher.Vulnerability{},
			ref:              errRef,
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			grafeasMock := mocks.NewMockGrafeasAPIService(ctrl)
			client, _ := NewClient(context.Background(), project, project, pgp.NewKeyRing(), grafeasMock)
			grafeasMock.EXPECT().ListOccurrences(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.returnOccs, nil).AnyTimes()
			attestations, err := client.GetVulnerabilities(ctx, test.ref)
			if test.expectedErrorStr == "" {
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.Equal(t, test.expectedErrorStr, err.Error())
			}
			assert.Equal(t, test.expectedResult, attestations)
			client.Close()
		})
	}
	defaultPollOptions()
}

func TestGetBuildDetail(t *testing.T) {
	ctx := context.Background()
	project := "project"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	occs := createAllOccurrences()
	errRef := getCanonicalRef(t, "us.gcr.io/grafeas/grafeas-server@sha256:c7303bdd6e36868d54b5b00dee125445a8d0f667c366420ccbe41dcf3b1c7733")
	validRef := getCanonicalRef(t, imgPath)
	tcs := map[string]struct {
		returnOccs       objects.ListOccurrencesResponse
		expectedResult   repository.BuildDetail
		expectedError    error
		ref              reference.Canonical
		expectedErrorStr string
	}{
		"valid input": {
			returnOccs: objects.ListOccurrencesResponse{
				Occurrences: occs,
			},
			ref: validRef,
			expectedResult: repository.BuildDetail{
				RepositoryURL: "https://github.com/grafeas/voucher",
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
				Err:  errNoOccurrences,
			},
			ref: validRef,
		},
		"image ref error": {
			expectedErrorStr: fmt.Sprintf("could not find project path in reference \"%s\"", errRef),
			expectedResult:   repository.BuildDetail{},
			ref:              errRef,
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			grafeasMock := mocks.NewMockGrafeasAPIService(ctrl)
			client, _ := NewClient(context.Background(), project, project, pgp.NewKeyRing(), grafeasMock)
			grafeasMock.EXPECT().ListOccurrences(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.returnOccs, nil).AnyTimes()
			attestations, err := client.GetBuildDetail(ctx, test.ref)
			if test.expectedErrorStr == "" {
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.Equal(t, test.expectedErrorStr, err.Error())
			}
			assert.Equal(t, test.expectedResult, attestations)
			client.Close()
		})
	}
}

func createAllOccurrences() []objects.Occurrence {
	noteKindVuln := objects.NoteKindVulnerability
	noteKindAtt := objects.NoteKindAttestation
	noteKindB := objects.NoteKindBuild
	noteKindD := objects.NoteKindDiscovery
	contentType := objects.AttestationUnspecified
	vulnSeverity := objects.SeverityLow
	vulnEffectiveSeverity := objects.SeverityMinimal
	packageKind := objects.VersionKindNormal
	successStatus := objects.DiscoveredAnalysisStatusFinishedSuccess
	occs := []objects.Occurrence{
		{Name: "name1", Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
			NoteName: "notename", Kind: &noteKindVuln,
			Vulnerability: &objects.VulnerabilityDetails{Type: "vulnmedium", Severity: &vulnSeverity,
				EffectiveSeverity: &vulnEffectiveSeverity,
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
					Git: &objects.GitSourceContext{URL: "https://github.com/grafeas/voucher",
						RevisionID: "2"}}}}}},

		{Name: "name4", Resource: &objects.Resource{URI: "https://gcr.io/project/image@sha256:foo"},
			NoteName: "notename", Kind: &noteKindD,
			Discovered: &objects.DiscoveryDetails{
				Discovered: &objects.DiscoveryDiscovered{AnalysisStatus: &successStatus}}},
	}
	return occs
}
