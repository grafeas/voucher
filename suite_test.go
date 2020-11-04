package voucher

import (
	"context"
	"errors"
	"testing"

	"github.com/grafeas/voucher/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

	results := suite.Run(context.Background(), &metrics.NoopClient{}, imageData)
	require.Equal(t, []CheckResult{}, results)

	errBrokenTest := errors.New("this test is broken")

	checks := make(map[string]Check)
	for _, check := range []struct {
		name string
		pass bool
		err  error
	}{
		{"passer", true, nil},
		{"failer", false, nil},
		{"broken", false, errBrokenTest},
	} {
		mockCheck := new(MockCheck)
		mockCheck.On("Check", mock.Anything, imageData).Return(check.pass, check.err)
		checks[check.name] = mockCheck
		suite.Add(check.name, mockCheck)
	}

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

	results = suite.Run(context.Background(), &metrics.NoopClient{}, imageData)
	assert.ElementsMatch(t, expectedResults, results)

	fixedCheck, err := suite.Get("fixed")
	assert.Nil(t, fixedCheck)
	if assert.NotNil(t, err) {
		assert.Equal(t, err, ErrNoCheck)
	}

	gottenCheck, err := suite.Get("broken")
	assert.Nil(t, err)

	if assert.NotNil(t, gottenCheck) {
		assert.Equal(t, checks["broken"], gottenCheck)
	}
}

func TestMakeSuccessfulSuite(t *testing.T) {
	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	imageData := newTestImageData(t)

	for _, name := range []string{"pass1", "pass2", "pass3"} {
		check := new(MockCheck)
		check.On("Check", mock.Anything, imageData).Return(true, nil)
		suite.Add(name, check)
	}

	results := suite.Run(context.Background(), &metrics.NoopClient{}, imageData)

	response := NewResponse(imageData, results)
	assert.Equal(t, true, response.Success)
}

func TestMakeFailingSuite(t *testing.T) {
	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	imageData := newTestImageData(t)

	for _, name := range []string{"fail1", "fail2", "fail3"} {
		check := new(MockCheck)
		check.On("Check", mock.Anything, imageData).Return(false, nil)
		suite.Add(name, check)
	}

	results := suite.Run(context.Background(), &metrics.NoopClient{}, imageData)

	response := NewResponse(imageData, results)
	assert.Equal(t, false, response.Success)
}

func TestAttestSuite(t *testing.T) {
	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	imageData := newTestImageData(t)
	errNoSigningEntity := errors.New("no signging entity exists for check")

	metadataClient := new(MockMetadataClient)
	metadataClient.
		On("NewPayloadBody", imageData).Return(imageData.String(), nil).
		On("AddAttestationToImage", mock.Anything, imageData, NewAttestation("snakeoil", imageData.String())).Return(SignedAttestation{
		Attestation: Attestation{
			CheckName: "snakeoil",
		},
	}, nil).
		On("AddAttestationToImage", mock.Anything, imageData, NewAttestation("pass2", imageData.String())).Return(SignedAttestation{
		Attestation: Attestation{
			CheckName: "pass2",
		},
	}, errNoSigningEntity).
		On("AddAttestationToImage", mock.Anything, imageData, NewAttestation("pass3", imageData.String())).Return(SignedAttestation{
		Attestation: Attestation{
			CheckName: "pass3",
		},
	}, errNoSigningEntity)

	for _, name := range []string{"snakeoil", "pass2", "pass3"} {
		check := new(MockCheck)
		check.On("Check", mock.Anything, imageData).Return(true, nil)
		suite.Add(name, check)
	}

	results := suite.RunAndAttest(context.Background(), metadataClient, &metrics.NoopClient{}, imageData)

	expectedResults := []CheckResult{
		{
			Name:      "snakeoil",
			ImageData: imageData,
			Err:       "",
			Success:   true,
			Attested:  true,
			Details: SignedAttestation{
				Attestation: Attestation{
					CheckName: "snakeoil",
				},
			},
		},
		{
			Name:      "pass2",
			ImageData: imageData,
			Err:       errNoSigningEntity.Error(),
			Success:   true,
			Attested:  false,
			Details: SignedAttestation{
				Attestation: Attestation{
					CheckName: "pass2",
				},
			},
		},
		{
			Name:      "pass3",
			ImageData: imageData,
			Err:       errNoSigningEntity.Error(),
			Success:   true,
			Attested:  false,
			Details: SignedAttestation{
				Attestation: Attestation{
					CheckName: "pass3",
				},
			},
		},
	}

	assert.ElementsMatch(t, expectedResults, results)
}

func TestNonattestingSuite(t *testing.T) {
	imageData := newTestImageData(t)
	errCreatingPayload := errors.New("cannot create payload body")

	metadataClient := new(MockMetadataClient)
	metadataClient.On("NewPayloadBody", imageData).Return("", errCreatingPayload)

	suite := NewSuite()
	assert.NotNilf(t, suite, "could not make CheckSuite")

	// only adding the snakeoil check, since that's the one we'll be attesting with
	check := new(MockCheck)
	check.On("Check", mock.Anything, imageData).Return(true, nil)
	suite.Add("snakeoil", check)

	results := suite.RunAndAttest(context.Background(), metadataClient, &metrics.NoopClient{}, imageData)

	expectedResult := CheckResult{
		Name:      "snakeoil",
		ImageData: imageData,
		Err:       errCreatingPayload.Error(),
		Success:   true,
		Attested:  false,
		Details:   nil,
	}

	assert.Contains(t, results, expectedResult)
}
