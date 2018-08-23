package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Shopify/ejson"
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/clair"
	"github.com/spf13/viper"
)

// ejsonFormat represents the format that the ejson configuration is structured
// in.
type ejsonFormat struct {
	Keys        map[string]string `json:"openpgpkeys"`
	ClairConfig clair.Config      `json:"clair"`
}

// readEjson reads from the ejson file and populates the passed interface.
func readEjson(data interface{}) error {
	if !viper.IsSet("ejson.dir") {
		return fmt.Errorf("ESON dir not set in the config file")
	}

	dir := viper.GetString("ejson.dir")
	if !viper.IsSet("ejson.secrets") {
		return fmt.Errorf("ESON secrets not set in the config file")
	}

	secrets := viper.GetString("ejson.secrets")

	decrypted, err := ejson.DecryptFile(secrets, dir, "")
	if err != nil {
		return err
	}

	err = json.NewDecoder(bytes.NewReader(decrypted)).Decode(data)
	return err
}

// getClairConfig uses the Command's configured ejson file to populate a
// clair.Config.
func getClairConfig() (clair.Config, error) {
	ejsonData := new(ejsonFormat)

	err := readEjson(ejsonData)
	if err != nil {
		return clair.Config{}, err
	}

	return ejsonData.ClairConfig, nil
}

// getKeyRing uses the Command's configured ejson file to populate a
// voucher.KeyRing.
func getKeyRing() (*voucher.KeyRing, error) {
	newKeyRing := voucher.NewKeyRing()

	ejsonData := new(ejsonFormat)

	err := readEjson(ejsonData)
	if err != nil {
		return nil, err
	}

	for name, key := range ejsonData.Keys {
		err = voucher.AddKeyToKeyRingFromReader(newKeyRing, name, bytes.NewReader([]byte(key)))
		if nil != err {
			return nil, err
		}
	}

	return newKeyRing, nil
}
