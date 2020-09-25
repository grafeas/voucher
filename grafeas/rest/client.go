package rest

import (
	"context"
	"errors"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/attestation"
	vgrafeas "github.com/Shopify/voucher/grafeas"
	"github.com/Shopify/voucher/grafeas/rest/objects"
	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/signer"
	"github.com/antihax/optional"
	"github.com/docker/distribution/reference"
)

var errCannotAttest = errors.New("cannot create attestations, keyring is empty")

// vulProject is the project that Google's Container Analysis writes vulnerability
// occurrences to.
const vulProject = "projects/goog-vulnz/notes/"

// Client implements voucher.MetadataClient, connecting to Grafeas.
type Client struct {
	grafeas        GrafeasAPIService        // The client reference.
	keyring        signer.AttestationSigner // The keyring used for signing metadata.
	binauthProject string                   // The project that Binauth Notes and Occurrences are written to.
	imageProject   string                   // The project that image information is stored.
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
	if !g.CanAttest() {
		return voucher.SignedAttestation{}, errCannotAttest
	}

	signedAttestation, err := voucher.SignAttestation(g.keyring, payload)
	if nil != err {
		return voucher.SignedAttestation{}, err
	}

	binauthProjectPath := vgrafeas.ProjectPath(g.binauthProject)
	contentType := objects.AttestationSigningJSON

	attestation := objects.AttestationDetails{Attestation: &objects.Attestation{
		PgpSignedAttestation: &objects.AttestationPgpSigned{Signature: signedAttestation.Signature,
			PgpKeyID: signedAttestation.KeyID, ContentType: &contentType}}}

	occurrence := objects.NewOccurrence(reference, payload.CheckName, &attestation, binauthProjectPath)
	_, err = g.grafeas.CreateOccurrence(ctx, binauthProjectPath, occurrence)

	if vgrafeas.IsAttestationExistsErr(err) {
		err = nil
		signedAttestation.Signature = ""
	}

	return signedAttestation, err
}

// GetAttestations returns all of the attestations associated with an image
func (g *Client) GetAttestations(ctx context.Context, reference reference.Canonical) ([]voucher.SignedAttestation, error) {
	var attestations []voucher.SignedAttestation

	occurrences, err := g.getAllOccurrences(ctx)
	if err != nil {
		return []voucher.SignedAttestation{}, err
	}
	for _, occ := range occurrences {
		if *occ.Kind != objects.NoteKindAttestation {
			continue
		}
		note := occ.NoteName
		attestations = append(
			attestations,
			objects.OccurrenceToAttestation(note, &occ),
		)
	}

	if 0 == len(attestations) && nil == err {
		err = &voucher.NoMetadataError{
			Type: voucher.AttestationType,
			Err:  vgrafeas.ErrNoOccurrences,
		}
	}

	return attestations, err
}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (g *Client) GetVulnerabilities(ctx context.Context, reference reference.Canonical) (items []voucher.Vulnerability, err error) {
	err = pollForDiscoveries(ctx, g)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}

	occurrences, err := g.getAllOccurrences(ctx)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}
	for _, occ := range occurrences {
		if *occ.Kind != objects.NoteKindVulnerability {
			continue
		}
		item := objects.OccurrenceToVulnerability(&occ, vulProject)
		items = append(items, item)
	}

	if 0 == len(items) && nil == err {
		err = &voucher.NoMetadataError{
			Type: voucher.VulnerabilityType,
			Err:  vgrafeas.ErrNoOccurrences,
		}
	}

	return
}

// Close closes the Grafeas client.
func (g *Client) Close() {}

// GetBuildDetail gets BuildDetails for the passed image.
func (g *Client) GetBuildDetail(ctx context.Context, reference reference.Canonical) (repository.BuildDetail, error) {
	items, err := g.getAllOccurrences(ctx)
	if nil != err {
		return repository.BuildDetail{}, err
	}
	occurrences := []objects.Occurrence{}
	for _, occ := range items {
		if *occ.Kind == objects.NoteKindBuild {
			occurrences = append(occurrences, occ)
		}
	}

	// we should only have 1 occurrence based on our kind specified
	if nil == err && len(occurrences) != 1 {
		if len(occurrences) == 0 {
			return repository.BuildDetail{}, &voucher.NoMetadataError{Type: voucher.BuildDetailsType, Err: vgrafeas.ErrNoOccurrences}
		}

		return repository.BuildDetail{}, errors.New("Found multiple Grafeas occurrences for " + reference.String())
	}

	return objects.OccurrenceToBuildDetail(&occurrences[0]), nil
}

func (g *Client) getAllOccurrences(ctx context.Context) (items []objects.Occurrence, err error) {
	project := vgrafeas.ProjectPath(g.binauthProject)

	occResp, err := g.grafeas.ListOccurrences(ctx, project, nil)
	if err != nil {
		return nil, err
	}

	items = append(items, occResp.Occurrences...)

	for occResp.NextPageToken != "" {
		occResp, err = g.grafeas.ListOccurrences(ctx, project, &objects.ListOpts{
			PageToken: optional.NewString(occResp.NextPageToken),
		})
		if err != nil {
			return nil, err
		}
		items = append(items, occResp.Occurrences...)
	}

	return
}

// NewClient creates a new Grafeas Client.
func NewClient(ctx context.Context, imageProject, binauthProject string, keyring signer.AttestationSigner, grafeas GrafeasAPIService) (*Client, error) {
	return &Client{
		grafeas:        grafeas,
		keyring:        keyring,
		binauthProject: binauthProject,
		imageProject:   imageProject,
	}, nil
}
