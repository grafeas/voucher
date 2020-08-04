package clair

import (
	v1 "github.com/coreos/clair/api/v1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
	"golang.org/x/oauth2"
)

func getSchema2Layers(m schema2.Manifest, config Config, tokenSrc oauth2.TokenSource, image reference.Canonical, parent digest.Digest) (map[string]v1.Vulnerability, error) {
	vulns := make(map[string]v1.Vulnerability)

	var err error

	for _, imageLayer := range m.Layers {
		current := imageLayer.Digest
		// send the current layer to Clair
		if err = sendLayerToClair(config, tokenSrc, NewLayerReference(image, current, parent)); nil != err {
			return vulns, err
		}

		parent = current
	}

	vulns, err = checkParentDigest(parent, config, vulns)
	if err != nil {
		return vulns, err
	}

	return vulns, nil
}
