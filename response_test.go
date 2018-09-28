package voucher

import (
	"bytes"
	"encoding/json"
	"testing"
)

const responseJson = `{
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

	buf := bytes.NewBufferString(responseJson)

	err := json.NewDecoder(buf).Decode(&response)

	if nil != err {
		t.Fatalf("failed to unmarshal valid data: %s", err)
	}
}
