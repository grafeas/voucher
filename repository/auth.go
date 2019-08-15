package repository

// KeyRing contains all the authentication keys
// needed to communicate an org's repository source
type KeyRing map[string]Auth

// Auth holds the necessary information to connect to a repository source
// It's exclusively a Token-based auth or Basic (user-password) auth
type Auth struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Type determines the authentication method being used to connect to a source
func (a *Auth) Type() string {
	if a.Token != "" {
		return TokenAuthType
	}

	if a.Username != "" && a.Password != "" {
		return UserPasswordAuthType
	}

	return ""
}
