package main

import (
	"fmt"
	"image"
	"image/gif"
	"os"
	"testing"
	"time"
)

func TestMaze(t *testing.T) {
	w := 10
	h := 10

	stack := NewStack()
	start := time.Now()
	cw, ch, ww := 20, 20, 5

	maze, path := NewMaze(w, h, point{0, 0}, point{w - 1, h - 1},
		DFS(stack, time.Now().Unix()),
	)

	fmt.Printf("Gen Time %v\n", time.Since(start))
	maze.ResetVisitedCells()

	start = time.Now()
	fmt.Println(maze)
	fmt.Printf("DrawTxt Time %v\n", time.Since(start))

	start = time.Now()
	mazeImg := &gif.GIF{
		Image: []*image.Paletted{Draw(maze, white, black, cw, ch, ww)},
		Delay: []int{10},
	}
	fmt.Printf("Draw Maze Time %v\n", time.Since(start))

	start = time.Now()
	mazeAnimGen := AnimatePath(maze, nil, path, nil, white, black, cw, ch, ww, 10)
	fmt.Printf("Draw Maze Gen Time %v\n", time.Since(start))

	mazePath := make([]*cell, 0, w)
	visited := make([]*cell, 0, w)

	start = time.Now()
	isPathFound := FindPath(maze, maze.Begin(), maze.End(), &mazePath, &visited)
	fmt.Printf("PathFound? with Recursion %v cells visited %d, path len %d Time %+v\n",
		isPathFound, len(visited), len(mazePath), time.Since(start))
	maze.ResetVisitedCells()

	start = time.Now()
	mazeAnimPath := AnimatePath(maze, visited, mazePath, green, blue, black, cw, ch, ww, 10)
	fmt.Printf("Draw Maze Find Path Time %v\n", time.Since(start))

	mazePath = mazePath[:0]
	visited = visited[:0]

	start = time.Now()
	isPathFound = FindShortestPath(maze, maze.Begin(), maze.End(), &mazePath, &visited)
	fmt.Printf("ShortestPathFound? with BFS %v cells visited %d path len %d Time %+v\n",
		isPathFound, len(visited), len(mazePath), time.Since(start))
	maze.ResetVisitedCells()

	start = time.Now()
	mazeAnimShortPath := AnimatePath(maze, visited, mazePath, yellow, red, black, cw, ch, ww, 10)
	fmt.Printf("Draw Maze Find Shortest Path Time %v\n", time.Since(start))

	// use maze as bg
	mazeAnimPath.Image[0] = mazeImg.Image[0]
	mazeAnimShortPath.Image[0] = mazeImg.Image[0]

	data := []struct {
		name string
		img  *gif.GIF
	}{
		{"maze.gif", mazeImg},
		{"maze-anim-gen-with-dfs.gif", mazeAnimGen},
		{"maze-anim-path-with-dfs-recur.gif", mazeAnimPath},
		{"maze-anim-path-with-bfs-weight.gif", mazeAnimShortPath},
	}

	var f *os.File
	var err error
	for _, img := range data {
		if f, err = os.Create(img.name); err == nil {
			if err = gif.EncodeAll(f, img.img); err != nil {
				panic(err)
			}
		}
	}

}
