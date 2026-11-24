package types

// IsFloat returns true if the supplied subject is one of the built-in Go float
// types or a type definition over a built-in Go float type.
func IsFloat(subject any) bool {
	if IsPrimitiveFloat(subject) {
		return true
	}
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	up := UnderlyingPrimitive(subject)
	return IsFloat(up)
}
