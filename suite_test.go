package voucher

import (
	"context"
	"os"
	"testing"

	"github.com/Shopify/voucher/signer/pgp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestKeyRing(t *testing.T) *pgp.KeyRing {
	t.Helper()

	require := require.New(t)

	keyring := pgp.NewKeyRing()

	keyFile, err := os.Open("testdata/testkey.asc")
	require.NoErrorf(err, "failed to open key file: %s", err)
	defer keyFile.Close()

	err = pgp.AddKeyToKeyRingFromReader(keyring, "snakeoil", keyFile)
	require.NoErrorf(err, "Failed to add key to keyring: %s", err)

	return keyring
}

func newTestImageData(t *testing.T) ImageData {
	t.Helper()
	imageData, err := NewImageData("localhost.local/path/to/image@sha256:b148c8af52ba402ed7dd98d73f5a41836ece508d1f4704b274562ac0c9b3b7da")
	assert.Nilf(t, err, "could not make ImageData")
	return imageData
}

func TestNewSuite(t *testing.T) {
	suite := NewSuite()
	require.NotNilf(t, suite, "could not make CheckSuite")

	imageData := newTestImageData(t)

	results := suite.Run(context.Background(), imageData)
	require.Equal(t, []CheckResult{}, results)

	brokenCheck := new(testBrokenCheck)

	// Add our three checks.
	suite.Add("passer", newTestCheck(true))
	suite.Add("failer", newTestCheck(false))
	suite.Add("broken", brokenCheck)

	expectedResults := []CheckResult{
		{
			Name:      "passer",
			ImageData: imageData,
			Err:       "",
			Success:   true,
			Attested:  false,
			Details:   nil,
		},
		{
			Name:      "failer",
			ImageData: imageData,
			Err:       "",
			Success:   false,
			Attested:  false,
			Details:   nil,
		},
		{
			Name:      "broken",
			ImageData: imageData,
			Err:       errBrokenTest.Error(),
			Success:   false,
			Attested:  false,
			Details:   nil,
		},
	}

	results = suite.Run(context.Background(), imageData)
	assert.ElementsMatch(t, expectedResults, results)

	fixedCheck, err := suite.Get("fixed")
	assert.Nil(t, fixedCheck)
	if assert.NotNil(t, err) {
		assert.Equal(t, err, ErrNoCheck)
	}

	gottenCheck, err := suite.Get("broken")
	assert.Nil(t, err)

	if assert.NotNil(t, gottenCheck) {
		assert.Equal(t, gottenCheck, brokenCheck)
	}
}

func TestMakeSuccessfulSuite(t *testing.T) {
	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	suite.Add("pass1", newTestCheck(true))
	suite.Add("pass2", newTestCheck(true))
	suite.Add("pass3", newTestCheck(true))

	imageData := newTestImageData(t)

	results := suite.Run(context.Background(), imageData)

	response := NewResponse(imageData, results)
	assert.Equal(t, true, response.Success)
}

func TestMakeFailingSuite(t *testing.T) {
	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	suite.Add("fail1", newTestCheck(false))
	suite.Add("fail2", newTestCheck(false))
	suite.Add("fail3", newTestCheck(false))

	imageData := newTestImageData(t)

	results := suite.Run(context.Background(), imageData)

	response := NewResponse(imageData, results)
	assert.Equal(t, false, response.Success)
}

func TestAttestSuite(t *testing.T) {
	keyring := newTestKeyRing(t)

	metadataClient := newTestMetadataClient(keyring, true)

	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	suite.Add("snakeoil", newTestCheck(true))
	suite.Add("pass2", newTestCheck(true))
	suite.Add("pass3", newTestCheck(true))

	imageData := newTestImageData(t)

	results := suite.RunAndAttest(context.Background(), metadataClient, imageData)

	expectedResults := []CheckResult{
		{
			Name:      "snakeoil",
			ImageData: imageData,
			Err:       "",
			Success:   true,
			Attested:  true,
			Details:   nil,
		},
		{
			Name:      "pass2",
			ImageData: imageData,
			Err:       "no signing entity exists for check name \"pass2\"",
			Success:   true,
			Attested:  false,
			Details:   nil,
		},
		{
			Name:      "pass3",
			ImageData: imageData,
			Err:       "no signing entity exists for check name \"pass3\"",
			Success:   true,
			Attested:  false,
			Details:   nil,
		},
	}

	assert.ElementsMatch(t, expectedResults, results)
}

func TestNonattestingSuite(t *testing.T) {
	keyring := newTestKeyRing(t)

	metadataClient := newTestMetadataClient(keyring, false)

	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	// only adding the snakeoil check, since that's the one we'll be attesting with
	suite.Add("snakeoil", newTestCheck(true))

	imageData := newTestImageData(t)

	results := suite.RunAndAttest(context.Background(), metadataClient, imageData)

	expectedResult := CheckResult{
		Name:      "snakeoil",
		ImageData: imageData,
		Err:       "cannot create payload body",
		Success:   true,
		Attested:  false,
		Details:   nil,
	}

	assert.Contains(t, results, expectedResult)
}
