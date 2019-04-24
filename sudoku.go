package main

import (
	"image"
	"image/color"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

// DrawSudokuBoard ...
func DrawSudokuBoard(cellSize int, num [3][3][3][3]int) (*image.RGBA, error) {
	black := color.RGBA{0, 0, 0, 255}
	//white := color.RGBA{255, 255, 255, 255}
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
	fontOptionSize := 30
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
					d.DrawString(strconv.Itoa(num[i][j][x][y]))

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
