package voucher

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const responseJSON = `{
	"image": "",
	"success": true,
	"results": [
		{
			"name": "diy",
			"success": true,
			"attested": true,
			"details": {
				"Name": "attested",
				"Details": "signature here"
			}

		}
	]
}`

func TestUnmarshalResponse(t *testing.T) {
	response := Response{}

	buf := bytes.NewBufferString(responseJSON)

	err := json.NewDecoder(buf).Decode(&response)
	assert.NoErrorf(t, err, "failed to unmarshal valid data: %s", err)
}
