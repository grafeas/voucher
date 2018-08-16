package voucher

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var errNotSigned = errors.New("contents were not signed")
var errNoSigner = errors.New("signer is not in keyring")

// signConfig is used for our Signer.
var signConfig = packet.Config{
	DefaultHash:            crypto.SHA512,
	DefaultCipher:          packet.CipherAES256,
	DefaultCompressionAlgo: packet.CompressionZLIB,
	CompressionConfig: &packet.CompressionConfig{
		Level: 9,
	},
	RSABits: 4096,
}

// Sign creates the signature for the attestation
func Sign(signer *openpgp.Entity, msg string) (string, error) {
	buf := new(bytes.Buffer)

	armor, err := armor.Encode(buf, openpgp.SignatureType, make(map[string]string))
	defer armor.Close()
	if err != nil {
		return "", fmt.Errorf("creating armor writer failed: %s", err)
	}

	signature, err := openpgp.Sign(armor, signer, nil, &signConfig)
	defer signature.Close()
	if nil != err {
		return "", err
	}

	_, err = signature.Write([]byte(msg))
	if nil != err {
		return "", fmt.Errorf("writing to signature writer failed: %s", err)
	}

	if cerr := signature.Close(); nil != cerr {
		return "", cerr
	}

	if cerr := armor.Close(); nil != cerr {
		return "", cerr
	}

	return buf.String(), nil
}

// Verify verifies a signed message's signature, and returns the message
// that was signed as well as an error if applicable.
func Verify(keyring openpgp.KeyRing, signed string) (string, error) {
	armoredBlock, err := armor.Decode(bytes.NewBufferString(signed))
	if nil != err {
		return "", fmt.Errorf("could not decode armor: %s", err)
	}

	messageDetails, err := openpgp.ReadMessage(armoredBlock.Body, keyring, nil, &signConfig)
	if nil != err {
		return "", err
	}

	if !messageDetails.IsSigned {
		return "", errNotSigned
	}

	if nil == messageDetails.SignedBy {
		return "", errNoSigner
	}

	body, err := ioutil.ReadAll(messageDetails.UnverifiedBody)
	if nil != err {
		if nil != messageDetails.SignatureError {
			err = messageDetails.SignatureError
		}
	}
	return string(body), err
}
