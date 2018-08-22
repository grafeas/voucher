package voucher

import "testing"

func TestAttestationPayload(t *testing.T) {
	payloadMessage := "test was successful"

	keyring, err := EjsonToKeyRing("tests/fixtures/key", "tests/fixtures/test.ejson")
	if nil != err {
		t.Fatalf("Failed to get keys from ejson: %s", err)
	}

	payload := AttestationPayload{
		CheckName: "snakeoil",
		Body:      payloadMessage,
	}

	result, fingerprint, err := payload.Sign(keyring)
	if nil != err {
		t.Fatalf("Failed to sign attestation: %s", err)
	}

	if snakeoilKeyFingerprint != fingerprint {
		t.Fatalf("Failed to get correct fingerprint, was %s vs %s", fingerprint, snakeoilKeyFingerprint)
	}

	message, err := Verify(keyring, result)
	if nil != err {
		t.Fatalf("Failed to verify result: %s", result)
	}

	if message != payloadMessage {
		t.Fatalf("Failed to get correct message, was \"%s\" instead of \"%s\"", message, payloadMessage)
	}

}
