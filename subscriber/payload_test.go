package subscriber

import "testing"

var payloadTestCases = []struct {
	name            string
	rawPayload      []byte
	expectedErr     error
	expectedPayload *Payload
}{
	{
		"successfully parses INSERT payload",
		[]byte(`{ "action":"INSERT", "digest": "some-digest-we-dont-validate" }`),
		nil,
		&Payload{Action: insertAction, Digest: "some-digest-dont-validate"},
	},
	{
		"returns nil payload with error for non-INSERT action payload",
		[]byte(`{ "action":"DELETE" }`),
		errInvalidPayload,
		nil,
	},
	{
		"returns an error when INSERT payload has no digest",
		[]byte(`{ "action":"INSERT" }`),
		errInvalidPayload,
		nil,
	},
}

func TestParsePayload(t *testing.T) {
	for _, tc := range payloadTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pl, err := parsePayload(tc.rawPayload)

			if (pl == nil && tc.expectedPayload != nil) || (pl != nil && tc.expectedPayload == nil) {
				t.Errorf("unexpected payload: got %v, want %v", pl, tc.expectedPayload)
			}

			if err != tc.expectedErr {
				t.Errorf("error: got %v, want %v", err, tc.expectedErr)
			}
		})
	}
}
