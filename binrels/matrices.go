package binrels

import "fmt"

func Copy(a [][]bool) [][]bool {
	if len(a) == 0 {
		return nil
	}

	return foreachcell(len(a), func(i int, j int) bool {
		return a[i][j]
	})
}

func Equal(a, b [][]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}

		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func Zero(n int) [][]bool {
	matrix := make([][]bool, n)
	for i := range matrix {
		matrix[i] = make([]bool, n)
	}
	return matrix
}

func Identity(n int) [][]bool {
	matrix := Zero(n)
	for i := range matrix {
		matrix[i][i] = true
	}
	return matrix
}

func PrintWithSource(source []string, relationship [][]bool) {
	minSourceNameLen := 0
	for _, s := range source {
		if len(s) > minSourceNameLen {
			minSourceNameLen = len(s)
		}
	}

	fmt.Printf("%-*s | ", minSourceNameLen, "")
	for i := range source {
		fmt.Printf("%-*s", minSourceNameLen, source[i])
		if i < len(source)-1 {
			fmt.Printf("| ")
		}
	}
	for i := range source {
		fmt.Println()
		fmt.Println("--------------------")
		fmt.Printf(" %-*s| ", minSourceNameLen, source[i])
		for j := range source {
			if relationship[i][j] {
				fmt.Printf("%-*s", minSourceNameLen, "1")
			} else {
				fmt.Printf("%-*s", minSourceNameLen, "0")
			}
			if j < len(source)-1 {
				fmt.Printf("| ")
			}
		}
	}
	fmt.Println()
}

func Print(relationship [][]bool) {
	for i := range relationship {
		fmt.Println()
		fmt.Println("--------------------")
		for j := range relationship[i] {
			if relationship[i][j] {
				fmt.Printf(" 1 ")
			} else {
				fmt.Printf(" 0 ")
			}
			if j < len(relationship[i])-1 {
				fmt.Printf("| ")
			}
		}
	}
	fmt.Println()
}
