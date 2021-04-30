package vtesting

import (
	"os"
	"testing"

	"github.com/grafeas/voucher/v2/signer"
	"github.com/grafeas/voucher/v2/signer/pgp"

	"github.com/stretchr/testify/require"
)

// NewPGPSigner creates a new signer using the test key in the testdata
// directory.
func NewPGPSigner(t *testing.T) signer.AttestationSigner {
	t.Helper()

	newKeyRing := pgp.NewKeyRing()

	keyFile, err := os.Open("../../testdata/testkey.asc")
	require.NoError(t, err, "failed to open key file")
	defer keyFile.Close()

	err = pgp.AddKeyToKeyRingFromReader(newKeyRing, "snakeoil", keyFile)
	require.NoError(t, err, "failed to add key to keyring")

	return newKeyRing
}
