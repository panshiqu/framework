package fiveinarow

import (
	"fmt"
	"testing"
)

func TestSame(t *testing.T) {
	board := [][]int{
		[]int{1, 0, 1, 0, 0},
		[]int{0, 1, 1, 0, 0},
		[]int{1, 1, 1, 1, 1},
		[]int{0, 0, 1, 1, 0},
		[]int{0, 0, 1, 0, 1},
	}

	fmt.Println("same:", same(board, 2, 2, 1, 1, 1))
}
