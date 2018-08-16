package docker

import (
	"net/http"

	reference "github.com/docker/distribution/reference"
)

// Auth uses a gcloudToken to authorize against GCR. Returns the OAuthToken
// GCR provides.
func Auth(gcloudToken string, repository reference.Named) (OAuthToken, error) {
	var oauthToken OAuthToken

	var err error

	request, err := http.NewRequest(http.MethodGet, GetTokenURI(repository), nil)
	if nil != err {
		return oauthToken, err
	}

	request.SetBasicAuth("oauth2accesstoken", gcloudToken)

	err = doDockerCall(request, &oauthToken)

	return oauthToken, err
}
