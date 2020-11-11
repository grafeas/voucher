package clair

import (
	v1 "github.com/coreos/clair/api/v1"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
	"golang.org/x/oauth2"

	dockerURI "github.com/grafeas/voucher/docker/uri"
)

// LayerReference is a structure containing a Layer digest, as well as the repository
// URI, to simplify loading a Layer from the server.
type LayerReference struct {
	Image   reference.Canonical // The Image's reference.
	Current digest.Digest       // The digest of the current layer.
	Parent  digest.Digest       // The digest of the parent layer.
}

// GetURI gets the URI that is described in the LayerReference.
func (ref *LayerReference) GetURI() string {
	return dockerURI.GetBlobURI(ref.Image, ref.Current)
}

// GetLayer returns a layer description of the LayerReference.
func (ref *LayerReference) GetLayer() v1.Layer {
	return v1.Layer{
		Name:       string(ref.Current),
		Path:       ref.GetURI(),
		Headers:    make(map[string]string),
		ParentName: string(ref.Parent),
		Format:     "Docker",
	}
}

// NewLayerReference creates a new LayerReference based on the passed Image, and layer digest (the current digest)
// and that layer's parent digest.
func NewLayerReference(image reference.Canonical, current, parent digest.Digest) LayerReference {
	return LayerReference{
		Image:   image,
		Current: current,
		Parent:  parent,
	}
}

// AddAuthorization adds a Bearer token to the v1.Layer passed to it and
// returns a new v1.Layer.
func AddAuthorization(layer v1.Layer, token *oauth2.Token) v1.Layer {
	layer.Headers["Authorization"] = token.Type() + " " + token.AccessToken

	return layer
}
