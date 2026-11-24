package types

// Category is the top-level class of data and variables.
type Category uint8

const (
	// Nil is specifically for the `nil` value, nothing else.
	Nil Category = iota
	// Scalar is for things that represent a singular value
	Scalar
	// Slice is for things contain zero or more values addressable via an
	// integer index
	Slice
	// Record is for things that contain zero or more values addressable via a
	// named key or field
	Record
)

// String returns the string representation of the Category.
func (c Category) String() string {
	switch c {
	case Nil:
		return "Nil"
	case Scalar:
		return "Scalar"
	case Slice:
		return "Slice"
	default:
		return "Record"
	}
}

// CategoryOf returns the TypeCategory of the supplied subject.
func CategoryOf(subject any) Category {
	switch {
	case subject == nil:
		return Nil
	case IsScalar(subject):
		return Scalar
	case IsSlice(subject):
		return Slice
	default:
		return Record
	}
}

// CategoryOfElem returns the TypeCategory of the elements of the supplied
// Slice subject. Returns Nil, false when subject is not a Slice.
func CategoryOfElem(subject any) (Category, bool) {
	switch {
	case IsSlice(subject):
		switch {
		case IsSliceScalar(subject):
			return Scalar, true
		case IsSliceSlice(subject):
			return Slice, true
		default:
			return Record, true
		}
	default:
		return Nil, false
	}
}
