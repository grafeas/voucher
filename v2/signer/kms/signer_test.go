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
			expected: "ad111970708aaff07524d81f71582952b75ade74951ebb3c25e801fc4a3f17de8d3e8fbc7c271114462fe63f67d33536",
		},
		{
			algo:     kms.AlgoSHA512,
			expected: "5b722b307fce6c944905d132691d5e4a2214b7fe92b738920eb3fce3a90420a19511c3010a0e7712b054daef5b57bad59ecbd93b3280f210578f547f4aed4d25",
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

			// Build signer with mocked KMS, sign something:
			k := &mockKMS{}
			signer, err := kms.NewSigner(keys, kms.WithKMSClient(k))
			require.NoError(t, err)
			sig, path, err := signer.Sign("my-check", checkBody)
			require.NoError(t, err)

			// Verify request to kMS, and interpretation of response
			if assert.Len(t, k.reqs, 1) {
				kmsDigest := k.reqs[0].GetDigest()
				var digest string
				switch tc.algo {
				case kms.AlgoSHA256:
					digest = string(kmsDigest.GetSha256())
				case kms.AlgoSHA384:
					digest = string(kmsDigest.GetSha384())
				case kms.AlgoSHA512:
					digest = string(kmsDigest.GetSha512())
				default:
					require.FailNow(t, "unexpected algo")
				}
				assert.Equal(t, tc.expected, fmt.Sprintf("%x", digest))
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
