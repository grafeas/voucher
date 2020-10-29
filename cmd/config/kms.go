package config

import (
	"github.com/grafeas/voucher/signer/kms"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getKMSKeyRing() (*kms.Signer, error) {
	rows, ok := viper.Get("kms_keys").([]interface{})
	if !ok {
		log.Warning("KMS keys not configured")
		return nil, nil
	}

	keys := make(map[string]kms.Key)
	for _, row := range rows {
		if m, ok := row.(map[string]interface{}); ok {
			check := m["check"].(string)
			path := m["path"].(string)
			algo := m["algo"].(string)
			keys[check] = kms.Key{Path: path, Algo: algo}
		} else {
			continue
		}
	}

	return kms.NewSigner(keys)
}
