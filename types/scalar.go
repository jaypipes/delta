package types

// IsScalar returns true if the supplied subject belongs to the Scalar type
// category.
//
// The Scalar type category encompasses all built-in Go types that are
// primitives (see [IsPrimitive]) and anything that is a Go type definition
// over a primitive (see [IsPrimitiveTypeDef])
//
// Note `nil` is not a Scalar.
func IsScalar(subject any) bool {
	return IsPrimitive(subject) || IsPrimitiveTypeDef(subject)
}
