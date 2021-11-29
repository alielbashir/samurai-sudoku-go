package main

import (
	"fmt"
	. "github.com/alielbashir/samurai-sudoku-go"
)

func main() {
	samuraiGrid := SamuraiGridFromFile("sudoku.txt")

	var samuraiSudoku SamuraiSudoku

	samuraiSudoku.SetGrid(samuraiGrid)

	fmt.Println(samuraiSudoku.Grid())

	ConcurrentSolveSamuraiSudoku(&samuraiSudoku)

	fmt.Println(samuraiSudoku.Grid())
}
