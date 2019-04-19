package main

import (
	"errors"
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
func TheEightQueensProblemMySolution() []ChessPiece {
	// convert pos to chessboard pos
	toChessPos := func(i, j int) string {
		return string([]byte{byte(i) + 'a' - 1, byte(j) + '0'})
	}
	// i is letters axis, j is numbers axis
	check := func(i1, j1, i2, j2 int) bool {
		// check horizontal and vertical coords
		if i1 == i2 || j1 == j2 {
			return false
		}

		// check diagonal
		c1pos := j1 - i1
		c1neg := j1 + i1
		c2pos := j2 - i2
		c2neg := j2 + i2
		if c1pos == c2pos || c1neg == c2neg {
			return false
		}

		return true
	}
	type piece struct {
		piece ChessPiece
		pos   [2]int
	}

	pieces := []piece{}
	var search func(i, j int) error

	search = func(i, j int) error {
		// found all pieces
		if len(pieces) > 7 {
			return nil
		}
		// change row
		if j > 8 {
			j = 1
			i++
		}
		// end of chessboard
		if i > 8 {
			return errors.New("no valid position")
		}
		invalidPos := false
		for _, p := range pieces {
			if !check(p.pos[0], p.pos[1], i, j) {
				invalidPos = true
				break
			}
		}
		if !invalidPos {
			// valid position
			p, err := NewChessPiece(Queen, toChessPos(i, j))
			if err != nil {
				panic(err)
			}
			pieces = append(pieces, piece{*p, [2]int{i, j}})
		}

		j++
		err := search(i, j)
		if err != nil {
			// on err take last element
			// and continue from that position with shift
			p := pieces[len(pieces)-1]
			pieces = pieces[:len(pieces)-1]
			return search(p.pos[0], p.pos[1]+1)
		}

		return nil
	}

	chessPieces := make([]ChessPiece, 0, 8)
	// start search
	_ = search(1, 1)

	for i := range pieces {
		p, err := NewChessPiece(Queen, toChessPos(pieces[i][0], pieces[i][1]))
		if err != nil {
			panic(err)
		}

		chessPieces = append(chessPieces, *p)
	}

	return chessPieces
}

// TheEightQueensProblemNW solution to 8-queen puzzle
// described by Niklaus Wirth in
// http://plbpc001.ouhk.edu.hk/~mt311/optional-reading/stepwise.pdf
func TheEightQueensProblemNW() {
	var board, pointer, safe

	tryColumn := func() {
		for advancePointer {
			testSquare
			if safe || lastSquare {
				break
			}
		}
	}

	regress := func() {
		reconsiderPriorColumn
		if !regressOutOfFirstCol {
			removeQueen
			if lastSquare {
				reconsiderPriorColumn
				if !regressOutOfFirstCol {
					removeQueen
				}
			}
		}
	}
	
	considerFirstColumn
	for tryColumn {
		if safe {
			setQueen
			considerNextColumn
		} else {
			regress
		}
		if lastColDone || regressOutOfFirstCol {
			break
		}
	}
}
