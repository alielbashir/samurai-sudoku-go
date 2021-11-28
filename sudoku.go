package sudoku

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var logger = log.New(os.Stdout, "", 0)

type Grid [][]int

//isSolved tells if this sudoku has been solved or not
//works for both 9x9 sudokus and 21*21 samurai sudokus
func (g Grid) isSolved() bool {
	for _, row := range g {
		for _, num := range row {
			if num == 0 {
				return false
			}
		}
	}
	return true
}

type SamuraiSudoku struct {
	mu   sync.Mutex
	grid Grid
}

func (s *SamuraiSudoku) Grid() Grid {
	return s.grid
}

func (s *SamuraiSudoku) SetGrid(grid Grid) {
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

func (g Grid) String() string {
	var buf bytes.Buffer
	var char string
	for _, row := range g {
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
func (s *SamuraiSudoku) GetSubSudoku(position Position) Grid {
	var grid = s.grid
	subSudoku := make(Grid, 9)
	var tmp Grid
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

//SolveSamuraiSudoku solves 21*21 samurai sudoku
func SolveSamuraiSudoku(samurai *SamuraiSudoku) Grid {

	// get all subsudokus
	subSudokus := []struct {
		position Position
		sudoku   Grid
	}{
		{TopLeft, samurai.GetSubSudoku(TopLeft)},
		{TopRight, samurai.GetSubSudoku(TopRight)},
		{Centre, samurai.GetSubSudoku(Centre)},
		{BottomLeft, samurai.GetSubSudoku(BottomLeft)},
		{BottomRight, samurai.GetSubSudoku(BottomRight)},
	}

	// iterate over the map until all subsudokus are solved
	for _, subSudoku := range subSudokus {
		solution := SolveSudoku(subSudoku.sudoku, subSudoku.position, samurai)
		logger.Printf("%s\n%v\n", subSudoku.position, solution)
	}

	return samurai.Grid()
}

//ConcurrentSolveSamuraiSudoku solves 21*21 samurai sudoku concurrently
func ConcurrentSolveSamuraiSudoku(samurai *SamuraiSudoku) Grid {

	// get all subsudokus
	subSudokus := []struct {
		position Position
		sudoku   Grid
	}{
		{TopLeft, samurai.GetSubSudoku(TopLeft)},
		{TopRight, samurai.GetSubSudoku(TopRight)},
		{Centre, samurai.GetSubSudoku(Centre)},
		{BottomLeft, samurai.GetSubSudoku(BottomLeft)},
		{BottomRight, samurai.GetSubSudoku(BottomRight)},
	}

	wg := new(sync.WaitGroup)

	// iterate over the map until all subsudokus are solved
	solvingLoop(samurai, subSudokus, wg)

	return samurai.Grid()
}

func solvingLoop(samurai *SamuraiSudoku, subSudokus []struct {
	position Position
	sudoku   Grid
}, wg *sync.WaitGroup) {
	for _, subSudoku := range subSudokus {
		// increment WaitGroup counter
		subSudoku := subSudoku
		wg.Add(1)
		go concurrentSolveSudoku(subSudoku.sudoku, subSudoku.position, samurai, wg)
	}
	wg.Wait()
}

//possible checks if index y,x in grid position can be filled with n in all subsudokus it's in
func possible(sudoku Grid, y int, x int, n int, position Position, samuraiSudoku *SamuraiSudoku) bool {
	var sharedSudoku Grid
	var yShared, xShared int
	switch position {
	case TopLeft:
		if 6 <= y && y < 9 && 6 <= x && x < 9 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(Centre)
			yShared, xShared = y-6, x-6
		}
	case TopRight:
		if 6 <= y && y < 9 && 0 <= x && x < 3 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(Centre)
			yShared, xShared = y-6, x+6
		}

	case Centre:
		if 0 <= y && y < 3 && 0 <= x && x < 3 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(TopLeft)
			yShared, xShared = y+6, x+6
		} else if 0 <= y && y < 3 && 6 <= x && x < 9 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(TopRight)
			yShared, xShared = y+6, x-6
		} else if 6 <= y && y < 9 && 0 <= x && x < 3 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(BottomLeft)
			yShared, xShared = y-6, x+6
		} else if 6 <= y && y < 9 && 6 <= x && x < 9 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(BottomRight)
			yShared, xShared = y-6, x-6
		}

	case BottomLeft:
		if 0 <= y && y < 3 && 6 <= x && x < 9 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(Centre)
			yShared, xShared = y+6, x-6
		}
	case BottomRight:
		if 0 <= y && y < 3 && 0 <= x && x < 3 {
			sharedSudoku = samuraiSudoku.GetSubSudoku(Centre)
			yShared, xShared = y+6, x+6
		}

	}
	if len(sharedSudoku) == 0 {
		return possibleSudoku(sudoku, y, x, n)
	} else {
		return possibleSudoku(sudoku, y, x, n) && possibleSudoku(sharedSudoku, yShared, xShared, n)
	}
}

//possibleSudoku checks if sudoku can be filled in position y,x with n
func possibleSudoku(sudoku Grid, y int, x int, n int) bool {
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

//concurrentSolveSudoku solves 9x9 subsudoku in specified position within samuraiSudoku, concurrently
func concurrentSolveSudoku(sudoku Grid, position Position, samuraiSudoku *SamuraiSudoku, wg *sync.WaitGroup) Grid {
	// TODO: fix some sudokus not solving.
	for !sudoku.isSolved() {
		time.Sleep(time.Duration(time.Duration.Milliseconds(20000)))
		wg.Add(1)
		logger.Printf("%s attempting solve on %v\n", position, sudoku)
		samuraiSudoku.mu.Lock()
		//logger.Printf("%s routine acquired lock\n", position)
		backtrack(sudoku, position, samuraiSudoku)
		samuraiSudoku.mu.Unlock()
		//logger.Printf("%s routine released lock\n", position)
		wg.Done()
	}

	logger.Printf("%s\n%v\n", position, sudoku)
	wg.Done()

	return sudoku
}

//SolveSudoku solves 9x9 subsudoku in position within samuraiSudoku
func SolveSudoku(sudoku Grid, position Position, samuraiSudoku *SamuraiSudoku) Grid {
	backtrack(sudoku, position, samuraiSudoku)
	return sudoku
}

//backtrack keeps attempting values recursively until 9x9 sudoku is solved completely
func backtrack(sudoku Grid, position Position, samuraiSudoku *SamuraiSudoku) bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			// if cell is empty
			if sudoku[y][x] == 0 {
				for n := 1; n < 10; n++ {
					if possible(sudoku, y, x, n, position, samuraiSudoku) {
						sudoku[y][x] = n
						if backtrack(sudoku, position, samuraiSudoku) {
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
