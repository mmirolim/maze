package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

// Sudoku board
type Sudoku struct {
	initPos [9][9]int
	result  [9][9]int
}

// NewSudoku returns sudoku with initPos set
func NewSudoku(initPos [9][9]int) *Sudoku {
	sudoku := new(Sudoku)
	sudoku.initPos = initPos
	return sudoku
}

// GetSolution returns sudoku solution matrix
func (s *Sudoku) GetSolution() [9][9]int {
	solution := [9][9]int{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			solution[i][j] = s.result[i][j]
		}
	}
	return solution
}

// SolveWithBacktracking ---
func (s *Sudoku) SolveWithBacktracking() error {
	result := [9][9]int{} // result matrix
	var i, j int          // slot coord in a box
	var candidate int
	digitProvidedH := [9][10]bool{}
	digitProvidedV := [9][10]bool{}
	digitProvidedBox := [9][10]bool{}
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			d := s.initPos[x][y]
			if d > 0 {
				digitProvidedH[y][d] = true
				digitProvidedV[x][d] = true
				digitProvidedBox[(x/3)*3+y/3][d] = true
			}
		}

	}

	availDigitsH := [9][10]bool{}
	availDigitsV := [9][10]bool{}
	availDigitBox := [9][10]bool{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 10; j++ {
			availDigitBox[i][j] = true
			availDigitsH[i][j] = true
			availDigitsV[i][j] = true
		}
	}

	advanceSlot := func() {
		j++
	}
	nextColumn := func() {
		i++
		j = 0
	}

	digitProvided := func(i, j int) int {
		return s.initPos[i][j]
	}

	testH := func(d int) bool {
		return availDigitsH[j][d] && !digitProvidedH[j][d]
	}

	testV := func(d int) bool {
		return availDigitsV[i][d] && !digitProvidedV[i][d]
	}

	testInBox := func(d int) bool {
		index := (i/3)*3 + j/3
		return availDigitBox[index][d] && !digitProvidedBox[index][d]
	}

	testDigit := func(d int) bool {
		return testH(d) && testV(d) && testInBox(d)
	}

	setCandidate := func(d int) {
		candidate = d
	}

	testSlot := func() bool {
		if digitProvided(i, j) > 0 {
			setCandidate(s.initPos[i][j])
			return true
		}

		val := candidate
		for val < 9 {
			val++
			ok := testDigit(val)
			if ok {
				setCandidate(val)
				return true
			}
		}
		return false
	}

	priorSlot := func() {
		j--
	}
	priorColumn := func() {
		i--
		j = 8
	}

	startColumn := func() {
		i = 0
	}
	lastSlotDone := func() bool {
		return i == 8 && j > 8
	}
	lastSlotInColumnDone := func() bool {
		return j > 8
	}
	isFirstSlotInColumn := func() bool {
		return j == 0
	}

	noRegressSlotLeft := func() bool {
		return i < 0
	}
	regressSlot := func() {
		for {
			if isFirstSlotInColumn() {
				priorColumn()
			} else {
				priorSlot()
			}
			if !noRegressSlotLeft() && digitProvided(i, j) > 0 {
				continue
			} else {
				break
			}
		}
	}
	setDigit := func() {
		availDigitsH[j][candidate] = false
		availDigitsV[i][candidate] = false
		availDigitBox[(i/3)*3+j/3][candidate] = false
		result[i][j] = candidate
		setCandidate(0)
	}

	removeDigit := func() {
		setCandidate(result[i][j])
		availDigitsH[j][candidate] = true
		availDigitsV[i][candidate] = true
		availDigitBox[(i/3)*3+j/3][candidate] = true
		result[i][j] = 0
	}
	startColumn()
	for {
		if testSlot() {
			setDigit()
			advanceSlot()
			if lastSlotDone() {
				break
			}
			if lastSlotInColumnDone() {
				nextColumn()
			}

		} else {
			regressSlot()
			removeDigit()
			if noRegressSlotLeft() {
				return fmt.Errorf("no solution found for init pos %v", s.initPos)
			}
		}

	}

	s.result = result
	return nil
}

// SolveWithSA solves sudoku with Simulated Annealing metaheuristics
func (s *Sudoku) SolveWithSA() error {
	state := [9][9]int{}

	// isFixedValueCell := func(i, j int) bool {
	// 	return s.initPos[i][j] > 0
	// }
	randomValuesExcluding := func(val []int) []int {
		out := []int{}
		r := rand.Perm(9)
	LOOP:
		for i := range r {
			for _, v := range val {
				if r[i]+1 == v {
					continue LOOP
				}
			}
			out = append(out, r[i]+1)
		}
		return out
	}

	fixedValuesFromSquare := func(i, j int) []int {
		r := []int{}
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				val := s.initPos[i*3+x][j*3+y]
				if val > 0 {
					r = append(r, val)
				}
			}
		}
		return r
	}

	// generate initial state
	// square 3x3 contains the integers 1 through to 9 exactly once
	fillSquare := func(i, j int) {
		digits := randomValuesExcluding(fixedValuesFromSquare(i, j))
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				val := s.initPos[i*3+x][j*3+y]
				if val == 0 {
					val = digits[len(digits)-1]
					digits = digits[:len(digits)-1]
				}

				state[i*3+x][j*3+y] = val

			}
		}
	}

	generateInitState := func() {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				fillSquare(i, j)
			}
		}
	}

	generateInitState()

	// E state energy function, cost function
	// calculate number of not present numbers in each
	// row and column, E is sum(Er, Ec)
	const uint16Max uint16 = 65535
	rowE := [9]uint16{uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max}
	columnE := [9]uint16{uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max, uint16Max}

	E := func() int {
		e := 0

		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				v := state[x][y]
				columnE[x] &^= 1 << uint16(v-1)
				rowE[y] &^= 1 << uint16(v-1)
			}
		}
		for x := range columnE {
			for i := 0; i < 9; i++ {
				if columnE[x]&(1<<uint16(i)) != 0 {
					e++
				}
				if rowE[x]&(1<<uint16(i)) != 0 {
					e++
				}
			}
			columnE[x] = uint16Max
			rowE[x] = uint16Max
		}

		return e
	}

	isFixed := func(i, j int) bool {
		return s.initPos[i][j] > 0
	}
	type cell struct {
		i, j int
	}
	var cell1, cell2 cell
	listOfNonFixedCells := func() []cell {
		r := []cell{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if !isFixed(i, j) {
					r = append(r, cell{i, j})
				}
			}
		}
		return r
	}()

	listOfNonFixedNeighoursInSq := func(i, j int) []cell {
		r := []cell{}
		sqx, sqy := (i/3)*3, (j/3)*3
		for x := sqx; x < sqx+3; x++ {
			for y := sqy; y < sqy+3; y++ {
				if isFixed(x, y) || (x == i && y == j) {
					continue
				}
				r = append(r, cell{x, y})
			}

		}
		return r
	}

	// neighbour candidate generator procedure
	// mutate state to new state
	neighbour := func() {
		cell1 = listOfNonFixedCells[rand.Intn(len(listOfNonFixedCells))]
		list2 := listOfNonFixedNeighoursInSq(cell1.i, cell1.j)
		cell2 = list2[rand.Intn(len(list2))]
		// swap cells
		state[cell1.i][cell1.j], state[cell2.i][cell2.j] = state[cell2.i][cell2.j], state[cell1.i][cell1.j]
	}
	revert := func() {
		state[cell1.i][cell1.j], state[cell2.i][cell2.j] = state[cell2.i][cell2.j], state[cell1.i][cell1.j]

	}
	// P acceptance probability function
	P := func(e1, e2 int, t float64) float64 {
		if e2 < e1 {
			return 1.0
		}
		return math.Exp(-float64(e2-e1) / t)
	}

	alpha := 0.9
	t0 := 1.0
	// T temperature non increasing function
	T := func(t float64) float64 {
		return alpha * t
	}
	var err error
	Tmin := 0.0001
	t := t0
	Elast := E()
	var Enew int

	// TODO maybe use restarts?
	// TODO choose X
	X := len(listOfNonFixedCells)
LOOP:
	for {
		// update temp
		t = T(t)
		for i := 0; i < X; i++ {
			// create new state
			neighbour()
			// compute new state Energy
			Enew = E()
			if Enew == 0 {
				// solution found
				break LOOP
			}
			// test transition probability
			if P(Elast, Enew, t) > rand.Float64() {
				Elast = Enew
			} else {
				// revert neighbour mutation
				revert()
			}

		}
		// stop if temp too low
		if Tmin > t {
			err = errors.New("temp limit reached")
			break
		}

	}

	s.result = state
	return err
}

// PrintBoard prints board to stdout
func (s *Sudoku) PrintBoard() {
	printState(s.result)
}

func printState(s [9][9]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%v ", s[j][i]) // output for debug
		}
		fmt.Println("")

	}
}

// Draw SudokuBoard with cellSize in px
func (s *Sudoku) Draw(cellSize int) (*image.RGBA, error) {
	black := color.RGBA{0, 0, 0, 255}
	cellPadding := cellSize / 20
	box3x3Padding := cellSize / 20
	box3x3Size := cellSize*3 + 2*box3x3Padding + 2*cellPadding
	// draw board
	board := image.NewRGBA(image.Rect(0, 0, box3x3Size*3, box3x3Size*3))
	draw.Draw(board, board.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	box3x3 := image.NewRGBA(image.Rect(0, 0, box3x3Size, box3x3Size))
	cell := image.NewRGBA(image.Rect(0, 0, cellSize, cellSize))
	// number's font
	fontVal, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	fontOptionSize := cellSize / 2
	fontFace := truetype.NewFace(fontVal, &truetype.Options{Size: float64(fontOptionSize)})
	d := &font.Drawer{
		Dst:  cell,
		Src:  image.NewUniform(black),
		Face: fontFace,
		Dot:  fixed.Point26_6{},
	}

	digitPos := fixed.P(cellSize/2-fontOptionSize/4, cellSize/2+fontOptionSize/3)
	// draw boxes on a board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// draw box
			draw.Draw(box3x3, box3x3.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
			// draw cells in a box
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					// draw cell
					draw.Draw(cell, cell.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
					// draw number
					d.Dot = digitPos
					d.DrawString(strconv.Itoa(s.result[i*3+x][j*3+y]))

					draw.Draw(
						box3x3,
						cell.Bounds().Add(image.Point{X: x * (cellSize + cellPadding), Y: y * (cellSize + cellPadding)}),
						cell,
						image.ZP,
						draw.Src,
					)

				}
			}
			draw.Draw(board,
				box3x3.Bounds().Add(
					image.Point{X: i * (box3x3Size + box3x3Padding), Y: j * (box3x3Size + box3x3Padding)},
				),
				box3x3, image.ZP, draw.Src)
		}
	}

	return board, nil
}
