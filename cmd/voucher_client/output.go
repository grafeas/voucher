package main

import (
	"fmt"
	"os"

	"github.com/Shopify/voucher"
)

// errorf prints a formatted string to standard error.
func errorf(format string, v interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", v)
}

// formatResponse returns the response as a string.
func formatResponse(resp *voucher.Response) string {
	output := ""
	for _, result := range resp.Results {
		if result.Success {
			output += fmt.Sprintf("   ✓ %s succeeded", result.Name)
			if !result.Attested {
				output += ", but wasn't attested"
			}
		} else {
			output += fmt.Sprintf("   ✗ %s failed", result.Name)
		}

		if "" != result.Err {
			output += fmt.Sprintf(", err: %s", result.Err)
		}
		output += "\n"
	}

	return output
}
