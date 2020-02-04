package signer

import (
	"fmt"
)

var NoKeyForCheckError = fmt.Errorf("no signing entity exists for check")
