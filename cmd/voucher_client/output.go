package main

import (
	"fmt"

	"github.com/Shopify/voucher"
)

// formatResponse returns the response as a string.
func formatResponse(resp *voucher.Response) string {
	output := ""
	for _, result := range resp.Results {
		if result.Success && "" == result.Err {
			output += fmt.Sprintf("   ✓ %s succeeded\n", result.Name)
			continue
		}
		if !result.Success {
			output += fmt.Sprintf("   ✗ %s failed", result.Name)
		} else if !result.Attested {
			output += fmt.Sprintf("   ✓ %s succeeded, but wasn't attested", result.Name)
		}

		if "" != result.Err {
			output += fmt.Sprintf(", err: %s", result.Err)
		}
		output += "\n"
	}

	return output
}
