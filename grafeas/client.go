package grafeas

import (
	"context"
	"errors"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/attestation"
	"github.com/Shopify/voucher/docker/uri"
	"github.com/Shopify/voucher/grafeas/objects"
	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/signer"
	"github.com/antihax/optional"
	"github.com/docker/distribution/reference"
)

var errCannotAttest = errors.New("cannot create attestations, keyring is empty")

// Client implements voucher.MetadataClient, connecting to Grafeas.
type Client struct {
	service        APIService               // The client reference.
	keyring        signer.AttestationSigner // The keyring used for signing metadata.
	binauthProject string                   // The project that Binauth Notes and Occurrences are written to.
	vulProject     string                   // The project to read vulnerability occurrences from.
}

// CanAttest returns true if the client can create and sign attestations.
func (g *Client) CanAttest() bool {
	return nil != g.keyring
}

// NewPayloadBody returns a payload body appropriate for this MetadataClient.
func (g *Client) NewPayloadBody(ref reference.Canonical) (string, error) {
	payload, err := attestation.NewPayload(ref).ToString()
	if err != nil {
		return "", err
	}

	return payload, err
}

// AddAttestationToImage adds a new attestation with the passed Attestation
// to the image described by ImageData.
func (g *Client) AddAttestationToImage(ctx context.Context, ref reference.Canonical, payload voucher.Attestation) (voucher.SignedAttestation, error) {
	if !g.CanAttest() {
		return voucher.SignedAttestation{}, errCannotAttest
	}

	signedAttestation, err := voucher.SignAttestation(g.keyring, payload)
	if nil != err {
		return voucher.SignedAttestation{}, err
	}

	binauthProjectPath := projectPath(g.binauthProject)

	occurrence := objects.NewOccurrence(ref, payload.CheckName, objects.NewAttestation(signedAttestation), binauthProjectPath)
	_, err = g.service.CreateOccurrence(ctx, binauthProjectPath, occurrence)

	if isAttestationExistsErr(err) {
		err = nil
		signedAttestation.Signature = ""
	}

	return signedAttestation, err
}

// GetAttestations returns all of the attestations associated with an image
func (g *Client) GetAttestations(ctx context.Context, ref reference.Canonical) ([]voucher.SignedAttestation, error) {
	var attestations []voucher.SignedAttestation

	occurrences, err := g.getAllOccurrences(ctx, g.binauthProject)
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
			occ.Attestation.AsVoucherAttestation(note),
		)
	}

	if 0 == len(attestations) && nil == err {
		err = &voucher.NoMetadataError{
			Type: voucher.AttestationType,
			Err:  errNoOccurrences,
		}
	}

	return attestations, err
}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (g *Client) GetVulnerabilities(ctx context.Context, ref reference.Canonical) (items []voucher.Vulnerability, err error) {
	err = pollForDiscoveries(ctx, g, ref)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}

	project, err := uri.ReferenceToProjectName(ref)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}

	occurrences, err := g.getAllOccurrences(ctx, project)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}
	for _, occ := range occurrences {
		if *occ.Kind != objects.NoteKindVulnerability {
			continue
		}
		item := occ.Vulnerability.AsVoucherVulnerability(occ.NoteName, g.vulProject)
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
func (g *Client) GetBuildDetail(ctx context.Context, ref reference.Canonical) (repository.BuildDetail, error) {
	project, err := uri.ReferenceToProjectName(ref)
	if nil != err {
		return repository.BuildDetail{}, err
	}

	items, err := g.getAllOccurrences(ctx, project)
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
			return repository.BuildDetail{}, &voucher.NoMetadataError{Type: voucher.BuildDetailsType, Err: errNoOccurrences}
		}

		return repository.BuildDetail{}, errors.New("Found multiple occurrences for " + ref.String())
	}

	return occurrences[0].Build.AsVoucherBuildDetail(), nil
}

func (g *Client) getAllOccurrences(ctx context.Context, path string) (items []objects.Occurrence, err error) {
	project := projectPath(path)

	occResp, err := g.service.ListOccurrences(ctx, project, nil)
	if err != nil {
		return nil, err
	}

	items = append(items, occResp.Occurrences...)

	for occResp.NextPageToken != "" {
		occResp, err = g.service.ListOccurrences(ctx, project, &objects.ListOpts{
			PageToken: optional.NewString(occResp.NextPageToken),
		})
		if err != nil {
			return nil, err
		}
		items = append(items, occResp.Occurrences...)
	}

	return
}

func projectPath(project string) string {
	return "projects/" + project
}

// NewClient creates a new Grafeas Client.
func NewClient(ctx context.Context, binauthProject, vulProject string, keyring signer.AttestationSigner, service APIService) (*Client, error) {
	return &Client{
		service:        service,
		keyring:        keyring,
		binauthProject: binauthProject,
		vulProject:     vulProject,
	}, nil
}
