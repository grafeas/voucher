package grafeas

import (
	"context"
	"errors"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	"github.com/Shopify/voucher"
	binauth "github.com/Shopify/voucher/grafeas/binauth"
	"github.com/docker/distribution/reference"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/attestation"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/common"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
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
	c, err := containeranalysis.NewGrafeasV1Beta1Client(g.ctx)
	if err != nil {
		return
	}
	defer c.Close()

	filterStr := resourceURL(reference)

	kind := getNoteKind(metadataType)
	if kind != common.NoteKind_NOTE_KIND_UNSPECIFIED {
		filterStr = kindFilterStr(reference, kind)
	}

	project := projectPath(g.imageProject)
	req := &grafeaspb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	occIterator := c.ListOccurrences(g.ctx, req)
	for {
		var occ *grafeaspb.Occurrence
		occ, err = occIterator.Next()
		if nil != err {
			if iterator.Done == err {
				err = nil
			}
			break
		}
		item := new(Item)
		item.Occurrence = occ
		items = append(items, item)
	}

	if 0 == len(items) && nil == err {
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
	c, err := containeranalysis.NewGrafeasV1Beta1Client(g.ctx)
	if err != nil {
		return nil, err
	}
	occ, err := c.CreateOccurrence(g.ctx, occurrenceRequest)

	item := new(Item)
	item.Occurrence = occ

	return item, err
}

func (g *Client) getOccurrenceAttestation(signature string, keyID string) *grafeaspb.Occurrence_Attestation {
	pgpKeyID := attestation.PgpSignedAttestation_PgpKeyId{
		PgpKeyId: keyID,
	}

	pgpSignedAttestation := attestation.PgpSignedAttestation{
		Signature:   signature,
		ContentType: attestation.PgpSignedAttestation_SIMPLE_SIGNING_JSON,
		KeyId:       &pgpKeyID,
	}

	attestationPgpSignedAttestation := attestation.Attestation_PgpSignedAttestation{
		PgpSignedAttestation: &pgpSignedAttestation,
	}

	newAttestation := attestation.Attestation{
		Signature: &attestationPgpSignedAttestation,
	}

	details := attestation.Details{
		Attestation: &newAttestation,
	}

	occurrenceAttestation := grafeaspb.Occurrence_Attestation{
		Attestation: &details,
	}

	return &occurrenceAttestation
}

func (g *Client) getCreateOccurrenceRequest(reference reference.Reference, parentNoteID string, attestation *grafeaspb.Occurrence_Attestation) *grafeaspb.CreateOccurrenceRequest {
	binauthProjectPath := "projects/" + g.binauthProject
	noteName := binauthProjectPath + "/notes/" + parentNoteID

	resource := grafeaspb.Resource{
		Uri: "https://" + reference.String(),
	}

	occurrence := grafeaspb.Occurrence{NoteName: noteName, Resource: &resource, Details: attestation}
	req := &grafeaspb.CreateOccurrenceRequest{Parent: binauthProjectPath, Occurrence: &occurrence}
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
