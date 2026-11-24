package cmp

import "fmt"

var (
	// ErrInvalidTypeComparison is returned when two incompatible types are compared.
	ErrInvalidTypeComparison = fmt.Errorf("invalid type comparison")
)

// InvalidTypeComparison returns ErrInvalidTypeComparison specifying the types
// of the two supplied subjects.
func InvalidTypeComparison(a, b any) error {
	return fmt.Errorf(
		"%w: %T not comparable to %T",
		ErrInvalidTypeComparison, a, b,
	)
}
