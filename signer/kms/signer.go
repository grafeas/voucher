package kms

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	apiv1 "cloud.google.com/go/kms/apiv1"
	"github.com/grafeas/voucher/signer"
	kms_pb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	AlgoSHA256 = "SHA256"
	AlgoSHA384 = "SHA384"
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

	for checkName, key := range keys {
		if key.Algo != AlgoSHA256 && key.Algo != AlgoSHA384 && key.Algo != AlgoSHA512 {
			return nil, fmt.Errorf("Unsupported digest algorithm %v for check %v", key.Algo, checkName)
		}
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

	var digest hash.Hash
	if _, err := digest.Write([]byte(body)); err != nil {
		return "", "", err
	}

	var d kms_pb.Digest
	switch key.Algo {
	case AlgoSHA256:
		digest = sha256.New()
		d.Digest = &kms_pb.Digest_Sha256{
			Sha256: digest.Sum(nil),
		}
	case AlgoSHA384:
		digest = sha512.New384()
		d.Digest = &kms_pb.Digest_Sha384{
			Sha384: digest.Sum(nil),
		}
	case AlgoSHA512:
		digest = sha512.New()
		d.Digest = &kms_pb.Digest_Sha512{
			Sha512: digest.Sum(nil),
		}
	default:
		return "", "", fmt.Errorf("Unsupported digest algorithm %v", key.Algo)
	}

	resp, err := s.client.AsymmetricSign(context.Background(), &kms_pb.AsymmetricSignRequest{
		Name: key.Path,
		Digest: &d,
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
