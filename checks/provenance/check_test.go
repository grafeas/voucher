package provenance

import (
	"testing"

	"github.com/Shopify/voucher"
	"github.com/stretchr/testify/assert"
	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

var (
	builderIdentityTestData = "trusted-person@email.com"
	imageSHA256TestData     = "sha256:1234c923e00e0fd2ba78041bfb64a105e1ecb7678916d1f7776311e45bf57890"
	imageURLTestData        = "gcr.io/" + projectTestData + "/name@" + imageSHA256TestData
	projectTestData         = "test"
)

var buildDetailsTestData = &containeranalysispb.BuildDetails{
	Provenance: &containeranalysispb.BuildProvenance{
		Id:        "foo",
		ProjectId: projectTestData,
		Creator:   builderIdentityTestData,
		BuiltArtifacts: []*containeranalysispb.Artifact{
			{
				Id:       imageURLTestData,
				Checksum: imageSHA256TestData,
				Names:    []string{"foo", "bar"},
			},
		},
	},
	ProvenanceBytes: "base64blob",
}

func TestArtifactIsImage(t *testing.T) {
	imageDataTestData, err := voucher.NewImageData(imageURLTestData)
	if err != nil {
		return
	}
	assert := assert.New(t)
	result := validateArtifacts(imageDataTestData, buildDetailsTestData)
	assert.True(result)
}
