package main

import "fmt"

func printSudoku(sudoku [][]int) {
	fmt.Println()
	for _, row := range sudoku {
		for _, num := range row {
			fmt.Print(num, " ")
		}
		fmt.Println()
	}
}

//func printSamuraiSudoku(samurai [][]int) {
//	fmt.Println()
//	for _, row := range samurai {
//		for _, num := range row {
//			if num == -1 {
//				fmt.Print("  ")
//			} else {
//				fmt.Print(num, " ")
//			}
//		}
//		fmt.Println()
//	}
//}

func possible(sudoku [][]int, y int, x int, n int) bool {
	for i := 0; i < 9; i++ {
		if sudoku[y][i] == n {
			return false
		}
	}
	for i := 0; i < 9; i++ {
		if sudoku[i][x] == n {
			return false
		}
	}
	x0 := (x / 3) * 3
	y0 := (y / 3) * 3

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if sudoku[y0+i][x0+j] == n {
				return false
			}
		}
	}
	return true
}

//func solveSudoku(bigSudoku [][]int) {
//
//}

func backtrack(sudoku [][]int) bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			// if cell is empty
			if sudoku[y][x] == 0 {
				for n := 1; n < 10; n++ {
					if possible(sudoku, y, x, n) {
						sudoku[y][x] = n
						if backtrack(sudoku) {
							return true
						}
						sudoku[y][x] = 0
					}
				}
				return false
			}
		}
	}
	fmt.Println("Solved!")
	printSudoku(sudoku)
	return true
}

func main() {
	var sudoku = [][]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 0, 0},
	}
	printSudoku(sudoku)
	backtrack(sudoku)
	//printSudoku()
}
