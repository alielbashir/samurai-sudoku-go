package main

import (
	. "github.com/alielbashir/samurai-sudoku-go"
)

func main() {
	samuraiGrid := SamuraiGridFromFile("sudoku.txt")

	var samuraiSudoku SamuraiSudoku

	samuraiSudoku.SetGrid(samuraiGrid)

	DoubleThreadSolveSamuraiSudoku(&samuraiSudoku)
	WriteGraph(&samuraiSudoku)

}
