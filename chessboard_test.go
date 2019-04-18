package main

import (
	"image/png"
	"os"
	"testing"
)

func TestDrawChessboard(t *testing.T) {
	pieces := []ChessPiece{}
	for _, pos := range []string{"d2", "a8", "g4"} {
		p, err := NewChessPiece(Queen, pos, "images/queen.png")
		if err != nil {
			t.Fatal(err)
		}
		pieces = append(pieces, *p)
	}

	board := NewChessBoard(60, pieces)
	img, err := board.Draw()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("chessboard.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
