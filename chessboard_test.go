package main

import (
	"image/png"
	"os"
	"testing"
)

func TestDrawChessboard(t *testing.T) {
	pieces := []ChessPiece{}
	for _, pos := range []string{"d2", "a8", "g4"} {
		p, err := NewChessPiece(Queen, pos)
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

func TestTheEightQueensProblemMySolution(t *testing.T) {
	pieces := TheEightQueensProblemMySolution()
	if len(pieces) != 8 {
		t.Errorf("solution not found, expected 8 got %v", len(pieces))
	}
	board := NewChessBoard(60, pieces)
	img, err := board.Draw()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("8-queen-solution.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTheEightQueensProblemNW(t *testing.T) {
	pieces := TheEightQueensProblemNW()
	if len(pieces) != 8 {
		t.Errorf("solution not found, expected 8 got %v", len(pieces))
	}

	// convert pos to chessboard pos
	toChessPos := func(i, j int) string {
		return string([]byte{byte(i) + 'a' - 1, byte(j) + '0'})
	}

	chessPieces := make([]ChessPiece, 0, 8)

	for i := range pieces {
		p, err := NewChessPiece(Queen, toChessPos(i+1, pieces[i]))
		if err != nil {
			panic(err)
		}

		chessPieces = append(chessPieces, *p)
	}

	board := NewChessBoard(60, chessPieces)
	img, err := board.Draw()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("8-queen-nw-solution.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkTheEightQueensProblemNW(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TheEightQueensProblemNW()
	}
}

func BenchmarkTheEightQueensProblemMySolution(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TheEightQueensProblemMySolution()
	}
}
