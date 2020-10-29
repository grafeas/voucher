package kms

import (
	"context"
	"crypto/sha512"
	"fmt"

	apiv1 "cloud.google.com/go/kms/apiv1"
	"github.com/grafeas/voucher/signer"
	kms_pb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	AlgoSHA512 = "SHA512"
	APIPath    = "//cloudkms.googleapis.com/v1"
)

type Key struct {
	Path string
	Algo string
}

// Signer is an AttestationSigner that uses Google's Cloud KMS to sign attestations
// Only supports SHA512 digests.
type Signer struct {
	keys   map[string]Key
	client *apiv1.KeyManagementClient
}

func NewSigner(keys map[string]Key) (*Signer, error) {
	client, err := apiv1.NewKeyManagementClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Signer{
		keys:   keys,
		client: client,
	}, nil
}

func (s *Signer) Sign(checkName, body string) (string, string, error) {
	key, ok := s.keys[checkName]
	if !ok {
		return "", "", signer.ErrNoKeyForCheck
	}

	if key.Algo != AlgoSHA512 {
		return "", "", fmt.Errorf("unable to hash algorithm %q, must be SHA512", key.Algo)
	}

	resp, err := s.client.AsymmetricSign(context.Background(), &kms_pb.AsymmetricSignRequest{
		Name: key.Path,
		Digest: &kms_pb.Digest{
			Digest: &kms_pb.Digest_Sha512{Sha512: sha512Digest([]byte(body))},
		},
	})

	if err != nil {
		return "", "", err
	}

	return string(resp.Signature), fmt.Sprintf(APIPath+"/%v", key.Path), nil
}

// Close closes the KMS signer's connections.
func (s *Signer) Close() error {
	return s.client.Close()
}

func sha512Digest(data []byte) []byte {
	output := make([]byte, sha512.Size)
	sha := sha512.Sum512(data)
	copy(output, sha[:])
	return output
}
