package sudoku

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

// Move A single move in sudoku
type Move struct {
	thread   int
	position Position // Position of the sudoku this Move was done in
	row      int
	column   int
	num      int // Number inserted
	time     int64
}

func (m Move) String() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "%d,%d, %s,%d,%d,%d", m.time, m.thread, m.position, m.row, m.column, m.num)
	return buf.String()
}

type Tracker struct {
	moves []Move
}

func (t *Tracker) resetMoves() {
	t.moves = nil
}

//SamuraiGridFromFile reads a samurai sudoku grid from a given file
func SamuraiGridFromFile(filePath string) Grid {
	const samuraiLength = 21
	grid := make(Grid, samuraiLength)

	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Printf("File couldn't be read!")
	}

	sudokuContents := string(buffer)

	for j, line := range strings.Split(sudokuContents, "\n") {
		charRow := strings.Split(line, "")
		intRow := make([]int, samuraiLength, samuraiLength)
		offset := 0
		for i := 0; i < samuraiLength; i++ {
			switch len(charRow) {
			case 9:
				if i < 6 || 15 <= i {
					intRow[i] = -1
					offset++
				} else {
					num, _ := strconv.Atoi(charRow[i-offset])

					intRow[i] = num
				}
			case 18:
				if 9 <= i && i < 12 {
					intRow[i] = -1
					offset++
				} else {
					num, _ := strconv.Atoi(charRow[i-offset])
					intRow[i] = num
				}

			default:
				num, _ := strconv.Atoi(charRow[i-offset])
				intRow[i] = num
			}
		}

		grid[j] = intRow
	}
	logger.Printf("Read \n%v\n", grid)
	return grid
}

type SamuraiSudoku struct {
	mu          sync.Mutex
	grid        Grid
	initialGrid Grid
	tracker     Tracker
}

func (s *SamuraiSudoku) ResetGrid() {
	s.tracker.resetMoves()
	for i, row := range s.initialGrid {
		for j, num := range row {
			s.grid[i][j] = num
		}
	}
}

func (s *SamuraiSudoku) Grid() Grid {
	return s.grid
}

func (s *SamuraiSudoku) SetGrid(grid Grid) {
	if s.initialGrid == nil {
		s.initialGrid = make(Grid, len(grid))
		for i := range grid {
			s.initialGrid[i] = make([]int, len(grid[i]))
			copy(s.initialGrid[i], grid[i])
		}
	}
	s.grid = grid
}

type Position int

const (
	TopLeft Position = iota + 1
	TopRight
	Centre
	BottomLeft
	BottomRight
)

type ThreadId int

const (
	Thread1 ThreadId = iota + 1
	Thread2
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

func (s *SamuraiSudoku) recordMove(id ThreadId, position Position, y int, x int, n int) {
	s.tracker.moves = append(s.tracker.moves, Move{
		thread:   int(position) * int(id),
		position: position,
		row:      y,
		column:   x,
		num:      n,
		time:     time.Now().UnixMicro(),
	})
}

func (s *SamuraiSudoku) moves() bytes.Buffer {
	buf := bytes.Buffer{}
	for _, move := range s.tracker.moves {
		fmt.Fprintf(&buf, "%s\n", move.String())
	}
	return buf
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
		SolveSudoku(subSudoku.sudoku, subSudoku.position, samurai)
	}

	return samurai.Grid()
}

//ConcurrentSolveSamuraiSudoku solves 21*21 samurai sudoku concurrently
func ConcurrentSolveSamuraiSudoku(samurai *SamuraiSudoku) Grid {
	rand.Seed(time.Now().UnixNano())
	// get all subsudokus
	getSubSudokus := func() []struct {
		position Position
		sudoku   Grid
	} {
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
		rand.Shuffle(len(subSudokus), func(i, j int) {
			subSudokus[i], subSudokus[j] = subSudokus[j], subSudokus[i]
		})
		return subSudokus
	}

	wg := new(sync.WaitGroup)

	// iterate over the map until all subsudokus are solved
	for !samurai.Grid().isSolved() {
		samurai.mu.Lock()
		samurai.ResetGrid()
		subSudokus := getSubSudokus()
		// reset samurai grid
		samurai.mu.Unlock()
		solvingLoop(samurai, subSudokus, wg)
		SolvingAttempts++
	}

	moves := samurai.moves()
	os.WriteFile("sudoku.log", moves.Bytes(), 0666)
	logger.Printf("attempt %d\n%v\n", SolvingAttempts, samurai.Grid())

	return samurai.Grid()
}

//DoubleThreadSolveSamuraiSudoku solves 21*21 samurai sudoku concurrently
func DoubleThreadSolveSamuraiSudoku(samurai *SamuraiSudoku) Grid {
	rand.Seed(time.Now().UnixNano())
	// get all subsudokus
	getSubSudokus := func() []struct {
		position Position
		sudoku   Grid
	} {
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
		rand.Shuffle(len(subSudokus), func(i, j int) {
			subSudokus[i], subSudokus[j] = subSudokus[j], subSudokus[i]
		})
		return subSudokus
	}

	wg := new(sync.WaitGroup)

	// iterate over the map until all subsudokus are solved
	for !samurai.Grid().isSolved() {
		samurai.mu.Lock()
		samurai.ResetGrid()
		subSudokus := getSubSudokus()
		// reset samurai grid
		samurai.mu.Unlock()
		doubleSolvingLoop(samurai, subSudokus, wg)
		SolvingAttempts++
	}

	moves := samurai.moves()
	os.WriteFile("sudoku.log", moves.Bytes(), 0666)
	logger.Printf("attempt %d\n%v\n", SolvingAttempts, samurai.Grid())

	return samurai.Grid()
}

var SolvingAttempts = 0

func solvingLoop(samurai *SamuraiSudoku, subSudokus []struct {
	position Position
	sudoku   Grid
}, wg *sync.WaitGroup) {
	wg.Add(len(subSudokus))
	for _, subSudoku := range subSudokus {
		// increment WaitGroup counter
		go concurrentSolveSudoku(Thread1, subSudoku.sudoku, subSudoku.position, samurai, wg)
	}
	wg.Wait()
	var order bytes.Buffer
	// populate order buffer for debugging purposes
	for _, sudokus := range subSudokus {
		_, err := fmt.Fprintf(&order, "%d, ", sudokus.position)
		if err != nil {
			return
		}
	}
}

func doubleSolvingLoop(samurai *SamuraiSudoku, subSudokus []struct {
	position Position
	sudoku   Grid
}, wg *sync.WaitGroup) {
	wg.Add(len(subSudokus) * 2)
	for _, subSudoku := range subSudokus {
		// increment WaitGroup counter
		go concurrentSolveSudoku(Thread1, subSudoku.sudoku, subSudoku.position, samurai, wg)
		go concurrentSolveSudoku(Thread2, subSudoku.sudoku, subSudoku.position, samurai, wg)
	}
	wg.Wait()
	var order bytes.Buffer
	// populate order buffer for debugging purposes
	for _, sudokus := range subSudokus {
		_, err := fmt.Fprintf(&order, "%d, ", sudokus.position)
		if err != nil {
			return
		}
	}
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
func concurrentSolveSudoku(threadId ThreadId, sudoku Grid, position Position, samuraiSudoku *SamuraiSudoku, wg *sync.WaitGroup) Grid {
	// TODO: fix some sudokus not solving.
	if sudoku.isSolved() {
		wg.Done()
		return sudoku
	}

	backtrack(threadId, sudoku, position, samuraiSudoku)

	wg.Done()

	return sudoku
}

//SolveSudoku solves 9x9 subsudoku in position within samuraiSudoku
func SolveSudoku(sudoku Grid, position Position, samuraiSudoku *SamuraiSudoku) Grid {
	backtrack(Thread1, sudoku, position, samuraiSudoku)
	return sudoku
}

//backtrack keeps attempting values recursively until 9x9 sudoku is solved completely
func backtrack(threadId ThreadId, sudoku Grid, position Position, samuraiSudoku *SamuraiSudoku) bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			// if cell is empty
			//logger.Printf("%s: waiting for lock...", position)
			samuraiSudoku.mu.Lock()
			//logger.Printf("%s: locked", position)
			if sudoku[y][x] == 0 {
				for n := 1; n < 10; n++ {
					if possible(sudoku, y, x, n, position, samuraiSudoku) {
						samuraiSudoku.recordMove(threadId, position, y, x, n)
						sudoku[y][x] = n
						samuraiSudoku.mu.Unlock()
						//logger.Printf("%s: set sudoku[%d, %d] = %d", position, y, x, n)
						if backtrack(threadId, sudoku, position, samuraiSudoku) {
							// should be unlocked here, but could get locked by other threads
							return true
						}
						//logger.Printf("%s: waiting for lock for 0", position)
						samuraiSudoku.mu.Lock()
						//logger.Printf("%s: acquired lock for 0", position)
						samuraiSudoku.recordMove(threadId, position, y, x, 0)
						sudoku[y][x] = 0
						//logger.Printf("%s: releasing lock after 0", position)
						//samuraiSudoku.mu.Unlock()
					}
				}
				//logger.Printf("%s: returning false, %d %d", position, y, x)
				samuraiSudoku.mu.Unlock()
				return false
			}
			//logger.Printf("%s: released lock, %d, %d", position, y, x)
			samuraiSudoku.mu.Unlock()
		}
	}
	//samuraiSudoku.mu.Unlock()
	return true
}
