package types

// IsByteArray returns true if the supplied subject is a `[]byte` or a type
// definition over `[]byte`
func IsByteArray(subject any) bool {
	switch subject.(type) {
	case []byte:
		return true
	default:
		return IsByteArrayTypeDef(subject)
	}
}

// IsByteArrayTypeDef returns true if the supplied subject is a type
// definition over `[]byte`.
//
// Note that Both `string` and `MyByteArrayType` are scalars.  However, `string`
// is not a slice, but `MyByteArrayType` *is* a slice.
func IsByteArrayTypeDef(subject any) bool {
	return IsPrimitiveTypeDef(subject) && IsSlice(subject)
}
