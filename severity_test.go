package voucher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSeverities = map[string]Severity{
	negligibleSeverityString: NegligibleSeverity,
	lowSeverityString:        LowSeverity,
	mediumSeverityString:     MediumSeverity,
	unknownSeverityString:    UnknownSeverity,
	highSeverityString:       HighSeverity,
	criticalSeverityString:   CriticalSeverity,
	"whatever":               UnknownSeverity,
}

func TestSeverityToString(t *testing.T) {
	assert := assert.New(t)

	for expected, severity := range testSeverities {
		if "whatever" == expected {
			continue
		}
		value := severity.String()
		assert.Equalf(value, expected, "Severity.String() returned the wrong output, should be: %v, was %v", expected, value)
	}
}

func TestStringToSeverity(t *testing.T) {
	assert := assert.New(t)

	for name, expected := range testSeverities {
		value, err := StringToSeverity(name)

		if nil != err {
			assert.Equal(string(name), "whatever", "got error converting severities: ", err)
			continue
		}

		assert.Equalf(value, expected, "StringToSeverity returned the wrong Severity, should be: %v, was %v", expected, value)

		if "whatever" == name {
			assert.Error(err)
		}
	}
}
