package subscriber

import (
	"encoding/json"
	"errors"

	"github.com/docker/distribution/reference"
)

const insertAction string = "INSERT"

var (
	errNotInsertAction  error = errors.New("ignoring; not an INSERT action")
	errNoDigest         error = errors.New("no digest specified")
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
		return nil, errNotInsertAction
	}

	// an image without a digest is only tagged; opt to skip this image
	// as the digest version has already been pushed and will be processed soon
	if pl.Digest == "" {
		return nil, errNoDigest
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
