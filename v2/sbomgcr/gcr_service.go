package sbomgcr

import (
	"context"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	gcr "github.com/google/go-containerregistry/pkg/v1/google"
)

// GCRService is an interface for the GCR service
type GCRService interface {
	PullImage(src string) (v1.Image, error)
	ListTags(ctx context.Context, repo name.Repository) (*gcr.Tags, error)
}

//  gcrServiceImpl implements GCRService
type gcrServiceImpl struct {
	auth authn.Authenticator
}

// NewGCRService returns a new GCRService
func NewGCRService() GCRService {
	auth, err := gcr.NewEnvAuthenticator()

	if err != nil {
		return nil
	}

	return &gcrServiceImpl{auth: auth}
}

// PullImage returns the image described by src
func (g *gcrServiceImpl) PullImage(src string) (v1.Image, error) {
	return crane.Pull(src)
}

// ListTags returns the tags for the repo
func (g *gcrServiceImpl) ListTags(ctx context.Context, repo name.Repository) (*gcr.Tags, error) {
	return gcr.List(repo, gcr.WithAuth(g.auth), gcr.WithContext(ctx))
}
