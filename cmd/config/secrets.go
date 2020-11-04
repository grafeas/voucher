package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Shopify/ejson"
	"github.com/spf13/viper"

	"github.com/grafeas/voucher/clair"
	"github.com/grafeas/voucher/repository"
	"github.com/grafeas/voucher/signer/pgp"
)

// Secrets represents the format that the ejson configuration is structured
// in.
type Secrets struct {
	Keys                     map[string]string  `json:"openpgpkeys"`
	ClairConfig              clair.Config       `json:"clair"`
	RepositoryAuthentication repository.KeyRing `json:"repositories"`
}

// ReadSecrets reads from the ejson file and populates the passed interface.
func ReadSecrets() (*Secrets, error) {
	if !viper.IsSet("ejson.dir") {
		return nil, fmt.Errorf("EJSON dir not set in the config file")
	}

	dir := viper.GetString("ejson.dir")
	if !viper.IsSet("ejson.secrets") {
		return nil, fmt.Errorf("EJSON secrets not set in the config file")
	}

	secrets := viper.GetString("ejson.secrets")

	decrypted, err := ejson.DecryptFile(secrets, dir, "")
	if err != nil {
		return nil, err
	}

	var data Secrets
	err = json.NewDecoder(bytes.NewReader(decrypted)).Decode(&data)
	return &data, err
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
