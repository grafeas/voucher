package provenance

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/voucher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

// check holds the required data for the check
type check struct {
	metadataClient voucher.MetadataClient
}

// SetMetadataClient sets the MetadataClient for this Check.
func (p *check) SetMetadataClient(metadataClient voucher.MetadataClient) {
	p.metadataClient = metadataClient
}

// Check runs the check :)
func (p *check) Check(i voucher.ImageData) (bool, error) {
	occurrences, err := p.metadataClient.GetMetadata(i, voucher.BuildDetailsType)
	if err != nil {
		return false, err
	}

	// there should only be one occurrence
	if len(occurrences) != 1 {
		return false, fmt.Errorf("Got %d occurrences for: %s", len(occurrences), i.String())
	}

	buildDetails := occurrences[0].GetBuildDetails()
	if validateProvenance(buildDetails) && validateArtifacts(i, buildDetails) {
		log.Infof("Validated image provenance and artifacts for: %s", i.String())
		return true, nil
	}

	return false, nil
}

func validateProvenance(details *containeranalysispb.BuildDetails) (trusted bool) {
	// get trusted things
	trustedBuilderIdentities := voucher.ToMapStringBool(viper.GetStringMap("trusted-builder-identities"))
	trustedBuilderProjects := voucher.ToMapStringBool(viper.GetStringMap("trusted-projects"))

	provenance, err := base64.StdEncoding.DecodeString(details.ProvenanceBytes)
	if err != nil {
		log.Errorf("Error decoding provenance: %v", err)
		return
	}

	// unmarshal against provenance bytes
	if err := json.Unmarshal(provenance, details); err != nil {
		log.Errorf("Provenance bytes do not match details: %v", err)
		return
	}

	if !trustedBuilderIdentities[details.Provenance.Creator] {
		log.Errorf("Builder identity not trusted: %s", details.Provenance.Creator)
		return
	}

	if !trustedBuilderProjects[details.Provenance.ProjectId] {
		log.Errorf("Builder project not trusted: %s", details.Provenance.ProjectId)
		return
	}

	trusted = true
	return
}

func validateArtifacts(i voucher.ImageData, details *containeranalysispb.BuildDetails) (matched bool) {
	// if an artifact built by this Build is the image, validate the SHAs match
	for _, artifact := range details.Provenance.BuiltArtifacts {
		if strings.HasSuffix(i.Digest().String(), artifact.Checksum) {
			matched = true
		}
	}
	return
}

func init() {
	voucher.RegisterCheckFactory("provenance", func() voucher.Check {
		return new(check)
	})
}
