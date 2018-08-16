package docker

// OAuthToken represents an OAuthToken returned by the
// Docker Registry API.
type OAuthToken struct {
	ExpiresIn int64  `json:"expires_in"`
	IssuedAt  string `json:"issued_at"`
	Token     string `json:"token"`
}
