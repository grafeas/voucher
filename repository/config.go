package repository

// Config contains the necessary parameters to authenticate/communicate with a source repository
type Config struct {
	Auth            KeyRing
	Organization    string `json:"org-name"`
	OrganizationURL string `json:"org-url"`
}
