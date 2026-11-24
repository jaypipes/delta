package genutil

// UniquePairs returns pairs of unique values for the supplied slice of
// things.
//
// Important: The supplied slice should have a stable sort order for
// UniquePairs to have a deterministic result.
func UniquePairs[T any](s []T) [][2]T {
	res := make([][2]T, 0, len(s)*(len(s)-1))
	for i, a := range s {
		for j, b := range s {
			if i != j {
				res = append(res, [2]T{a, b})
			}
		}
	}
	return res
}
