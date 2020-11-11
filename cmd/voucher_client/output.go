package main

import (
	"fmt"
	"os"

	"github.com/grafeas/voucher"
)

// errorf prints a formatted string to standard error.
func errorf(format string, v interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", v)
}

// formatResponse returns the response as a string.
func formatResponse(resp *voucher.Response) string {
	output := ""
	if resp.Success {
		fmt.Println("image is approved")
	} else {
		fmt.Println("image was rejected")
	}
	for _, result := range resp.Results {
		if result.Success {
			output += fmt.Sprintf("   ✓ passed %s", result.Name)
			if !result.Attested {
				output += ", but wasn't attested"
			}
		} else {
			output += fmt.Sprintf("   ✗ failed %s", result.Name)
		}

		if "" != result.Err {
			output += fmt.Sprintf(", err: %s", result.Err)
		}
		output += "\n"
	}

	return output
}
