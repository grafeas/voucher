package sbomgcr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/docker/distribution/reference"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	gcr "github.com/google/go-containerregistry/pkg/v1/google"
	voucher "github.com/grafeas/voucher/v2"
)

// structs for unmarshalling

type Envelope struct {
	PayloadType string      `json:"payloadType"`
	Payload     string      `json:"payload"`
	Signatures  []Signature `json:"signatures"`
}

type Signature struct {
	KeyID string `json:"keyid"`
	Sig   string `json:"sig"`
}

type CustomPredicate struct {
	Type          string `json:"_type"`
	PredicateType string `json:"predicateType"`
	Subject       []struct {
		Name   string `json:"name"`
		Digest struct {
			Sha256 string `json:"sha256"`
		} `json:"digest"`
	} `json:"subject"`
	Predicate struct {
		Data      string    `json:"Data"`
		Timestamp time.Time `json:"Timestamp"`
	} `json:"predicate"`
}

// Client connects to GCR
type Client struct{}

// GetVulnerabilities returns the detected vulnerabilities for the Image described by voucher.ImageData.
func (c *Client) GetVulnerabilities(ctx context.Context, ref reference.Canonical) (vulnerabilities []voucher.Vulnerability, err error) {
	return []voucher.Vulnerability{}, nil
}

// GetSBOM gets the SBOM for the passed image.
func (c *Client) GetSBOM(ctx context.Context, ref reference.Canonical) (cyclonedx.BOM, error) {
	repoName := ref.Name()
	imageSHA := string(ref.Digest())

	tag := strings.Replace(imageSHA, ":", "-", 1) + ".att"

	sbomDigest, err := c.GetSBOMDigestWithTag(ctx, repoName, tag)

	if err != nil {
		return cyclonedx.BOM{}, fmt.Errorf("error getting sbom digest with tag %w", err)
	}

	sbomImageName := repoName + "@" + sbomDigest

	sbom, err := crane.Pull(sbomImageName)

	if err != nil {
		return cyclonedx.BOM{}, fmt.Errorf("error pulling image from gcr with crane %w", err)
	}

	cycloneDX, err := c.GetSbomFromImage(sbom)

	if err != nil {
		return cyclonedx.BOM{}, fmt.Errorf("error getting SBOM from image %w", err)
	}

	return cycloneDX, nil
}

// GetSBOMDigestWithTag gets the gcr tags for the passed image.
func (c *Client) GetSBOMDigestWithTag(ctx context.Context, repoName string, tag string) (string, error) {
	repository, err := name.NewRepository(repoName)

	if err != nil {
		return "", fmt.Errorf("error returning repo name %w", err)
	}

	auth, err := gcr.NewEnvAuthenticator()

	if err != nil {
		return "", fmt.Errorf("error returning auth %w", err)
	}

	tags, err := gcr.List(repository, gcr.WithAuth(auth), gcr.WithContext(ctx))

	if err != nil {
		return "", fmt.Errorf("error occurred when trying to retrieve tags from this repo %w", err)
	}

	for digest, manifest := range tags.Manifests {
		for _, t := range manifest.Tags {
			if tag == t {
				return digest, nil
			}
		}
	}

	return "", fmt.Errorf("no digest found in Client.GetSBOMDigestWithTag")
}

func (c *Client) GetSbomFromImage(image v1.Image) (cyclonedx.BOM, error) {
	var cyclonedxBOM cyclonedx.BOM
	layer, err := image.Layers()

	if err != nil {
		return cyclonedxBOM, fmt.Errorf("error getting layers from image %w", err)
	}

	readCloser, _ := layer[0].Uncompressed()

	envelope, err := getEnvelopeFromReader(readCloser)

	if err != nil {
		return cyclonedxBOM, fmt.Errorf("error getting envelope %w", err)
	}

	customPredicate, err := getCustomPredicateFromEnvelope(envelope)

	if err != nil {
		return cyclonedxBOM, fmt.Errorf("error getting custom predicate %w", err)
	}

	err = json.Unmarshal([]byte(customPredicate.Predicate.Data), &cyclonedxBOM)

	if err != nil {
		return cyclonedxBOM, fmt.Errorf("error unmarshalling into cycloneDX SBOM %w", err)
	}

	return cyclonedxBOM, nil
}

func getEnvelopeFromReader(reader io.ReadCloser) (Envelope, error) {
	bt, _ := io.ReadAll(reader)
	var envelope Envelope

	err := json.Unmarshal(bt, &envelope)

	if err != nil {
		return envelope, fmt.Errorf("error unmarshalling into envelope %w", err)
	}

	return envelope, nil
}

func getCustomPredicateFromEnvelope(envelope Envelope) (CustomPredicate, error) {
	decoded, _ := base64.StdEncoding.DecodeString(string(envelope.Payload))
	var predicate CustomPredicate

	err := json.Unmarshal(decoded, &predicate)

	if err != nil {
		return predicate, fmt.Errorf("error unmarshalling into custom predicate %w", err)
	}

	return predicate, nil
}

// NewClient creates a new sbomgcr
func NewClient() *Client {
	client := new(Client)
	return client
}
