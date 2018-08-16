package binauth

import (
	"encoding/json"

	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
)

const payloadType = "Google cloud binauthz container signature"

// PayloadIdentity represents the identity block in an Payload message.
type PayloadIdentity struct {
	DockerReference string `json:"docker-reference"`
}

// PayloadImage represents the image block in an Payload message.
type PayloadImage struct {
	DockerManifestDigest digest.Digest `json:"docker-manifest-digest"`
}

// PayloadCritical represents the critical block in the Payload message.
type PayloadCritical struct {
	Identity PayloadIdentity `json:"identity"`
	Image    PayloadImage    `json:"image"`
	Type     string          `json:"type"`
}

// Payload represents an Payload message.
type Payload struct {
	Critical PayloadCritical `json:"critical"`
}

// ToString returns the payload as a JSON encoded string, or returns an error.
func (p Payload) ToString() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// NewPayload creates a new Binauth specific payload for the image at
// the passed URL.
func NewPayload(reference reference.Canonical) Payload {
	payload := Payload{
		Critical: PayloadCritical{
			Identity: PayloadIdentity{
				DockerReference: reference.Name(),
			},
			Image: PayloadImage{
				DockerManifestDigest: reference.Digest(),
			},
			Type: payloadType,
		},
	}

	return payload
}
