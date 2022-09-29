package docker

import (
	"io"
	"net/http"

	"github.com/docker/distribution"
)

// getDockerManifest executes an API call to Docker using the passed http.Client, and unmarshals
// the resulting data into the passed interface, or returns an error if there's an issue.
func getDockerManifest(client *http.Client, request *http.Request) (distribution.Manifest, error) {
	resp, err := client.Do(request)
	if nil != err {
		return nil, NewManifestError(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if nil != err {
		return nil, NewManifestError(err)
	}

	if resp.StatusCode >= 300 {
		return nil, NewManifestErrorWithRequest(resp.Status, b)
	}

	contentType := resp.Header.Get("Content-Type")
	if !isValidManifest(contentType) {
		return nil, NewManifestErrorWithRequest(resp.Status, b)
	}

	manifest, _, err := distribution.UnmarshalManifest(contentType, b)
	if nil != err {
		return nil, NewManifestError(err)
	}

	return manifest, nil
}

// isValidManifest ensures that we don't try to unmarshal an invalid manifest.
func isValidManifest(contentType string) bool {
	for _, mediaType := range distribution.ManifestMediaTypes() {
		if mediaType == contentType {
			return true
		}
	}

	return false
}
