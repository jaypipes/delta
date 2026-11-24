package types

// IsBool returns true if the supplied subject is the built-in Go bool type or
// a type definition over the built-in Go bool type.
func IsBool(subject any) bool {
	if IsPrimitiveBool(subject) {
		return true
	}
	if !IsPrimitiveTypeDef(subject) {
		return false
	}
	up := UnderlyingPrimitive(subject)
	return IsBool(up)
}
