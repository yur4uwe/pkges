package binrels

import "slices"

func Power(a [][]bool, n int) [][]bool {
	if n == 0 {
		return Identity(len(a))
	}

	if n == 1 {
		return a
	}

	if n&1 == 0 {
		half := Power(a, n/2)
		return Composition(half, half)
	}

	return Composition(a, Power(a, n-1))
}

func Transpose(a [][]bool) [][]bool {
	if len(a) == 0 {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return a[j][i]
	})
}

func Complement(a [][]bool) [][]bool {
	if len(a) == 0 {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return !a[i][j]
	})
}

func DefinitionDomain(a [][]bool) []int {
	if len(a) == 0 {
		return nil
	}

	res := make([]int, 0, len(a))
	for i := range a {
		for j := range a[0] {
			if a[i][j] && !slices.Contains(res, i) {
				res = append(res, i)
				break
			}
		}
	}
	return res
}

func MeaningDomain(a [][]bool) []int {
	if len(a) == 0 {
		return nil
	}

	res := make([]int, 0, len(a))
	for i := range a {
		for j := range a[0] {
			if a[i][j] && !slices.Contains(res, j) {
				res = append(res, j)
				break
			}
		}
	}
	return res
}

func BottomIntersection(a [][]bool, x int) []int {
	if len(a) == 0 || x < 0 || x >= len(a) {
		return nil
	}
	res := make([]int, 0, len(a[x]))
	for y := range a[x] {
		if a[x][y] {
			res = append(res, y)
		}
	}
	return res
}

func TopIntersection(a [][]bool, x int) []int {
	if len(a) == 0 || x < 0 || x >= len(a) {
		return nil
	}
	res := make([]int, 0, len(a))
	for y := range a {
		if a[y][x] {
			res = append(res, y)
		}
	}
	return res
}

func TransitiveClosure(a [][]bool) [][]bool {
	if len(a) == 0 {
		return nil
	}

	accumulator := Copy(a)
	for i := 2; ; i++ {
		next_a := Power(a, i)

		before := Copy(accumulator)
		accumulator = Union(accumulator, next_a)

		if Equal(before, accumulator) {
			break
		}
	}

	return accumulator
}

func Reachability(a [][]bool) [][]bool {
	return Union(Identity(len(a)), TransitiveClosure(a))
}

func MutualReachability(a [][]bool) [][]bool {
	reach := Reachability(a)

	return Intersection(reach, Transpose(reach))
}
