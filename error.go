package delta

import "fmt"

var (
	// ErrTypeComparison is returned when two incompatible types are compared.
	ErrTypeComparison = fmt.Errorf("incompatible type comparison")
)

// CannotCompareTypes returns ErrTypeComparison specifying the types of the
// two supplied subjects.
func CannotCompareTypes(a, b any) error {
	return fmt.Errorf("%w: %T not comparable to %T", ErrTypeComparison, a, b)
}
