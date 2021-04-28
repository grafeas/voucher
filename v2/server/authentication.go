package server

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// isAuthorized returns true if the request's basic authentication header matches the
// configured username and password. The password in the configuration is assumed to
// have hashed using the bcrypt algorithm.
func (s *Server) isAuthorized(r *http.Request) error {
	// If the server does not require auth, the user is always authorized.
	if !s.serverConfig.RequireAuth {
		return nil
	}

	if s.serverConfig.Username == "" || s.serverConfig.PassHash == "" {
		return errors.New("username or password misconfigured in configuration")
	}

	username, password, ok := r.BasicAuth()
	if ok {
		if username == s.serverConfig.Username {
			if err := bcrypt.CompareHashAndPassword([]byte(s.serverConfig.PassHash), []byte(password)); nil != err {
				return err
			}
			return nil
		}
	}

	return errors.New("user failed to authenticate, username and/or password is incorrect")
}
