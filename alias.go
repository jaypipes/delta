package delta

import (
	"github.com/jaypipes/delta/cmp"
	"github.com/jaypipes/delta/diff"
	"github.com/jaypipes/delta/types"
)

var (
	IsScalar           = types.IsScalar
	IsSlice            = types.IsSlice
	IsPrimitive        = types.IsPrimitive
	IsPrimitiveTypeDef = types.IsPrimitiveTypeDef

	IsInteger              = types.IsInteger
	IsIntegerTypeDef       = types.IsIntegerTypeDef
	IsSignedIntger         = types.IsSignedInteger
	IsSignedIntegerTypeDef = types.IsSignedIntegerTypeDef

	IsFloat        = types.IsFloat
	IsFloatTypeDef = types.IsFloatTypeDef

	IsString           = types.IsString
	IsByteArray        = types.IsByteArray
	IsByteArrayTypeDef = types.IsByteArrayTypeDef
)

var (
	Comparable      = cmp.Comparable
	Compare         = cmp.Compare
	Equal           = cmp.Equal
	EqualComparable = cmp.EqualComparable
)

var (
	Diff = diff.Diff
)
