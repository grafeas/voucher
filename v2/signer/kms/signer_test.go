package kms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/googleapis/gax-go/v2"
	"github.com/grafeas/voucher/v2/signer/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kms_pb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	checkName     = "my-check"
	mockSignature = "signed"
	keyPath       = "project/my-project/locations/global/keyRings/my-keyring/cryptoKeys/my-key"
)

func TestSigner_Sign(t *testing.T) {
	const checkBody = "pass"
	cases := []struct {
		algo     string
		expected string
	}{
		{
			algo:     kms.AlgoSHA256,
			expected: "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1",
		},
		{
			algo:     kms.AlgoSHA384,
			expected: "8517aca50b29dbeb9f5d8964a8adf64e2f6592b0aff41eade47f68ba9a3849254fac883bf15836da007086e73a145d7b",
		},
		{
			algo:     kms.AlgoSHA512,
			expected: "5510ebbda5ed4da007c55a62fd7075c722ec031f07398ef3e90b9b50e0fe950985476c474414d2b386e8f08cd505fb506b528006a30abfe9ca0eb0b67b7e760b",
		},
	}

	for _, tc := range cases {
		t.Run(tc.algo, func(t *testing.T) {
			keys := map[string]kms.Key{
				checkName: {
					Path: keyPath,
					Algo: tc.algo,
				},
			}

			// Build signer with mocked KMS:
			k := &mockKMS{}
			signer, err := kms.NewSigner(keys, kms.WithKMSClient(k))
			require.NoError(t, err)

			sig, path, err := signer.Sign("my-check", checkBody)
			require.NoError(t, err)

			// Verify request to kMS, and interpretation of response
			if assert.Len(t, k.reqs, 1) {
				kmsReq := k.reqs[0]
				assert.Equal(t, "", fmt.Sprintf("%x", kmsReq.GetDigest().GetSha256()))
			}
			assert.Equal(t, mockSignature, sig)
			assert.Equal(t, "//cloudkms.googleapis.com/v1/project/my-project/locations/global/keyRings/my-keyring/cryptoKeys/my-key", path)
		})
	}
}

type mockKMS struct {
	reqs []*kms_pb.AsymmetricSignRequest
}

func (k *mockKMS) AsymmetricSign(_ context.Context, req *kms_pb.AsymmetricSignRequest, _ ...gax.CallOption) (*kms_pb.AsymmetricSignResponse, error) {
	k.reqs = append(k.reqs, req)
	return &kms_pb.AsymmetricSignResponse{
		Signature: []byte(mockSignature),
	}, nil
}
func (k *mockKMS) Close() error { return nil }
