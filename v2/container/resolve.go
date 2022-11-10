package container

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Resolver struct {
	auth authn.Authenticator
}

func NewResolver(auth authn.Authenticator) *Resolver {
	return &Resolver{
		auth: auth,
	}
}

func (r *Resolver) ToDigest(ctx context.Context, ref string) (*name.Digest, error) {
	// Parse reference, it might already contain a digest:
	parsedRef, err := name.ParseReference(ref, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("parsing reference failed: %w", err)
	}
	if digest, ok := parsedRef.(name.Digest); ok {
		return &digest, nil
	}

	// Fetch the remote ImageIndex/manifest list, and return its digest:
	opts := []remote.Option{remote.WithContext(ctx)}
	if r.auth != nil {
		opts = append(opts, remote.WithAuth(r.auth))
	}
	remoteImg, err := remote.Index(parsedRef, opts...)
	if err != nil {
		return nil, fmt.Errorf("getting remote image index failed: %w", err)
	}
	digest, err := remoteImg.Digest()
	if err != nil {
		return nil, fmt.Errorf("getting image digest failed: %w", err)
	}
	digestedRef, err := name.NewDigest(fmt.Sprintf("%s:%s@%s", parsedRef.Context().Name(), parsedRef.Identifier(), digest))
	if err != nil {
		return nil, fmt.Errorf("parsing digest failed: %w", err)
	}
	return &digestedRef, nil
}
