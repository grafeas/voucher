package voucher

import "fmt"

// Severity is a integer that represents how severe a vulnerability
// is.
type Severity int

// Severity constants, which represent the severities that we track. Other systems'
// severities should be converted to one of the following.
const (
	NegligibleSeverity Severity = iota
	LowSeverity        Severity = iota
	MediumSeverity     Severity = iota
	UnknownSeverity    Severity = iota
	HighSeverity       Severity = iota
	CriticalSeverity   Severity = iota
)

const (
	negligibleSeverityString = "negligible"
	lowSeverityString        = "low"
	mediumSeverityString     = "medium"
	highSeverityString       = "high"
	criticalSeverityString   = "critical"
	unknownSeverityString    = "unknown"
)

// String returns a string representation of a Severity.
func (s Severity) String() string {
	switch s {
	case NegligibleSeverity:
		return negligibleSeverityString
	case LowSeverity:
		return lowSeverityString
	case MediumSeverity:
		return mediumSeverityString
	case HighSeverity:
		return highSeverityString
	case CriticalSeverity:
		return criticalSeverityString
	}
	return unknownSeverityString
}

// StringToSeverity returns the matching Severity to the passed string.
// Returns an error if there isn't a matching Severity.
func StringToSeverity(s string) (Severity, error) {
	switch s {
	case negligibleSeverityString:
		return NegligibleSeverity, nil
	case lowSeverityString:
		return LowSeverity, nil
	case mediumSeverityString:
		return MediumSeverity, nil
	case highSeverityString:
		return HighSeverity, nil
	case criticalSeverityString:
		return CriticalSeverity, nil
	case unknownSeverityString:
		return UnknownSeverity, nil
	}
	return UnknownSeverity, fmt.Errorf("severity %s doesn't exist", s)
}
