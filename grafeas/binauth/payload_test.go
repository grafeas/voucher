package binauth

import (
	"encoding/json"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/assert"
)

const testPayloadURL = "gcr.io/test/image/we/are/testing@sha256:8c733d1c02c464f484bccd5fda4bc560040b80e105e431f17e1b1c3fce9b7d27"

const testPayloadOutput = `{
  "critical": {
    "identity": {
      "docker-reference": "gcr.io/test/image/we/are/testing"
    },
    "image": {
      "docker-manifest-digest": "sha256:8c733d1c02c464f484bccd5fda4bc560040b80e105e431f17e1b1c3fce9b7d27"
    },
    "type": "Google cloud binauthz container signature"
  }
}`

func TestNewPayload(t *testing.T) {
	assert := assert.New(t)

	rawRef, err := reference.Parse(testPayloadURL)
	assert.Nil(err)

	canonicalRef, isCanonical := rawRef.(reference.Canonical)
	assert.True(isCanonical)

	payload := NewPayload(canonicalRef)

	b, err := json.MarshalIndent(&payload, "", "  ")
	assert.Nil(err)

	assert.Equal(testPayloadOutput, string(b))
}
