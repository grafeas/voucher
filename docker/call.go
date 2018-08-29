package docker

import (
	"encoding/json"
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

	return json.NewDecoder(resp.Body).Decode(&data)
}
