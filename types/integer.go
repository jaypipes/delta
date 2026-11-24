package types

// IsInteger returns true if the supplied subject is one of the built-in Go
// integer types or a type definition over a built-in Go integer type.
func IsInteger(subject any) bool {
	if IsPrimitiveInteger(subject) {
		return true
	}
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	up := UnderlyingPrimitive(subject)
	return IsInteger(up)
}

// IsSigned returns true if the supplied subject is one of the built-in signed
// Go integer types or a type definition over a built-in signed Go integer
// type.
func IsSigned(subject any) bool {
	if IsPrimitiveSigned(subject) {
		return true
	}
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	up := UnderlyingPrimitive(subject)
	return IsSigned(up)
}
