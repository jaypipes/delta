package types_test

import (
	"fmt"

	testtypes "github.com/jaypipes/delta/internal/testing/types"
	"github.com/jaypipes/delta/types"
)

func ExampleCategoryOf() {
	subjects := []any{
		nil,
		int(0),
		uint(0),
		float32(0.0),
		int8(0),
		testtypes.MyIntType(0),
		[]int64{0, 1, 2},
		testtypes.MyStructType{
			A: 0,
		},
		map[string]string{
			"A": "a",
		},
	}

	for _, subject := range subjects {
		cat := types.CategoryOf(subject)
		fmt.Printf(
			"%v (%T) is types.%s\n",
			subject, subject, cat,
		)
	}

	// Output:
	// <nil> (<nil>) is types.Nil
	// 0 (int) is types.Scalar
	// 0 (uint) is types.Scalar
	// 0 (float32) is types.Scalar
	// 0 (int8) is types.Scalar
	// 0 (types.MyIntType) is types.Scalar
	// [0 1 2] ([]int64) is types.Slice
	// {0} (types.MyStructType) is types.Record
	// map[A:a] (map[string]string) is types.Record
}

func ExampleCategoryOfElem() {
	subjects := []any{
		nil,
		int(0),
		[]uint{0},
		[]testtypes.MyIntType{0},
		testtypes.MyStructType{
			A: 0,
		},
		[]testtypes.MyStructType{
			{
				A: 0,
			},
		},
		map[string]string{
			"A": "a",
		},
		[]map[string]string{
			{
				"A": "a",
			},
		},
	}

	for _, subject := range subjects {
		elCat, isSlice := types.CategoryOfElem(subject)
		if !isSlice {
			fmt.Printf(
				"%v (%T) is not a types.Slice\n",
				subject, subject,
			)
		} else {
			fmt.Printf(
				"%v (%T) element type is types.%s\n",
				subject, subject, elCat,
			)
		}
	}

	// Output:
	// <nil> (<nil>) is not a types.Slice
	// 0 (int) is not a types.Slice
	// [0] ([]uint) element type is types.Scalar
	// [0] ([]types.MyIntType) element type is types.Scalar
	// {0} (types.MyStructType) is not a types.Slice
	// [{0}] ([]types.MyStructType) element type is types.Record
	// map[A:a] (map[string]string) is not a types.Slice
	// [map[A:a]] ([]map[string]string) element type is types.Record
}
