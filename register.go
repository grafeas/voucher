package voucher

import (
	"fmt"
)

// CheckFactory is a type of function that creates a new Check.
type CheckFactory func() Check

// CheckFactories is a map of registered CheckFactories.
type CheckFactories map[string]CheckFactory

// DefaultCheckFactories is the default CheckFactory collection.
var DefaultCheckFactories = make(CheckFactories)

// Register adds a new CheckFactory to this CheckFactories.
func (cf CheckFactories) Register(name string, creator CheckFactory) {
	if nil == cf[name] {
		cf[name] = creator
	}
}

// Get returns the CheckFactory with the passed name.
func (cf CheckFactories) Get(name string) CheckFactory {
	return cf[name]
}

// RegisterCheckFactory adds a CheckFactory to the DefaultCheckFactories
// that can be run. Once a Check is added, it can be referenced by the name
// that was passed in when this function was called.
func RegisterCheckFactory(name string, creator CheckFactory) {
	DefaultCheckFactories.Register(name, creator)
}

// GetCheckFactories gets new copies of the Checks from their registered
// CheckCheckFactories.
func GetCheckFactories(names ...string) (map[string]Check, error) {
	checks := make(map[string]Check, len(DefaultCheckFactories))
	for _, name := range names {
		creator := DefaultCheckFactories.Get(name)
		if nil == creator {
			return checks, fmt.Errorf("requested check \"%s\" does not exist", name)
		}
		checks[name] = creator()
	}
	return checks, nil
}

// IsCheckFactoryRegistered returns true if the passed CheckFactory was
// registered.
func IsCheckFactoryRegistered(name string) bool {
	return (nil != DefaultCheckFactories.Get(name))
}
