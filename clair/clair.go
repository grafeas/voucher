package clair

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Shopify/voucher/docker"
	"github.com/coreos/clair/api/v1"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
)

// sendLayerToClair sends a layer from the passed repository (with the passed LayerReference).
func sendLayerToClair(hostname string, token docker.OAuthToken, layerRef LayerReference) (err error) {
	layer := AddAuthorization(layerRef.GetLayer(), token.Token)

	data := map[string]v1.Layer{
		"Layer": layer,
	}

	var buffer bytes.Buffer

	err = json.NewEncoder(&buffer).Encode(data)
	if nil != err {
		return
	}

	request, err := http.NewRequest(http.MethodPost, "http://"+hostname+"/v1/layers", &buffer)
	if nil != err {
		return
	}

	resp, err := http.DefaultClient.Do(request)
	if nil != err {
		return
	}

	err = resp.Body.Close()

	return
}

// getLayerFromClair gets the description of the Layer with the passed digest from Clair,
// using the passed digest.
func getLayerFromClair(hostname string, token docker.OAuthToken, digest digest.Digest) (layer v1.Layer, err error) {
	url := "http://" + hostname + "/v1/layers/" + string(digest) + "?vulnerabilities"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if nil != err {
		return
	}

	resp, err := http.DefaultClient.Do(request)
	if nil != err {
		return
	}

	defer resp.Body.Close()

	var layerEnv v1.LayerEnvelope

	err = json.NewDecoder(resp.Body).Decode(&layerEnv)
	if nil != layerEnv.Layer {
		layer = *layerEnv.Layer
	}

	return
}

// getVulnerabilities gets a map[string]v1.Vulnerability from Clair, so that we can convert
// them to Voucher Vulnerabilities all at once.
func getVulnerabilities(hostname string, oauthToken docker.OAuthToken, image reference.Canonical) (map[string]v1.Vulnerability, error) {
	vulns := make(map[string]v1.Vulnerability)

	manifest, err := docker.RequestManifest(oauthToken, image)
	if nil != err {
		return vulns, err
	}

	parent := digest.Digest("")

	for _, imageLayer := range manifest.Layers {
		var layer v1.Layer

		current := imageLayer.Digest

		ref := LayerReference{
			Image:   image,
			Current: current,
			Parent:  parent,
		}

		if err = sendLayerToClair(hostname, oauthToken, ref); nil != err {
			return vulns, err
		}

		if layer, err = getLayerFromClair(hostname, oauthToken, current); nil != err {
			return vulns, err
		}

		for _, feature := range layer.Features {
			for _, vul := range feature.Vulnerabilities {
				vulns[vul.Name] = vul
			}
		}

		parent = current
	}
	return vulns, err
}
