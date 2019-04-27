package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const queenImg = "images/queen.png"

type ChessPieceNames string

var (
	Queen ChessPieceNames = "Queen"
)

type ChessPiece struct {
	name ChessPieceNames
	img  image.Image
	pos  string
}

func NewChessPiece(name ChessPieceNames, pos string) (*ChessPiece, error) {
	var imageName string
	switch name {
	case Queen:
		imageName = queenImg
	}
	f, err := os.Open(imageName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	if len(pos) != 2 {
		return nil, fmt.Errorf("wrong position format %+v", pos)
	}

	return &ChessPiece{
		name: name,
		img:  img,
		pos:  pos,
	}, nil

}

type ChessBoard struct {
	cellSize int
	pieces   []ChessPiece
}

func NewChessBoard(cellSize int, pieces []ChessPiece) *ChessBoard {
	return &ChessBoard{
		cellSize: cellSize,
		pieces:   pieces,
	}
}

func (b *ChessBoard) Draw() (*image.RGBA, error) {
	bgColor := color.RGBA{242, 244, 247, 255}
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	size := b.cellSize * 8
	border := b.cellSize * 2
	m := image.NewRGBA(image.Rect(0, 0, size+border, size+border))
	board := image.NewRGBA(image.Rect(0, 0, size, size)) // 8 cells  + border
	draw.Draw(m, m.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)
	var cellColor color.RGBA
	// TODO implement image interface for chessboard struct
	for x := 1; x < 9; x++ {
		for y := 1; y < 9; y++ {
			if (x+y)%2 == 0 {
				cellColor = white
			} else {
				cellColor = black
			}
			for i := 0; i < b.cellSize; i++ {
				for j := 0; j < b.cellSize; j++ {
					board.SetRGBA((x-1)*b.cellSize+i, (y-1)*b.cellSize+j, cellColor)
				}
			}
		}
	}

	draw.Draw(m,
		board.Bounds().Add(image.Point{X: b.cellSize, Y: b.cellSize}),
		board, image.ZP, draw.Src)

	// returns x and y coords of cell
	piecePosition := func(p *ChessPiece) (int, int) {
		x := int(p.pos[0] - 'a' + 1)
		y := int(p.pos[1] - '0')
		if x < 0 || x > 8 || y < 0 || y > 8 {
			panic(fmt.Sprintf("wrong piece pos format %+v", p.pos))
		}

		return x, 8 - y + 1
	}
	// draw pieces on the board
	for _, piece := range b.pieces {
		scaledPieceImg := image.NewRGBA(image.Rect(0, 0, b.cellSize, b.cellSize))
		draw.BiLinear.Scale(
			scaledPieceImg, scaledPieceImg.Bounds(), piece.img, piece.img.Bounds(),
			draw.Over, nil,
		)

		x, y := piecePosition(&piece)
		draw.Draw(m,
			scaledPieceImg.Bounds().Add(image.Point{X: x * b.cellSize, Y: y * b.cellSize}),
			scaledPieceImg, image.ZP, draw.Over)
	}
	// draw legends
	d := &font.Drawer{
		Dst:  m,
		Src:  image.NewUniform(black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{},
	}
	cellSizeHalf := b.cellSize / 2
	Ycoords := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	for k, v := range Ycoords {
		// draw numbers
		d.Dot = fixed.P(b.cellSize-10, (9-k)*b.cellSize-cellSizeHalf)
		d.DrawString(strconv.Itoa(k + 1))
		// draw letters
		d.Dot = fixed.P((k+2)*b.cellSize-cellSizeHalf, b.cellSize*9+10)
		d.DrawString(v)
	}

	return m, nil
}

// TheEightQueensProblemMySolution to the 8-queen puzzle
// “The eight queens puzzle is the problem
// of placing eight chess queens on an 8×8 chessboard
// so that no two queens threaten each other.
// Thus, a solution requires that no two queens
// share the same row, column, or diagonal.”
func TheEightQueensProblemMySolution() [][2]int {
	// auxiliary structures
	isNotFreeRow := [9]bool{}
	isNotFreeDig1 := [17]bool{}
	isNotFreeDig2 := [15]bool{}

	pieces := []int{}
	var search func(i, j int) bool
	search = func(i, j int) bool {
		if j > 8 {
			return false
		}
		// end of chessboard
		if i > 8 {
			return true
		}

		if !isNotFreeRow[j] && !isNotFreeDig1[i+j] && !isNotFreeDig2[j-i+7] {
			// valid position
			isNotFreeRow[j] = true
			isNotFreeDig1[i+j] = true
			isNotFreeDig2[j-i+7] = true
			pieces = append(pieces, j)
			i++
			j = 1
		} else {
			j++
		}

		if !search(i, j) {
			// on err take last element
			// and continue from that position with shift
			i := len(pieces)
			j := pieces[len(pieces)-1]
			isNotFreeRow[j] = false
			isNotFreeDig1[i+j] = false
			isNotFreeDig2[j-i+7] = false
			pieces = pieces[:len(pieces)-1]
			return search(i, j+1)
		}

		return true
	}

	// start search
	_ = search(1, 1)

	sl := [][2]int{}
	for i := range pieces {
		sl = append(sl, [2]int{i + 1, pieces[i]})
	}

	return sl
}

// TheEightQueensProblemNW solution to 8-queen puzzle
// described by Niklaus Wirth in
// http://plbpc001.ouhk.edu.hk/~mt311/optional-reading/stepwise.pdf
// returns slice of values where index is column and value is row
func TheEightQueensProblemNW() []int {
	// current column
	var j int // 0 <= j <= 9

	// solutions of rows to place queens
	// size defined for max j value
	var x = make([]int, 10) // vals 0 <= v <= 8
	i := 0                  // data to store x[j] for efficiency
	// predicates
	var safe bool
	lastSquare := func() bool {
		return i == 8
	}

	lastColDone := func() bool {
		return j > 8
	}

	regressOutOfFirstCol := func() bool {
		return j < 1
	}
	// instructions
	considerFirstColumn := func() {
		j = 1
		i = 0
	}
	considerNextColumn := func() {
		x[j] = i
		j = j + 1
		i = 0
	}
	reconsiderPriorColumn := func() {
		j = j - 1
		i = x[j]
	}
	advancePointer := func() {
		i = i + 1
	}

	// auxiliary variables for efficience of testSquare
	rowIsFree := make([]bool, 9) // vals 1:8 number of rows
	for i := range rowIsFree {
		rowIsFree[i] = true
	}
	// /-diagonal is free
	diagonal1IsFree := make([]bool, 17) // vals 2:16, sum of coords
	for i := range diagonal1IsFree {
		diagonal1IsFree[i] = true
	}
	// \-diagonal is free
	diagonal2IsFree := make([]bool, 15) // values -7:7 diff of coords, 7 should be added on index check
	for i := range diagonal2IsFree {
		diagonal2IsFree[i] = true
	}

	testSquare := func() {
		safe = rowIsFree[i] && diagonal1IsFree[j+i] && diagonal2IsFree[j-i+7]
	}

	setQueen := func() {
		rowIsFree[i] = false
		diagonal1IsFree[j+i] = false
		diagonal2IsFree[j-i+7] = false
	}
	removeQueen := func() {
		rowIsFree[i] = true
		diagonal1IsFree[j+i] = true
		diagonal2IsFree[j-i+7] = true
	}
	tryColumn := func() {
		for {
			advancePointer()
			testSquare()
			if safe || lastSquare() {
				break
			}
		}
	}

	regress := func() {
		reconsiderPriorColumn()
		if !regressOutOfFirstCol() {
			removeQueen()
			if lastSquare() {
				reconsiderPriorColumn()
				if !regressOutOfFirstCol() {
					removeQueen()
				}
			}
		}
	}

	considerFirstColumn()
	for {
		tryColumn()
		if safe {
			setQueen()
			considerNextColumn()
		} else {
			regress()
		}
		if lastColDone() || regressOutOfFirstCol() {
			break
		}
	}

	// trim redundant
	return x[1:9]
}
