package voucher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Shopify/ejson"
	"golang.org/x/crypto/openpgp"
)

var errNoKeys = errors.New("no keys in ejson file")
var errKeysNotMap = errors.New("keys in ejson is not a map")

// KeyRing wraps an OpenPGP EntityList (which implements openpgp.KeyRing),
// adding support for determining which key is associated with which check.
// KeyRing implements openpgp.KeyRing, thus can be used in place of it where
// appropriate.
type KeyRing struct {
	keyIds   map[string]uint64
	entities openpgp.EntityList
}

// KeysById returns the set of keys that have the given key id.
func (keyring *KeyRing) KeysById(id uint64) []openpgp.Key {
	if nil == keyring {
		return nil
	}
	return keyring.entities.KeysById(id)
}

// KeysByIdUsage returns the set of keys with the given id
// that also meet the key usage given by requiredUsage.
// The requiredUsage is expressed as the bitwise-OR of
// packet.KeyFlag* values.
func (keyring *KeyRing) KeysByIdUsage(id uint64, requiredUsage byte) []openpgp.Key {
	if nil == keyring {
		return nil
	}
	return keyring.entities.KeysByIdUsage(id, requiredUsage)
}

// DecryptionKeys returns all private keys that are valid for
// decryption.
func (keyring *KeyRing) DecryptionKeys() []openpgp.Key {
	if nil == keyring {
		return nil
	}
	return keyring.entities.DecryptionKeys()
}

// getSigner gets the first available signing Entity with the passed id.
func (keyring *KeyRing) getSigner(id uint64) (*openpgp.Entity, error) {
	if nil == keyring {
		return nil, fmt.Errorf("keyring is empty")
	}

	for _, entity := range keyring.entities {
		if nil != entity.PrimaryKey {
			if id == entity.PrimaryKey.KeyId {
				return entity, nil
			}
		}
	}

	return nil, fmt.Errorf("no signing entity exists for id \"%#x\"", id)
}

// GetSignerByName gets the first available signing key associated with the passed name.
func (keyring *KeyRing) GetSignerByName(name string) (*openpgp.Entity, error) {
	if nil != keyring {
		keyID := keyring.keyIds[name]
		if 0 != keyID {
			return keyring.getSigner(keyID)
		}
	}
	return nil, fmt.Errorf("no signing entity exists for check name \"%s\"", name)
}

// AddEntities adds new keys from the passed EntityList to the keyring for
// use.
func (keyring *KeyRing) AddEntities(name string, input openpgp.EntityList) {
	if nil == keyring {
		keyring = NewKeyRing()
	}

	for _, entity := range input {
		if nil != entity.PrimaryKey {
			keyring.entities = append(keyring.entities, entity)
			keyring.keyIds[name] = entity.PrimaryKey.KeyId
		}
	}
}

// NewKeyRing creates a new keyring from the passed EntityList. The keys in
// the input EntityList are then associated with the
func NewKeyRing() *KeyRing {
	keyring := new(KeyRing)

	keyring.keyIds = make(map[string]uint64)
	keyring.entities = make(openpgp.EntityList, 0)

	return keyring
}

// AddEntitiesToKeyRingFromString imports the PGP keys stored in a string into the
// passed KeyRing.
func addEntitiesToKeyRingFromReader(keyring *KeyRing, name string, reader io.Reader) error {
	var err error
	var newKeyring openpgp.EntityList

	newKeyring, err = openpgp.ReadArmoredKeyRing(reader)
	if nil != err {
		return err
	}

	keyring.AddEntities(name, newKeyring)

	return nil
}

// EjsonToKeyRing takes an ejson directory path, and an ejson filename, and returns
// a new KeyRing with the keys located in that file.
func EjsonToKeyRing(dir, filename string) (*KeyRing, error) {
	newKeyRing := NewKeyRing()
	secrets, err := readEjson(dir, filename)
	if err != nil {
		return nil, err
	}

	keys, err := extractKeys(secrets)
	if err != nil {
		return nil, err
	}

	for name, key := range keys {
		err = addEntitiesToKeyRingFromReader(newKeyRing, name, bytes.NewReader([]byte(key)))
		if nil != err {
			return nil, err
		}
	}

	return newKeyRing, nil
}

// readEjson reads from the ejson file and returns a map[string]interface{}.
func readEjson(dir, filename string) (map[string]interface{}, error) {
	secrets := make(map[string]interface{})

	decrypted, err := ejson.DecryptFile(filename, dir, "")
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(bytes.NewReader(decrypted)).Decode(&secrets)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

// extractEjsonKeys extracts the OpenPGP keys from the map[string]interface{}
// containing all secrets, and returns a map[string]string containing the
// key value pairs. If there's an issue (the environment key doesn't exist, for
// example), returns an error.
func extractKeys(secrets map[string]interface{}) (map[string]string, error) {
	rawKeys, ok := secrets["openpgpkeys"]
	if !ok {
		return nil, errNoKeys
	}

	keysMap, ok := rawKeys.(map[string]interface{})
	if !ok {
		return nil, errKeysNotMap
	}

	keysSecrets := make(map[string]string, len(keysMap))

	for key, rawValue := range keysMap {

		// Only export values that convert to strings properly.
		if value, ok := rawValue.(string); ok {
			keysSecrets[key] = value
		}
	}

	return keysSecrets, nil
}
