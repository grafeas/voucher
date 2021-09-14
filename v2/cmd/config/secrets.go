package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Shopify/ejson"
	"github.com/spf13/viper"
	"go.mozilla.org/sops/v3/decrypt"

	"github.com/grafeas/voucher/v2/clair"
	"github.com/grafeas/voucher/v2/repository"
	"github.com/grafeas/voucher/v2/signer/pgp"
)

// Secrets represents the format that the ejson configuration is structured
// in.
type Secrets struct {
	Keys                     map[string]string  `json:"openpgpkeys"`
	ClairConfig              clair.Config       `json:"clair"`
	RepositoryAuthentication repository.KeyRing `json:"repositories"`
	Datadog                  DatadogSecrets     `json:"datadog"`
}

type DatadogSecrets struct {
	APIKey string `json:"api_key"`
	AppKey string `json:"app_key"`
}

// ReadSecrets reads from the ejson file and populates the passed interface.
func ReadSecrets() (*Secrets, error) {
	decrypted, err := decryptSecrets()
	if err != nil {
		return nil, err
	}

	var data Secrets
	if err := json.Unmarshal(decrypted, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func decryptSecrets() ([]byte, error) {
	ejDir := viper.GetString("ejson.dir")
	ejSecrets := viper.GetString("ejson.secrets")
	if ejDir != "" && ejSecrets != "" {
		return ejson.DecryptFile(ejSecrets, ejDir, "")
	}

	sops := viper.GetString("sops.file")
	if sops != "" {
		return decrypt.File(sops, "json")
	}

	return nil, fmt.Errorf("secrets not provided via ejson or sops")
}

// getPGPKeyRing uses the Command's configured ejson file to populate a
// voucher.KeyRing.
func (s *Secrets) getPGPKeyRing() (*pgp.KeyRing, error) {
	newKeyRing := pgp.NewKeyRing()

	for name, key := range s.Keys {
		err := pgp.AddKeyToKeyRingFromReader(newKeyRing, name, bytes.NewReader([]byte(key)))
		if nil != err {
			return nil, err
		}
	}

	return newKeyRing, nil
}
