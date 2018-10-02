package voucher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttestationPayload(t *testing.T) {
	payloadMessage := "test was successful"

	keyring := newTestKeyRing(t)

	payload := AttestationPayload{
		CheckName: "snakeoil",
		Body:      payloadMessage,
	}

	result, fingerprint, err := payload.Sign(keyring)
	if assert.NoErrorf(t, err, "Failed to sign attestation: %s", err) {
		assert.Equalf(t, snakeoilKeyFingerprint, fingerprint, "Failed to get correct fingerprint, was %s vs %s", fingerprint, snakeoilKeyFingerprint)
	}

	message, err := Verify(keyring, result)
	if assert.NoErrorf(t, err, "Failed to verify result: %s", result) {
		assert.Equalf(t, message, payloadMessage, "Failed to get correct message, was \"%s\" instead of \"%s\"", message, payloadMessage)
	}
}

func TestAttestationPayloadWithEmptyKeyRing(t *testing.T) {
	var keyring *KeyRing = nil

	payload := AttestationPayload{
		CheckName: "snakeoil",
		Body:      "this should fail",
	}

	// try to sign with
	_, _, err := payload.Sign(keyring)
	if assert.Error(t, err) {
		assert.Containsf(t, err.Error(), errEmptyKeyring.Error(), "Did not return correct error when signing with empty keyring: %s", err)
	}
}
