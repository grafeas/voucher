package grafeas

import (
	"context"
	"errors"

	containeranalysis "cloud.google.com/go/devtools/containeranalysis/apiv1alpha1"
	"github.com/Shopify/voucher"
	binauth "github.com/Shopify/voucher/grafeas/binauth"
	"github.com/docker/distribution/reference"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

var errCannotAttest = errors.New("cannot create attestations, keyring is empty")

// Client implements voucher.MetadataClient, connecting to Grafeas.
type Client struct {
	ctx            context.Context  // The context to use for new connections created in this client.
	keyring        *voucher.KeyRing // The keyring used for signing metadata.
	binauthProject string           // The project that Binauth Notes and Occurrences are written to.
	imageProject   string           // The project that image information is stored.
}

// CanAttest returns true if the client can create and sign attestations.
func (g *Client) CanAttest() bool {
	return nil != g.keyring
}

// NewPayloadBody returns a payload body appropriate for this MetadataClient.
func (g *Client) NewPayloadBody(reference reference.Canonical) (string, error) {
	payload, err := binauth.NewPayload(reference).ToString()
	if err != nil {
		return "", err
	}
	return payload, err
}

// GetMetadata gets metadata of the requested type for the passed image.
func (g *Client) GetMetadata(reference reference.Canonical, metadataType voucher.MetadataType) (items []voucher.MetadataItem, err error) {
	c, err := containeranalysis.NewClient(g.ctx)
	if err != nil {
		return
	}

	filterStr := resourceURL(reference)

	kind := getNoteKind(metadataType)
	if kind != containeranalysispb.Note_KIND_UNSPECIFIED {
		filterStr = kindFilterStr(reference, kind)
	}

	project := projectPath(g.imageProject)
	req := &containeranalysispb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	iterator := c.ListOccurrences(g.ctx, req)
	for occ, complete := iterator.Next(); complete == nil; occ, complete = iterator.Next() {
		item := new(Item)
		item.Occurrence = occ
		items = append(items, item)
	}

	if 0 == len(items) {
		err = errNoOccurrences
	}

	return
}

// AddAttestationToImage adds a new attestation with the passed AttestationPayload
// to the image described by ImageData.
func (g *Client) AddAttestationToImage(reference reference.Canonical, payload voucher.AttestationPayload) (voucher.MetadataItem, error) {
	if !g.CanAttest() {
		return nil, errCannotAttest
	}

	signed, keyID, err := payload.Sign(g.keyring)
	if nil != err {
		return nil, err
	}

	attestation := g.getOccurrenceAttestation(signed, keyID)
	occurrenceRequest := g.getCreateOccurrenceRequest(reference, payload.CheckName, attestation)
	c, err := containeranalysis.NewClient(g.ctx)
	if err != nil {
		return nil, err
	}
	occ, err := c.CreateOccurrence(g.ctx, occurrenceRequest)

	item := new(Item)
	item.Occurrence = occ

	return item, err
}

func (g *Client) getOccurrenceAttestation(signature string, keyID string) *containeranalysispb.Occurrence_Attestation {
	pgpKeyID := containeranalysispb.PgpSignedAttestation_PgpKeyId{keyID}
	pgpSignedAttestation := containeranalysispb.PgpSignedAttestation{signature, 1, &pgpKeyID}
	attestationAuthoritySignedAttestation := containeranalysispb.AttestationAuthority_Attestation_PgpSignedAttestation{&pgpSignedAttestation}
	attestationAuthorityAttestation := containeranalysispb.AttestationAuthority_Attestation{&attestationAuthoritySignedAttestation}
	occurrenceAttestation := containeranalysispb.Occurrence_Attestation{&attestationAuthorityAttestation}
	return &occurrenceAttestation
}

func (g *Client) getCreateOccurrenceRequest(reference reference.Reference, parentNoteID string, attestation *containeranalysispb.Occurrence_Attestation) *containeranalysispb.CreateOccurrenceRequest {
	binauthProjectPath := "projects/" + g.binauthProject
	noteName := binauthProjectPath + "/notes/" + parentNoteID
	resourceURL := "https://" + reference.String()
	occurrence := containeranalysispb.Occurrence{NoteName: noteName, ResourceUrl: resourceURL, Details: attestation}
	req := &containeranalysispb.CreateOccurrenceRequest{Parent: binauthProjectPath, Occurrence: &occurrence}
	return req
}

// NewClient creates a new Grafeas Client.
func NewClient(ctx context.Context, imageProject, binauthProject string, keyring *voucher.KeyRing) *Client {
	client := new(Client)
	client.ctx = ctx
	client.keyring = keyring
	client.binauthProject = binauthProject
	client.imageProject = imageProject
	return client
}
