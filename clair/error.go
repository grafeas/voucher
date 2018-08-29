package clair

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/clair/api/v1"
)

func getErrorFromResponse(resp *http.Response) error {
	var errorStruct struct {
		Error v1.Error
	}

	errMsg := "(no message returned)"

	// attempt to unmarshal the body into a JSON structure, ignoring any errors if they occur.
	err := json.NewDecoder(resp.Body).Decode(&errorStruct)
	if nil != err {
		if "" != errorStruct.Error.Message {
			errMsg = errorStruct.Error.Message
		}
	}

	return fmt.Errorf("status code %d, %s", resp.StatusCode, errMsg)
}
