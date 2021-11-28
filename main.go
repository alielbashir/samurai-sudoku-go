package main

import (
	"fmt"
	. "samurai-sudoku/internal"
)

func main() {
	samuraiGrid := Grid{
		{0, 0, 5, 7, 0, 0, 0, 2, 0, -1, -1, -1, 0, 0, 9, 6, 0, 0, 0, 2, 0},
		{4, 9, 0, 0, 6, 0, 0, 1, 0, -1, -1, -1, 1, 4, 0, 0, 5, 0, 0, 3, 0},
		{0, 0, 7, 0, 0, 4, 9, 0, 6, -1, -1, -1, 0, 0, 2, 0, 0, 1, 7, 0, 8},
		{0, 0, 6, 0, 0, 0, 0, 0, 8, -1, -1, -1, 0, 0, 3, 0, 0, 0, 0, 0, 2},
		{0, 7, 0, 0, 0, 0, 0, 9, 0, -1, -1, -1, 0, 5, 0, 0, 0, 0, 0, 6, 0},
		{2, 0, 0, 0, 0, 0, 3, 0, 0, -1, -1, -1, 4, 0, 0, 0, 0, 0, 5, 0, 0},
		{5, 0, 8, 9, 0, 0, 7, 0, 0, 0, 0, 0, 6, 0, 5, 8, 0, 0, 2, 0, 0},
		{0, 1, 0, 0, 3, 0, 0, 8, 5, 0, 0, 0, 0, 1, 0, 0, 7, 0, 0, 8, 6},
		{0, 2, 0, 0, 0, 5, 6, 0, 0, 0, 1, 0, 0, 2, 0, 0, 0, 4, 3, 0, 0},
		{-1, -1, -1, -1, -1, -1, 0, 0, 0, 4, 0, 6, 0, 0, 0, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, 0, 0, 6, 0, 5, 0, 2, 0, 0, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, 0, 0, 0, 2, 0, 8, 0, 0, 0, -1, -1, -1, -1, -1, -1},
		{0, 0, 8, 5, 0, 0, 0, 2, 0, 0, 3, 0, 0, 0, 8, 9, 0, 0, 0, 6, 0},
		{6, 2, 0, 0, 4, 0, 0, 5, 0, 0, 0, 0, 9, 6, 0, 0, 2, 0, 0, 5, 0},
		{0, 0, 7, 0, 0, 8, 9, 0, 3, 0, 0, 0, 0, 0, 2, 0, 0, 8, 1, 0, 9},
		{0, 0, 6, 0, 0, 0, 0, 0, 2, -1, -1, -1, 0, 0, 1, 0, 0, 0, 0, 0, 6},
		{0, 5, 0, 0, 0, 0, 0, 4, 0, -1, -1, -1, 0, 8, 0, 0, 0, 0, 0, 2, 0},
		{8, 0, 0, 0, 0, 0, 3, 0, 0, -1, -1, -1, 7, 0, 0, 0, 0, 0, 5, 0, 0},
		{1, 0, 5, 9, 0, 0, 2, 0, 0, -1, -1, -1, 2, 0, 6, 7, 0, 0, 4, 0, 0},
		{0, 3, 0, 0, 6, 0, 0, 7, 1, -1, -1, -1, 0, 3, 0, 0, 9, 0, 0, 7, 8},
		{0, 6, 0, 0, 0, 3, 5, 0, 0, -1, -1, -1, 0, 9, 0, 0, 0, 4, 2, 0, 0},
	}

	var samuraiSudoku SamuraiSudoku

	samuraiSudoku.SetGrid(samuraiGrid)

	fmt.Println(samuraiSudoku.Grid())

	SolveSamuraiSudoku(&samuraiSudoku)

	fmt.Println(samuraiSudoku.Grid())
}
