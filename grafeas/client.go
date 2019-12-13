package grafeas

import (
	"context"
	"errors"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	"github.com/docker/distribution/reference"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/common"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"

	"github.com/Shopify/voucher"
	binauth "github.com/Shopify/voucher/grafeas/binauth"
	"github.com/Shopify/voucher/repository"
)

var errCannotAttest = errors.New("cannot create attestations, keyring is empty")

// Client implements voucher.MetadataClient, connecting to Grafeas.
type Client struct {
	grafeas        *containeranalysis.GrafeasV1Beta1Client // The client reference.
	keyring        *voucher.KeyRing                        // The keyring used for signing metadata.
	binauthProject string                                  // The project that Binauth Notes and Occurrences are written to.
	imageProject   string                                  // The project that image information is stored.
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
func (g *Client) GetMetadata(ctx context.Context, reference reference.Canonical, metadataType voucher.MetadataType) (items []voucher.MetadataItem, err error) {
	filterStr := resourceURL(reference)

	kind := getNoteKind(metadataType)
	if kind != common.NoteKind_NOTE_KIND_UNSPECIFIED {
		filterStr = kindFilterStr(reference, kind)
	}

	project := projectPath(g.imageProject)
	req := &grafeaspb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	occIterator := g.grafeas.ListOccurrences(ctx, req)
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
		err = &voucher.NoMetadataError{
			Type: metadataType,
			Err:  errNoOccurrences,
		}
	}

	return
}

// AddAttestationToImage adds a new attestation with the passed AttestationPayload
// to the image described by ImageData.
func (g *Client) AddAttestationToImage(ctx context.Context, reference reference.Canonical, payload voucher.AttestationPayload) (voucher.MetadataItem, error) {
	if !g.CanAttest() {
		return nil, errCannotAttest
	}

	signed, keyID, err := payload.Sign(g.keyring)
	if nil != err {
		return nil, err
	}

	attestation := newOccurrenceAttestation(signed, keyID)
	occurrenceRequest := g.getCreateOccurrenceRequest(reference, payload.CheckName, attestation)
	occ, err := g.grafeas.CreateOccurrence(ctx, occurrenceRequest)
	item := new(Item)
	item.Occurrence = occ

	if isAttestionExistsErr(err) {
		err = nil
		item.Occurrence = nil
	}

	return item, err
}

func (g *Client) getCreateOccurrenceRequest(reference reference.Canonical, parentNoteID string, attestation *grafeaspb.Occurrence_Attestation) *grafeaspb.CreateOccurrenceRequest {
	binauthProjectPath := "projects/" + g.binauthProject
	noteName := binauthProjectPath + "/notes/" + parentNoteID

	resource := grafeaspb.Resource{
		Uri: "https://" + reference.Name() + "@" + reference.Digest().String(),
	}

	occurrence := grafeaspb.Occurrence{NoteName: noteName, Resource: &resource, Details: attestation}
	req := &grafeaspb.CreateOccurrenceRequest{Parent: binauthProjectPath, Occurrence: &occurrence}
	return req
}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (g *Client) GetVulnerabilities(ctx context.Context, reference reference.Canonical) (items []voucher.Vulnerability, err error) {
	filterStr := kindFilterStr(reference, common.NoteKind_VULNERABILITY)
	err = pollForDiscoveries(ctx, g, reference)
	if nil != err {
		return []voucher.Vulnerability{}, err
	}

	project := projectPath(g.imageProject)
	req := &grafeaspb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	occIterator := g.grafeas.ListOccurrences(ctx, req)
	for {
		var occ *grafeaspb.Occurrence
		occ, err = occIterator.Next()
		if nil != err {
			if iterator.Done == err {
				err = nil
			}
			break
		}
		item := OccurrenceToVulnerability(occ)
		items = append(items, item)
	}

	return
}

// Close closes the Grafeas client.
func (g *Client) Close() {
	g.grafeas.Close()
}

// GetBuildDetails gets BuildDetails for the passed image.
func (g *Client) GetBuildDetails(ctx context.Context, reference reference.Canonical) (items []repository.BuildDetail, err error) {
	filterStr := kindFilterStr(reference, common.NoteKind_BUILD)

	project := projectPath(g.imageProject)
	req := &grafeaspb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	occIterator := g.grafeas.ListOccurrences(ctx, req)
	for {
		var occ *grafeaspb.Occurrence
		occ, err = occIterator.Next()
		if nil != err {
			if iterator.Done == err {
				err = nil
			}
			break
		}
		item := OccurrenceToBuildDetails(occ)
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

// NewClient creates a new Grafeas Client.
func NewClient(ctx context.Context, imageProject, binauthProject string, keyring *voucher.KeyRing) (*Client, error) {
	var err error

	client := new(Client)
	client.grafeas, err = containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}

	client.keyring = keyring
	client.binauthProject = binauthProject
	client.imageProject = imageProject

	return client, nil
}
