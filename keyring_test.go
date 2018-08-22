package voucher

import (
	"strconv"
	"testing"
)

const snakeoilKeyID = "1E92E2B4BB73E885"
const snakeoilKeyFingerprint = "90E942641C07A4C466BA97161E92E2B4BB73E885"
const testSignedValue = "test value to sign"

func TestGetKey(t *testing.T) {

	keyring, err := EjsonToKeyRing("tests/fixtures/key", "tests/fixtures/test.ejson")
	if nil != err {
		t.Fatalf("Failed to get keys from ejson: %s", err)
	}

	entity, err := keyring.GetSignerByName("snakeoil")
	if nil != err {
		t.Fatalf("Failed to get signing key from KeyRing: %s", err)
	}

	if nil == entity.PrimaryKey {
		t.Fatalf("Failed to get private key from KeyRing.")
	}

	keyID, err := strconv.ParseUint(snakeoilKeyID, 16, 64)
	if nil != err {
		t.Fatalf("Failed to convert snakeoilKeyID to uint: %s", err)
	}

	if entity.PrimaryKey.KeyId != keyID {
		t.Fatalf("Failed to get same key ID from KeyRing: \n%d vs \n%d", entity.PrimaryKey.KeyId, keyID)
	}

	signedValue, err := Sign(entity, testSignedValue)
	if nil != err {
		t.Fatalf("Failed to sign message: %s", err)
	}

	_, err = Verify(keyring, signedValue)
	if nil != err {
		t.Fatalf("Failed to verify signed message: %s", err)
	}
}
