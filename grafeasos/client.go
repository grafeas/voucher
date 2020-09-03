package grafeasos

import (
	"context"
	"errors"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/attestation"
	"github.com/Shopify/voucher/containeranalysis"
	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/signer"
	"github.com/antihax/optional"
	"github.com/docker/distribution/reference"
	grafeaspb "github.com/grafeas/client-go/0.1.0"
)

var errCannotAttest = errors.New("cannot create attestations, keyring is empty")

// Client implements voucher.MetadataClient, connecting to Grafeas.
type Client struct {
	grafeas        *grafeaspb.GrafeasV1Beta1ApiService // The client reference.
	keyring        signer.AttestationSigner            // The keyring used for signing metadata.
	binauthProject string                              // The project that Binauth Notes and Occurrences are written to.
	imageProject   string                              // The project that image information is stored.
}

// CanAttest returns true if the client can create and sign attestations.
func (g *Client) CanAttest() bool {
	return nil != g.keyring
}

// NewPayloadBody returns a payload body appropriate for this MetadataClient.
func (g *Client) NewPayloadBody(reference reference.Canonical) (string, error) {
	payload, err := attestation.NewPayload(reference).ToString()
	if err != nil {
		return "", err
	}

	return payload, err
}

// AddAttestationToImage adds a new attestation with the passed Attestation
// to the image described by ImageData.
func (g *Client) AddAttestationToImage(ctx context.Context, reference reference.Canonical, payload voucher.Attestation) (voucher.SignedAttestation, error) {
	signedAttestation := voucher.SignedAttestation{}

	if !g.CanAttest() {
		return signedAttestation, errCannotAttest
	}

	signed, err := voucher.SignAttestation(g.keyring, payload)
	if nil != err {
		return signedAttestation, err
	}

	binauthProjectPath := "projects/" + g.binauthProject

	contentType := grafeaspb.SIMPLE_SIGNING_JSON_AttestationPgpSignedAttestationContentType

	attestation := grafeaspb.V1beta1attestationDetails{Attestation: &grafeaspb.AttestationAttestation{PgpSignedAttestation: &grafeaspb.AttestationPgpSignedAttestation{Signature: signed.Signature,
		PgpKeyId: signed.KeyID, ContentType: &contentType}}}

	occurrence := g.getCreateOccurrence(reference, payload.CheckName, &attestation, binauthProjectPath)
	_, _, err = g.grafeas.CreateOccurrence(ctx, binauthProjectPath, occurrence)

	if isAttestionExistsErr(err) {
		err = nil
	}

	return signed, err
}

// OccurrenceToAttestation converts an Occurrence to a Attestation
func OccurrenceToAttestation(checkName string, occ *grafeaspb.V1beta1Occurrence) voucher.SignedAttestation {
	signedAttestation := voucher.SignedAttestation{
		Attestation: voucher.Attestation{
			CheckName: checkName,
		},
	}

	attestationDetails := occ.Attestation

	signedAttestation.Body = string(*attestationDetails.Attestation.GenericSignedAttestation.ContentType)

	return signedAttestation
}

// GetAttestations returns all of the attestations associated with an image
func (g *Client) GetAttestations(ctx context.Context, reference reference.Canonical) ([]voucher.SignedAttestation, error) {
	// filterStr := kindFilterStr(reference, grafeaspb.ATTESTATION_V1beta1NoteKind)

	var occ []grafeaspb.V1beta1Occurrence

	var attestations []voucher.SignedAttestation

	project := projectPath(g.binauthProject)

	createOccs := grafeaspb.V1beta1BatchCreateOccurrencesRequest{Parent: project, Occurrences: occ}

	occs, httpResponse, err := g.grafeas.BatchCreateOccurrences(ctx, project, createOccs)

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 200 {
		return nil, err
	}

	for _, oc := range occs.Occurrences {
		note := oc.NoteName
		attestations = append(
			attestations,
			OccurrenceToAttestation(note, &oc),
		)
	}

	return attestations, nil
}

func (g *Client) getVulnerabilityDiscoveries(ctx context.Context) (items []grafeaspb.V1beta1Occurrence, err error) {
	occurrences, err := g.getNotesOccurrencesForKind(ctx, grafeaspb.DISCOVERY_V1beta1NoteKind, isDiscoveryVulnerabilityNote)

	items = append(items, occurrences...)

	if 0 == len(items) && nil == err {
		err = &voucher.NoMetadataError{
			Type: containeranalysis.DiscoveryType,
			Err:  errNoOccurrences,
		}
	}
	return
}

func (g *Client) getNotesOccurrencesForKind(ctx context.Context, noteKind grafeaspb.V1beta1NoteKind, fn noteTest) (items []grafeaspb.V1beta1Occurrence, err error) {
	optsNotes := &grafeaspb.ListNotesOpts{}
	project := projectPath(g.imageProject)
	notesResponse, httpResponse, err := g.grafeas.ListNotes(ctx, project, optsNotes)
	if err != nil {
		return
	}

	if httpResponse.StatusCode != 200 {
		return
	}

	for {
		for _, note := range notesResponse.Notes {
			if fn(&note, noteKind) {
				noteOccurrences, errOcc := g.getOccurrencesForNote(ctx, note.Name, noteKind)
				if errOcc != nil && err == nil {
					err = errOcc
					continue
				}
				items = append(items, noteOccurrences...)
			}
		}
		if notesResponse.NextPageToken == "" {
			break
		}

		optsNotes.PageToken = optional.NewString(notesResponse.NextPageToken)
		notesResponse, _, err = g.grafeas.ListNotes(ctx, project, optsNotes)
		if nil != err {
			break
		}
	}

	return
}

func (g *Client) getOccurrencesForNote(ctx context.Context, noteName string, noteKind grafeaspb.V1beta1NoteKind) (items []grafeaspb.V1beta1Occurrence, err error) {
	optsOccurrences := &grafeaspb.ListNoteOccurrencesOpts{PageSize: optional.NewInt32(5)}
	occsResponse, _, errOcc := g.grafeas.ListNoteOccurrences(ctx, noteName, optsOccurrences)
	if errOcc != nil {
		err = errOcc
		return
	}
	for {
		for i, occ := range occsResponse.Occurrences {
			if noteKind == *occ.Kind {
				items = append(items, occsResponse.Occurrences[i])
			}
		}
		if occsResponse.NextPageToken == "" {
			break
		}

		optsOccurrences.PageToken = optional.NewString(occsResponse.NextPageToken)
		occsResponse, _, err = g.grafeas.ListNoteOccurrences(ctx, noteName, optsOccurrences)
		if nil != err {
			break
		}
	}
	return
}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (g *Client) GetVulnerabilities(ctx context.Context, reference reference.Canonical) (items []voucher.Vulnerability, err error) {
	err = pollForDiscoveries(ctx, g)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}

	kind := grafeaspb.VULNERABILITY_V1beta1NoteKind
	occurrences, err := g.getNotesOccurrencesForKind(ctx, kind, isTypeNote)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}
	for _, occ := range occurrences {
		item := OccurrenceToVulnerability(&occ)
		items = append(items, item)
	}

	if 0 == len(items) && nil == err {
		err = &voucher.NoMetadataError{
			Type: voucher.VulnerabilityType,
			Err:  errNoOccurrences,
		}
	}

	return
}

// Close closes the Grafeas client.
func (g *Client) Close() {}

// GetBuildDetail gets BuildDetails for the passed image.
func (g *Client) GetBuildDetail(ctx context.Context, reference reference.Canonical) (repository.BuildDetail, error) {
	kind := grafeaspb.BUILD_V1beta1NoteKind
	occurrences, err := g.getNotesOccurrencesForKind(ctx, kind, isTypeNote)
	if nil != err {
		return repository.BuildDetail{}, err
	}

	// we should only have 1 occurrence based on our kind specified
	if nil == err && len(occurrences) != 1 {
		if len(occurrences) == 0 {
			return repository.BuildDetail{}, &voucher.NoMetadataError{Type: voucher.BuildDetailsType, Err: errNoOccurrences}
		}

		return repository.BuildDetail{}, errors.New("Found multiple Grafeas occurrences for " + reference.String())
	}

	return OccurrenceToBuildDetail(&occurrences[0]), nil
}

// NewClient creates a new Grafeas Client.
func NewClient(ctx context.Context, imageProject, binauthProject string, keyring signer.AttestationSigner, config *grafeaspb.Configuration) (*Client, error) {
	client := new(Client)
	client.grafeas = grafeaspb.NewAPIClient(config).GrafeasV1Beta1Api

	client.keyring = keyring
	client.binauthProject = binauthProject
	client.imageProject = imageProject

	return client, nil
}
