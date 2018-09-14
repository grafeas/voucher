package voucher

import (
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp"
)

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

// AddKeyToKeyRingFromReader imports the PGP keys stored in a Reader into the
// passed KeyRing.
func AddKeyToKeyRingFromReader(keyring *KeyRing, name string, reader io.Reader) error {
	var err error
	var newKeyring openpgp.EntityList

	newKeyring, err = openpgp.ReadArmoredKeyRing(reader)
	if nil != err {
		return err
	}

	keyring.AddEntities(name, newKeyring)

	return nil
}
