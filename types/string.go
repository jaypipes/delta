package types

// IsString returns true if the supplied subject is the built-in Go string type
// or a type definition over the built-in Go string type.
func IsString(subject any) bool {
	if IsPrimitiveString(subject) {
		return true
	}
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	up := UnderlyingPrimitive(subject)
	return IsString(up)
}
