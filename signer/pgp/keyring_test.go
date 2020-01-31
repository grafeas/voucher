package pgp

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const snakeoilKeyID = "1E92E2B4BB73E885"
const snakeoilKeyFingerprint = "90E942641C07A4C466BA97161E92E2B4BB73E885"
const testSignedValue = "test value to sign"

func newTestKeyRing(t *testing.T) *KeyRing {
	t.Helper()

	require := require.New(t)

	keyring := NewKeyRing()

	keyFile, err := os.Open("../../testdata/testkey.asc")
	require.NoErrorf(err, "failed to open key file: %s", err)
	defer keyFile.Close()

	err = AddKeyToKeyRingFromReader(keyring, "snakeoil", keyFile)
	require.NoErrorf(err, "Failed to add key to keyring: %s", err)

	return keyring
}

func getTestKeyID(t *testing.T) uint64 {
	t.Helper()

	require := require.New(t)

	keyID, err := strconv.ParseUint(snakeoilKeyID, 16, 64)
	require.NoErrorf(err, "Failed to convert snakeoilKeyID to uint: %s", err)

	return keyID
}

func TestGetKeyAndSign(t *testing.T) {
	require := require.New(t)

	keyring := newTestKeyRing(t)

	entity, err := keyring.GetSignerByName("snakeoil")
	require.NoErrorf(err, "Failed to get signing key from KeyRing: %s", err)
	require.NotNilf(entity.PrimaryKey, "Failed to get private key from KeyRing.")

	keyID := getTestKeyID(t)
	require.Equal(entity.PrimaryKey.KeyId, keyID)

	signedValue, err := sign(entity, testSignedValue)
	require.NoErrorf(err, "Failed to sign message: %s", err)

	_, err = Verify(keyring, signedValue)
	require.NoErrorf(err, "Failed to verify signed message: %s", err)
}

func TestOpenpgpKeyRing(t *testing.T) {
	assert := assert.New(t)

	keyring := newTestKeyRing(t)

	keyID := getTestKeyID(t)

	keys := keyring.KeysById(keyID)

	assert.Lenf(keys, 1, "incorrect number of keys returned by KeysByID")

	for _, key := range keys {
		assert.Equal(key.PublicKey.KeyId, keyID, "returned key that shouldn't have been, key ID is %X, should be %s", key.PublicKey.Fingerprint, snakeoilKeyFingerprint)
	}

	encKeys := keyring.DecryptionKeys()
	assert.Lenf(encKeys, 0, "too many keys returned by DecryptionKeys")
}

func TestBadAddKey(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.Buffer{}

	keyring := newTestKeyRing(t)

	err := AddKeyToKeyRingFromReader(keyring, "badkey", &buffer)
	if assert.Error(err) {
		assert.Equal(err.Error(), "openpgp: invalid argument: no armored data found")
	}
}

func TestAttestationPayload(t *testing.T) {
	payloadMessage := "test was successful"

	keyring := newTestKeyRing(t)

	result, fingerprint, err := keyring.Sign("snakeoil", payloadMessage)
	if assert.NoErrorf(t, err, "Failed to sign attestation: %s", err) {
		assert.Equalf(t, snakeoilKeyFingerprint, fingerprint, "Failed to get correct fingerprint, was %s vs %s", fingerprint, snakeoilKeyFingerprint)
	}

	message, err := Verify(keyring, result)
	if assert.NoErrorf(t, err, "Failed to verify result: %s", result) {
		assert.Equalf(t, message, payloadMessage, "Failed to get correct message, was \"%s\" instead of \"%s\"", message, payloadMessage)
	}
}
