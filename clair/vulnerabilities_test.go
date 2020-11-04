package clair

import (
	v1 "github.com/coreos/clair/api/v1"

	"github.com/grafeas/voucher"
)

// ClairVulnerabilitiesV1 return a list of clair vulnerabilities
func ClairVulnerabilities() map[string][]v1.Vulnerability {
	vulns := map[string][]v1.Vulnerability{
		"sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b": {
			{
				Name:        "bad vul 1",
				Description: "some bad vul that i cant comprepend",
				FixedBy:     "Hackerman",
				Severity:    "Medium",
			},
			{
				Name:        "bad vul 2",
				Description: "Dont know what this is",
				FixedBy:     "Hackerwoman",
				Severity:    "Medium",
			},
		},
		"sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736": {
			{
				Name:        "Super bad vul",
				Description: "bark bark",
				FixedBy:     "Hackerdog",
				Severity:    "Low",
			},
			{
				Name:        "Super bad vul 2 meow",
				Description: "meow meow",
				FixedBy:     "Hackercat",
				Severity:    "High",
			},
		},
	}
	return vulns
}

// VoucherVulnerabilities return a list of voucher vulnerabilies
// corresponding to ClairVulnerabilites
func VoucherVulnerabilities(keys ...string) []voucher.Vulnerability {
	vulns := map[string][]voucher.Vulnerability{
		"low": {
			{
				Name:        "Super bad vul",
				Description: "bark bark",
				FixedBy:     "Hackerdog",
				Severity:    voucher.LowSeverity,
			},
		},
		"medium": {
			{
				Name:        "bad vul 1",
				Description: "some bad vul that i cant comprepend",
				FixedBy:     "Hackerman",
				Severity:    voucher.MediumSeverity,
			},
			{
				Name:        "bad vul 2",
				Description: "Dont know what this is",
				FixedBy:     "Hackerwoman",
				Severity:    voucher.MediumSeverity,
			},
		},
		"high": {
			{
				Name:        "Super bad vul 2 meow",
				Description: "meow meow",
				FixedBy:     "Hackercat",
				Severity:    voucher.HighSeverity,
			},
		},
	}
	vulnlist := []voucher.Vulnerability{}
	for _, key := range keys {
		vulnlist = append(vulnlist, vulns[key]...)
	}
	return vulnlist
}
