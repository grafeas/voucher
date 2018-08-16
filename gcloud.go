package voucher

import (
	"strings"
)

// GetAccessToken queries the gcloud command to get an access token.
// This token is then passed to the API to get a bearer token.
func GetAccessToken() (string, error) {
	result, err := RunShellCommand("gcloud", "auth", "print-access-token")
	if nil != err {
		return "", err
	}

	return strings.TrimSpace(result), err
}
