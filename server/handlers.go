package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleCheckImage is a request handler that executes an individual Check or
// all of the Checks in one CheckGroup and creates any attestations if
// applicable.
func (s *Server) HandleCheckImage(w http.ResponseWriter, r *http.Request) {
	var err error

	if err = s.isAuthorized(r); nil != err {
		http.Error(w, "username or password is incorrect", http.StatusUnauthorized)
		LogError("username or password is incorrect", err)
		return
	}

	variables := mux.Vars(r)
	checkName := variables["check"]

	if "" == checkName {
		http.Error(w, "failure", http.StatusInternalServerError)
		return
	}

	requiredChecks := []string{checkName}

	if s.HasCheckGroup(checkName) {
		requiredChecks = s.GetCheckGroup(checkName)
	}

	if err = verifiedRequiredChecksAreRegistered(requiredChecks...); err != nil {
		http.Error(w, fmt.Sprintf("check or group \"%s\" is not active: %s", checkName, err), http.StatusNotFound)
		return
	}

	s.handleChecks(w, r, requiredChecks...)
}

// HandleVerifyImage is a request handler that verifies an individual
// attestation or all of the attestations which would be created by one
// CheckGroup and creates any attestations if applicable.
func (s *Server) HandleVerifyImage(w http.ResponseWriter, r *http.Request) {
	var err error

	if err = s.isAuthorized(r); nil != err {
		http.Error(w, "username or password is incorrect", 401)
		LogError("username or password is incorrect", err)
		return
	}

	variables := mux.Vars(r)
	checkName := variables["check"]

	if "" == checkName {
		http.Error(w, "failure", http.StatusInternalServerError)
		return
	}

	requiredChecks := []string{checkName}

	if s.HasCheckGroup(checkName) {
		requiredChecks = s.GetCheckGroup(checkName)
	}

	if err = verifiedRequiredChecksAreRegistered(requiredChecks...); err != nil {
		http.Error(w, fmt.Sprintf("check or group \"%s\" is not active: %s", checkName, err), http.StatusNotFound)
		return
	}

	s.handleVerify(w, r, requiredChecks...)
}

// HandleHealthCheck is a request handler that returns HTTP Status Code 200
// when it is called. Can be used to determine uptime.
func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
