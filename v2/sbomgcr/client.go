package sbomgcr

import (
	gcr "github.com/google/go-containerregistry/pkg/v1/google"
)

func NewClient() {
	gcr.NewEnvAuthenticator()
}
