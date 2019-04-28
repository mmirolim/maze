package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

// PopulateSudokuBoard with digits according to initial digits provided
// by initPos
func PopulateSudokuBoard(initPos [9][9]int) ([9][9]int, error) {
	result := [9][9]int{} // result matrix
	var i, j int          // slot coord in a box
	var candidate int
	digitProvidedH := [9][10]bool{}
	digitProvidedV := [9][10]bool{}
	digitProvidedBox := [9][10]bool{}
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			d := initPos[x][y]
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
		return initPos[i][j]
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
			setCandidate(initPos[i][j])
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
				return result, fmt.Errorf("no solution found for init pos %v", initPos)
			}
		}

	}

	return result, nil
}

// DrawSudokuBoard ...
func DrawSudokuBoard(cellSize int, num [9][9]int) (*image.RGBA, error) {
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
					d.DrawString(strconv.Itoa(num[i*3+x][j*3+y]))

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
