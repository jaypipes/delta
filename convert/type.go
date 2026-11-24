package convert

import "github.com/jaypipes/delta/types"

// Type describes the conversion type of a thing.
//
// Type is different from [reflect.Type] in some important ways.
//
// Type combines signed and unsigned integer types of all sizes into
// a single Integer Type.
//
// This allows [Compare] and [Equal] to compare differently sized and signed
// integer types and differently size float types without excessive conversions
// by simply type-casting the underlying values to the largest sized integer or
// float type.
//
// Note that the zero value of Type is Invalid, which mirrors the behaviour in
// the [reflect] package.
type Type uint

const (
	TypeInvalid Type = iota
	TypeInteger
	TypeFloat
	TypeString
	TypeBool
)

// TypeOf returns the conversion type of the supplied subject.
func TypeOf(subject any) Type {
	switch subject := subject.(type) {
	case Type:
		return subject
	default:
		switch {
		case types.IsInteger(subject):
			return TypeInteger
		case types.IsString(subject):
			return TypeString
		case types.IsFloat(subject):
			return TypeFloat
		case types.IsBool(subject):
			return TypeBool
		default:
			return TypeInvalid
		}
	}
}

// Convertible returns true if the supplied subjects are convertible to each
// other.
func Convertible(a any, b any) bool {
	if types.IsScalar(a) {
		return types.IsScalar(b)
	}
	if types.IsSliceScalar(a) {
		return types.IsSliceScalar(b)
	}
	return false
}
