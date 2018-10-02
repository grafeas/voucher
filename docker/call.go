package docker

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// doDockerCall executes an API call to Docker using the passed http.Client, and unmarshals
// the resulting data into the passed interface, or returns an error if there's an issue.
func doDockerCall(client *http.Client, request *http.Request, data interface{}) error {
	resp, err := client.Do(request)
	if nil != err {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return responseToError(resp)
	}

	return json.NewDecoder(resp.Body).Decode(&data)
}

// responseToError converts the body of a response to an error.
func responseToError(resp *http.Response) error {
	b, _ := ioutil.ReadAll(resp.Body)
	return errors.New("failed to load resource with status \"" + resp.Status + "\": " + string(b))
}
