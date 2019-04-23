package main

import (
	"image/png"
	"os"
	"testing"
)

func TestDrawSudokuBoard(t *testing.T) {
	vals := [3][3][3][3]int{
		[3][3][3]int{
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
		},
		[3][3][3]int{
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
		},
		[3][3][3]int{
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
			[3][3]int{[3]int{1, 2, 3}, [3]int{4, 5, 6}, [3]int{7, 8, 9}},
		},
	}
	img, err := DrawSudokuBoard(60, vals)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("sudoku.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}

}
