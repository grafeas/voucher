package grafeas

import (
	"context"
	"errors"

	containeranalysis "cloud.google.com/go/devtools/containeranalysis/apiv1alpha1"
	"github.com/Shopify/voucher"
	binauth "github.com/Shopify/voucher/grafeas/binauth"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

var errCannotAttest = errors.New("cannot create attestations, keyring is empty")

// Client implements voucher.MetadataClient, connecting to Grafeas.
type Client struct {
	keyring        *voucher.KeyRing // The keyring used for signing metadata.
	binauthProject string           // The project that Binauth Notes and Occurrences are written to.
	imageProject   string           // The project that image information is stored.
}

// CanAttest returns true if the client can create and sign attestations.
func (g *Client) CanAttest() bool {
	return nil != g.keyring
}

// NewPayloadBody returns a payload body appropriate for this MetadataClient.
func (g *Client) NewPayloadBody(imageData voucher.ImageData) (string, error) {
	payload, err := binauth.NewPayload(imageData).ToString()
	if err != nil {
		return "", err
	}
	return payload, err
}

// GetOccurrencesForImage gets the occurrences
func (g *Client) GetOccurrencesForImage(imageData voucher.ImageData, kind voucher.NoteKind) (occurrences []voucher.Occurrence, err error) {
	ctx := context.Background()
	c, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return
	}

	filterStr := resourceURL(imageData)
	if kind != containeranalysispb.Note_KIND_UNSPECIFIED {
		filterStr = kindFilterStr(imageData, kind)
	}

	project := projectPath(g.imageProject)
	req := &containeranalysispb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	iterator := c.ListOccurrences(ctx, req)
	for occ, complete := iterator.Next(); complete == nil; occ, complete = iterator.Next() {
		occurrences = append(occurrences, occ)
	}

	if 0 == len(occurrences) {
		err = errNoOccurrences
	}

	return
}

// AddAttestationOccurrenceToImage adds a new attestation with the passed AttestationPayload
// to the image described by ImageData.
func (g *Client) AddAttestationOccurrenceToImage(imageData voucher.ImageData, payload voucher.AttestationPayload) (voucher.Occurrence, error) {
	if !g.CanAttest() {
		return nil, errCannotAttest
	}

	signed, keyId, err := payload.Sign(g.keyring)
	if nil != err {
		return nil, err
	}

	attestation := g.getOccurrenceAttestation(signed, keyId)
	occurrenceRequest := g.getCreateOccurrenceRequest(imageData, payload.CheckName, attestation)
	ctx := context.Background()
	c, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.CreateOccurrence(ctx, occurrenceRequest)

}

func (g *Client) getOccurrenceAttestation(signature string, keyID string) *containeranalysispb.Occurrence_Attestation {
	pgpKeyID := containeranalysispb.PgpSignedAttestation_PgpKeyId{keyID}
	pgpSignedAttestation := containeranalysispb.PgpSignedAttestation{signature, 1, &pgpKeyID}
	attestationAuthoritySignedAttestation := containeranalysispb.AttestationAuthority_Attestation_PgpSignedAttestation{&pgpSignedAttestation}
	attestationAuthorityAttestation := containeranalysispb.AttestationAuthority_Attestation{&attestationAuthoritySignedAttestation}
	occurrenceAttestation := containeranalysispb.Occurrence_Attestation{&attestationAuthorityAttestation}
	return &occurrenceAttestation
}

func (g *Client) getCreateOccurrenceRequest(imageData voucher.ImageData, parentNoteID string, attestation *containeranalysispb.Occurrence_Attestation) *containeranalysispb.CreateOccurrenceRequest {
	binauthProjectPath := "projects/" + g.binauthProject
	noteName := binauthProjectPath + "/notes/" + parentNoteID
	occurrence := containeranalysispb.Occurrence{NoteName: noteName, ResourceUrl: imageData.String(), Details: attestation}
	req := &containeranalysispb.CreateOccurrenceRequest{Parent: binauthProjectPath, Occurrence: &occurrence}
	return req
}

// NewClient creates a new Grafeas Client.
func NewClient(imageProject, binauthProject string, keyring *voucher.KeyRing) *Client {
	client := new(Client)
	client.keyring = keyring
	client.binauthProject = binauthProject
	client.imageProject = imageProject
	return client
}
