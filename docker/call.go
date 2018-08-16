package docker

import (
	"encoding/json"
	"net/http"
)

// setBearerToken sets the Bearer token in the Authorization header.
func setBearerToken(request *http.Request, token OAuthToken) {
	request.Header.Set("Authorization", "Bearer "+token.Token)
}

// doDockerCall executes an API call to Docker, and returns the resulting data.
// a schema2.Manifest, or an error if there's an issue.
func doDockerCall(request *http.Request, data interface{}) error {
	resp, err := http.DefaultClient.Do(request)
	if nil != err {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&data)
}
