package main

import (
	"image/png"
	"math/rand"
	"os"
	"testing"
)

var (
	initPos = [9][9]int{
		[9]int{5, 6, 0, 8, 4, 7, 0, 0, 0},
		[9]int{3, 0, 9, 0, 0, 0, 6, 0, 0},
		[9]int{0, 0, 8, 0, 0, 0, 0, 0, 0},
		[9]int{0, 1, 0, 0, 8, 0, 0, 4, 0},
		[9]int{7, 9, 0, 6, 0, 2, 0, 1, 8},
		[9]int{0, 5, 0, 0, 3, 0, 0, 9, 0},
		[9]int{0, 0, 0, 0, 0, 0, 2, 0, 0},
		[9]int{0, 0, 6, 0, 0, 0, 8, 0, 7},
		[9]int{0, 0, 0, 3, 1, 6, 0, 5, 9},
	}

	expectedResult = [9][9]int{
		[9]int{5, 6, 1, 8, 4, 7, 9, 2, 3},
		[9]int{3, 7, 9, 5, 2, 1, 6, 8, 4},
		[9]int{4, 2, 8, 9, 6, 3, 1, 7, 5},
		[9]int{6, 1, 3, 7, 8, 9, 5, 4, 2},
		[9]int{7, 9, 4, 6, 5, 2, 3, 1, 8},
		[9]int{8, 5, 2, 1, 3, 4, 7, 9, 6},
		[9]int{9, 3, 5, 4, 7, 8, 2, 6, 1},
		[9]int{1, 4, 6, 2, 9, 5, 8, 3, 7},
		[9]int{2, 8, 7, 3, 1, 6, 4, 5, 9},
	}
)

func TestSudokuSolveWithBacktrackingWithoutInitDigits(t *testing.T) {
	sudoku := NewSudoku([9][9]int{})
	err := sudoku.SolveWithBacktracking()
	if err != nil {
		t.Fatalf("error %v", err)
	}
	img, err := sudoku.Draw(60)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("sudoku-with-no-init-pos.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}

}

func TestSudokuSolveWithBacktrackingWithInitDigits(t *testing.T) {
	sudoku := NewSudoku(initPos)
	err := sudoku.SolveWithBacktracking()
	if err != nil {
		t.Fatalf("error %v", err)
	}
	r := sudoku.GetSolution()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if r[i][j] != expectedResult[i][j] {
				t.Fatalf("result != expected result at (%v, %v) %v != %v", i, j, r[i][j], expectedResult[i][j])
			}
		}
	}
	img, err := sudoku.Draw(60)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("sudoku-with-init-pos.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSudokuSolveWithSimulatedAnnealing(t *testing.T) {
	sudoku := NewSudoku(initPos)
	err := sudoku.SolveWithSA()
	if err != nil {
		t.Fatalf("error %v", err)
	}
	r := sudoku.GetSolution()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if r[i][j] != expectedResult[i][j] {
				t.Fatalf("result != expected result at (%v, %v) %v != %v", i, j, r[i][j], expectedResult[i][j])
			}
		}
	}
	img, err := sudoku.Draw(60)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("sudoku-with-sa.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}

}

func BenchmarkSudokuSolveWithBacktracking(b *testing.B) {
	sudoku := NewSudoku(initPos)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = sudoku.SolveWithBacktracking()
	}
}

func BenchmarkSudokuSolveWithSA(b *testing.B) {
	sudoku := NewSudoku(initPos)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = sudoku.SolveWithSA()
	}
}

func BenchmarkRandIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Intn(9)
	}
}

func BenchmarkRandShift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int((rand.Int() * 9) >> 32)
	}
}
