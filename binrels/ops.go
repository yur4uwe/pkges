package binrels

func foreachcell(size int, f func(i, j int) bool) [][]bool {
	result := Zero(size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			result[i][j] = f(i, j)
		}
	}
	return result
}

func Intersection(a, b [][]bool) [][]bool {
	if len(a) != len(b) || len(a) == 0 || len(a[0]) != len(b[0]) {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return a[i][j] && b[i][j]
	})
}

func Union(a, b [][]bool) [][]bool {
	if len(a) != len(b) || len(a) == 0 || len(a[0]) != len(b[0]) {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return a[i][j] || b[i][j]
	})
}

func Diff(a, b [][]bool) [][]bool {
	if len(a) != len(b) || len(a) == 0 || len(a[0]) != len(b[0]) {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return a[i][j] && !b[i][j]
	})
}

func SymmDiff(a, b [][]bool) [][]bool {
	if len(a) != len(b) || len(a) == 0 || len(a[0]) != len(b[0]) {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return a[i][j] && !b[i][j] || !a[i][j] && b[i][j]
	})
}

func Composition(a, b [][]bool) [][]bool {
	if len(a) != len(b) || len(a) == 0 || len(a[0]) != len(b[0]) {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		for k := 0; k < len(a); k++ {
			if a[i][k] && b[k][j] {
				return true
			}
		}
		return false
	})
}
