package main

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher/docker"
)

// getCanonicalReference gets the canonical image reference for the passed
// image reference. If the passed reference is not already a canonical image
// reference, this method will connect to the registry to get the current digest
// and create the canonical reference from the original reference and that digest.
//
// This is because Binary Authorization only supports canonical image references,
// as a non-canonical image reference could refer to multiple versions of the same
// image (with different contents).
func getCanonicalReference(client *http.Client, ref reference.Reference) (reference.Canonical, error) {
	if canonicalRef, ok := ref.(reference.Canonical); ok {
		return canonicalRef, nil
	}

	if taggedRef, ok := ref.(reference.NamedTagged); ok {
		imageDigest, err := docker.GetDigestFromTagged(client, taggedRef)
		if nil != err {
			return nil, fmt.Errorf("getting digest from tag failed: %s", err)
		}
		canonicalRef, err := reference.WithDigest(reference.TrimNamed(taggedRef), imageDigest)
		if nil != err {
			return nil, fmt.Errorf("making canonical reference failed: %s", err)
		}
		return canonicalRef, nil
	}
	return nil, fmt.Errorf("reference cannot be converted to a canonical reference")
}
