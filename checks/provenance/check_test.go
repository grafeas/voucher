package provenance

import (
	"testing"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	builderIdentityTestData = "trusted-person@email.com"
	imageSHA256TestData     = "sha256:1234c923e00e0fd2ba78041bfb64a105e1ecb7678916d1f7776311e45bf57890"
	imageURLTestData        = "gcr.io/" + projectTestData + "/name@" + imageSHA256TestData
	projectTestData         = "test"
)

var buildDetailsTestData = repository.BuildDetail{
	ProjectID:    projectTestData,
	BuildCreator: builderIdentityTestData,
	Artifacts: []repository.BuildArtifact{
		{
			ID:       imageURLTestData,
			Checksum: imageSHA256TestData,
		},
	},
}

func TestArtifactIsImage(t *testing.T) {
	imageDataTestData, err := voucher.NewImageData(imageURLTestData)
	require.NoError(t, err)

	assert := assert.New(t)
	result := validateArtifacts(imageDataTestData, buildDetailsTestData)
	assert.True(result)
}
