package subscriber

import (
	"encoding/json"
	"errors"

	"github.com/docker/distribution/reference"
)

const insertAction string = "INSERT"

var (
	errInvalidPayload   error = errors.New("ignoring; payload is invalid")
	errConversionFailed error = errors.New("error converting reference into type reference.Canonical")
)

// Payload contains the information from the pub/sub message
type Payload struct {
	Action string `json:"action"`
	Digest string `json:"digest"`
	Tag    string `json:"tag,omitempty"`
}

// parsePayload parses data from a pubsub message.
func parsePayload(message []byte) (*Payload, error) {
	var pl Payload

	err := json.Unmarshal(message, &pl)
	if err != nil {
		return nil, err
	}

	if pl.Action != insertAction {
		return nil, errInvalidPayload
	}

	// an image without a digest is only tagged; opt to skip this image
	// as the digest version has already been pushed and will be procssed soon
	if pl.Digest == "" {
		return nil, errInvalidPayload
	}

	return &pl, nil
}

func (p *Payload) asCanonicalImage() (reference.Canonical, error) {
	imageRef, err := reference.Parse(p.Digest)
	if nil != err {
		return nil, err
	}

	canonicalRef, ok := imageRef.(reference.Canonical)
	if !ok {
		return nil, errConversionFailed
	}

	return canonicalRef, nil
}
