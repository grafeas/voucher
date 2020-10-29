package server

import (
	"fmt"
	"strings"

	"github.com/grafeas/voucher"
)

func verifiedRequiredChecksAreRegistered(checks ...string) error {
	disabledChecks := make([]string, 0, len(checks))
	for _, check := range checks {
		if !voucher.IsCheckFactoryRegistered(check) {
			disabledChecks = append(disabledChecks, check)
		}
	}

	if len(disabledChecks) != 0 {
		return fmt.Errorf("required check(s) are not registered: %s", strings.Join(disabledChecks, ", "))
	}

	return nil
}
