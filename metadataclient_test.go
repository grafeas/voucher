package voucher

import "errors"

type testMetadataClient struct {
	canAttest bool
	keyring   *KeyRing
}

func (t *testMetadataClient) CanAttest() bool {
	return t.canAttest
}

func (t *testMetadataClient) NewPayloadBody(i ImageData) (string, error) {
	if t.canAttest {
		return i.String(), nil
	}
	return "", errors.New("cannot create payload body")
}

func (t *testMetadataClient) GetMetadata(i ImageData, metadataType MetadataType) ([]MetadataItem, error) {
	return []MetadataItem{}, nil
}

func (t *testMetadataClient) AddAttestationToImage(i ImageData, payload AttestationPayload) (MetadataItem, error) {
	_, _, err := payload.Sign(t.keyring)
	if nil != err {
		return nil, err
	}

	occ := new(MetadataItem)

	return *occ, nil
}

func newTestMetadataClient(keyring *KeyRing, canAttest bool) *testMetadataClient {
	client := new(testMetadataClient)
	client.keyring = keyring
	client.canAttest = canAttest
	return client
}
