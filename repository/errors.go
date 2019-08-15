package repository

import "fmt"

// TypeMismatchError represents a type mismatch between objects
type typeMismatchError struct {
	expectedType string
	actualType   string
}

func (t *typeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch found. Expected: %s, Actual: %s", t.expectedType, t.actualType)
}

// NewTypeMismatchError creates a new TypeMismatchError
func NewTypeMismatchError(expected string, actual string) error {
	return &typeMismatchError{
		expectedType: expected,
		actualType:   actual,
	}
}
