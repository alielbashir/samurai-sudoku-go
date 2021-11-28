package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

type Matrix [][]int

type SamuraiSudoku struct {
	mu   sync.Mutex
	grid Matrix
}

func (s *SamuraiSudoku) Grid() Matrix {
	return s.grid
}

func (s *SamuraiSudoku) SetGrid(grid Matrix) {
	s.grid = grid
}

type Position int

const (
	TopLeft Position = iota
	TopRight
	Centre
	BottomLeft
	BottomRight
)

func (p Position) String() string {
	switch p {
	case TopLeft:
		return "top left"
	case TopRight:
		return "top right"
	case Centre:
		return "centre"
	case BottomLeft:
		return "bottom left"
	case BottomRight:
		return "bottom right"
	}
	return "unknown"
}

func (m Matrix) String() string {
	var buf bytes.Buffer
	var char string
	for _, row := range m {
		for _, num := range row {
			if num == -1 {
				char = " "
			} else {
				char = strconv.Itoa(num)
			}
			_, err := fmt.Fprint(&buf, char, " ")
			if err != nil {
				return ""
			}
		}
		_, err := fmt.Fprint(&buf, "\n")
		if err != nil {
			return ""
		}
	}
	return buf.String()
}

//GetSubSudoku returns sub-sudoku for given position, assuming 21*21 samurai sudoku grid
func (s *SamuraiSudoku) GetSubSudoku(position Position) Matrix {
	var grid = s.grid
	subSudoku := make(Matrix, 9)
	var tmp Matrix
	switch position {
	case TopLeft:
		tmp = grid[0:9]
		for i := range tmp {
			subSudoku[i] = tmp[i][0:9]
		}
	case TopRight:
		tmp = grid[0:9]
		for i := range tmp {
			subSudoku[i] = tmp[i][12:21]
		}
	case Centre:
		tmp = grid[6:15]
		for i := range tmp {
			subSudoku[i] = tmp[i][6:15]
		}

	case BottomLeft:
		tmp = grid[12:21]
		for i := range tmp {
			subSudoku[i] = tmp[i][0:9]
		}

	case BottomRight:
		tmp = grid[12:21]
		for i := range tmp {
			subSudoku[i] = tmp[i][12:21]
		}
	}
	return subSudoku
}

func SolveSamuraiSudoku(samurai *SamuraiSudoku) Matrix {
	// get all subsudokus
	subSudokus := []Matrix{
		samurai.GetSubSudoku(TopLeft),
		samurai.GetSubSudoku(TopRight),
		samurai.GetSubSudoku(Centre),
		samurai.GetSubSudoku(BottomLeft),
		samurai.GetSubSudoku(BottomRight),
	}

	// iterate over the slice until everything is solved
	solutions := make([]Matrix, 5)
	// TODO: check possible for all shared sudokus
	for _, sudoku := range subSudokus {
		solutions = append(solutions, SolveSudoku(sudoku))
	}

	for _, solution := range solutions {
		fmt.Println(solution)
	}
	return samurai.Grid()
}

func possible(sudoku Matrix, y int, x int, n int) bool {
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

func SolveSudoku(sudoku Matrix) Matrix {
	backtrack(sudoku)
	return sudoku
}

func backtrack(sudoku Matrix) bool {
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
	return true
}