package kms

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	apiv1 "cloud.google.com/go/kms/apiv1"
	"github.com/googleapis/gax-go/v2"
	"github.com/grafeas/voucher/v2/signer"
	kms_pb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	AlgoSHA256 = "SHA256"
	AlgoSHA384 = "SHA384"
	AlgoSHA512 = "SHA512"
	APIPath    = "//cloudkms.googleapis.com/v1"
)

// kmsClient is a subset of cloud.google.com/go/kms/apiv1.KeyManagementClient
type kmsClient interface {
	AsymmetricSign(ctx context.Context, req *kms_pb.AsymmetricSignRequest, opts ...gax.CallOption) (*kms_pb.AsymmetricSignResponse, error)
	Close() error
}

var _ kmsClient = (*apiv1.KeyManagementClient)(nil)

type Key struct {
	Path string
	Algo string
}

// Signer is an AttestationSigner that uses Google's Cloud KMS to sign attestations
// Only supports SHA512 digests.
type Signer struct {
	keys   map[string]Key
	client kmsClient
}

func NewSigner(keys map[string]Key, opts ...SignerOpt) (*Signer, error) {
	for checkName, key := range keys {
		switch key.Algo {
		case AlgoSHA256, AlgoSHA384, AlgoSHA512:
			// supported
		default:
			return nil, fmt.Errorf("unsupported digest algorithm %v for check %v", key.Algo, checkName)
		}
	}

	s := &Signer{keys: keys}
	for _, o := range opts {
		o(s)
	}

	if s.client == nil {
		client, err := apiv1.NewKeyManagementClient(context.Background())
		if err != nil {
			return nil, err
		}
		s.client = client
	}
	return s, nil
}

type SignerOpt func(*Signer)

func WithKMSClient(client kmsClient) SignerOpt {
	return func(s *Signer) {
		s.client = client
	}
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
		return "", "", fmt.Errorf("unsupported digest algorithm %v", key.Algo)
	}

	resp, err := s.client.AsymmetricSign(context.Background(), &kms_pb.AsymmetricSignRequest{
		Name:   key.Path,
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
