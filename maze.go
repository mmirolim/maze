package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"math/rand"
	"strings"
)

var (
	red    = color.RGBA{R: 255, A: 255}
	blue   = color.RGBA{B: 255, A: 255}
	green  = color.RGBA{G: 255, A: 255}
	black  = color.RGBA{0, 0, 0, 0}
	white  = color.RGBA{255, 255, 255, 255}
	yellow = color.RGBA{255, 255, 102, 255}
)

func DFS(stack *stack, seed int64) func(*Maze) (*Maze, []*cell) {
	return func(m *Maze) (*Maze, []*cell) {
		rand.Seed(seed)
		genPath := make([]*cell, 0, m.w*m.h)
		current := m.Begin()
		current.visited = true
		filter := func(c *cell) bool {
			return !c.visited
		}
		for {
			unvisited := m.AdjacentCells(current, filter)
			if len(unvisited) > 0 {
				genPath = append(genPath, current)
				next := unvisited[rand.Intn(len(unvisited))]
				stack.Push(current)
				m.RmWall(current, next)
				current = next
				current.visited = true
			} else if stack.Len() > 0 {
				current = stack.Pop()
			} else {
				break
			}
		}
		return m, genPath
	}
}

// returns path and all visited cells
// DFS recur search
func FindPath(m *Maze, start, end *cell, path, visited *[]*cell) bool {
	if start.visited {
		return false
	}
	start.visited = true
	*path = append(*path, start)
	*visited = append(*visited, start)
	if start == end {
		return true
	}
	//filter connected cells
	unvisited := m.AdjacentCells(start, func(c *cell) bool {
		return isConnected(start, c) && !c.visited
	})

	for _, next := range unvisited {
		if FindPath(m, next, end, path, visited) {
			return true
		}
	}

	// cell is not part of path
	*path = (*path)[:len(*path)-1]

	return false
}

// returns path and all visited cells
// BFS search with storing dist for each cell from start point
func FindShortestPath(m *Maze, start, end *cell, path, visited *[]*cell) bool {
	q := make([]*cell, 0, 4)
	start.visited = true
	*visited = append(*visited, start)
	current := start
	// store distance from start point
	visitedCellDist := make(map[*cell]int)
	stepsFromStart := 0
	for {

		//filter connected, not visited cells
		unvisited := m.AdjacentCells(current, func(c *cell) bool {
			return isConnected(current, c) && !c.visited
		})
		stepsFromStart++
		// check all connected cells
		for _, c := range unvisited {
			c.visited = true
			q = append(q, c)
			*visited = append(*visited, c)
			// store dist
			visitedCellDist[c] = stepsFromStart

			if c != end {
				continue
			}

			// found end point
			// compute shortest path backwards
			min := 1<<31 - 1
			for {
				// add cell to shortest path
				*path = append(*path, c)
				if c == start {
					break
				}
				// get all connected and visited cells
				connected := m.AdjacentCells(c, func(c2 *cell) bool {
					return isConnected(c, c2) && c2.visited
				})
				// find cell which is closest to start point
				for _, next := range connected {
					if dist := visitedCellDist[next]; min > dist {
						c = next
						min = dist
					}
				}
			}

			return true

		}

		if len(q) > 0 {
			current = q[0]
			q[0] = nil
			q = q[1:]
		} else {
			break
		}

	}

	return false
}

type point struct {
	x, y int
}

type cell struct {
	left, up, right, down bool // doors if exits
	visited               bool
	point
}

type Maze struct {
	w, h        int
	entry, exit point
	cells       [][]*cell
}

// remove adjacent cells walls
func (m *Maze) RmWall(cell1, cell2 *cell) {
	dx := cell1.x - cell2.x
	dy := cell1.y - cell2.y
	// panic if cell not adjacent or same
	if dx+dy == 0 || dx+dy > 1 || dx+dy < -1 {
		panic("cells not adjacent or same")
	}

	switch {
	case dx > 0:
		// left wall
		cell1.left = true
		cell2.right = true
	case dx < 0:
		// right wall
		cell1.right = true
		cell2.left = true
	case dy > 0:
		// up wall
		cell1.up = true
		cell2.down = true
	case dy < 0:
		// down wall
		cell1.down = true
		cell2.up = true
	default:
		panic("bug")
	}

}

// return adjecent filtered cells
func (m *Maze) AdjacentCells(c *cell, filter func(*cell) bool) []*cell {
	p := c.point
	cells := []*cell{}
	data := []struct {
		cond bool
		point
	}{
		{p.x > 0, point{p.x - 1, p.y}},
		{p.x < m.w-1, point{p.x + 1, p.y}},
		{p.y > 0, point{p.x, p.y - 1}},
		{p.y < m.h-1, point{p.x, p.y + 1}},
	}
	for _, d := range data {
		if d.cond && filter(m.cells[d.x][d.y]) {
			cells = append(cells, m.cells[d.x][d.y])
		}
	}

	return cells
}

func isConnected(c1, c2 *cell) bool {
	if c1.y+1 == c2.y && c2.up && c1.down ||
		c1.x-1 == c2.x && c2.right && c1.left ||
		c1.y-1 == c2.y && c2.down && c1.up ||
		c1.x+1 == c2.x && c2.left && c1.right {
		return true
	}
	return false
}

func NewMaze(w, h int, entry, exit point,
	generator func(*Maze) (*Maze, []*cell)) (*Maze, []*cell) {

	if w < 1 || h < 1 {
		panic("w, h should be > 1")
	}
	for _, p := range []point{entry, exit} {
		if p.x > w-1 || p.x < 0 || p.y > h-1 || p.y < 0 {
			panic("start and end point should be inside maze")
		}
	}

	m := &Maze{w: w, h: h, entry: entry, exit: exit, cells: make([][]*cell, w)}

	for x := range m.cells {
		if m.cells[x] == nil {
			m.cells[x] = make([]*cell, h)
		}
		for y := range m.cells[x] {
			m.cells[x][y] = &cell{point: point{x, y}}
		}
	}

	if generator == nil {
		return m, nil
	}
	return generator(m)
}

func (m *Maze) ResetVisitedCells() {
	for x := range m.cells {
		for y := range m.cells[x] {
			m.cells[x][y].visited = false
		}
	}
}

func (m *Maze) Begin() *cell {
	return m.cells[m.entry.x][m.entry.y]
}

func (m *Maze) End() *cell {
	return m.cells[m.exit.x][m.exit.y]
}

func (m *Maze) String() string {
	output, hline, vline := []byte{}, []byte{}, []byte{}

	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			mark := " "
			if x == m.entry.x && y == m.entry.y {
				mark = "S"
			}
			if x == m.exit.x && y == m.exit.y {
				mark = "E"

			}
			hElm := "+---"
			vElm := "| " + mark + " "

			if m.cells[x][y].up {
				hElm = "+   "
			}

			if m.cells[x][y].left {
				vElm = "  " + mark + " "
			}

			hline = append(hline, []byte(hElm)...)
			vline = append(vline, []byte(vElm)...)
		}

		output = append(output, append(hline, []byte("+\n")...)...)
		output = append(output, append(vline, []byte("|\n")...)...)
		hline, vline = hline[:0], vline[:0]
	}

	return strings.Join([]string{string(output), strings.Repeat("+---", m.w), "+\n"}, "")
}

func Draw(m *Maze, fill, border color.Color, cw, ch, ww int) *image.Paletted {
	r := image.Rect(0, 0, m.w*cw, m.h*cw)
	img := image.NewPaletted(r, palette.WebSafe)

	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			cell := m.cells[x][y]
			rect := image.Rect(cell.x*cw, cell.y*ch, cell.x*cw+cw, cell.y*ch+ch)
			DrawCell(cell, img.SubImage(rect).(*image.Paletted), fill, border, cw, ch, ww)
		}
	}
	return img
}

func DrawCell(cell *cell, img *image.Paletted,
	fill, border color.Color,
	cw, ch, ww int) {

	rect := img.Bounds()
	x0, y0, x1, y1 := rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, fill)
			// check walls
			if (!cell.left && x < x0+ww) ||
				(!cell.right && x > x1-ww) ||
				(!cell.up && y < y0+ww) ||
				(!cell.down && y > y1-ww) {
				img.Set(x, y, border)
			}

		}
	}
}

// speed in 100th of second
func AnimatePath(m *Maze, visited, path []*cell,
	fillVis, fillPath, border color.Color,
	cw, ch, ww, speed int) *gif.GIF {

	r := image.Rect(0, 0, m.w*cw, m.h*ch)
	img := image.NewPaletted(r, palette.WebSafe)
	imgs := []image.Image{img}

	lenVisited := len(visited)
	// join visited and path cells
	visited = append(visited, path...)
	for i, cell := range visited {
		fill := fillVis
		rect := image.Rect(cell.x*cw, cell.y*ch, cell.x*cw+cw, cell.y*ch+ch)
		cellImg := image.NewPaletted(rect, palette.WebSafe)
		if i >= lenVisited {
			fill = fillPath // fill path diff
		}
		DrawCell(cell, cellImg, fill, border, cw, ch, ww)
		imgs = append(imgs, cellImg)

	}

	gifAnim := &gif.GIF{
		Image:     make([]*image.Paletted, len(imgs)),
		Delay:     make([]int, len(imgs)),
		LoopCount: -1,
	}
	for i := range imgs {
		gifAnim.Image[i] = imgs[i].(*image.Paletted)
		gifAnim.Delay[i] = speed
	}

	return gifAnim
}

type stack struct {
	cells []*cell
}

func NewStack() *stack {
	return &stack{make([]*cell, 0, 10)}
}

func (s *stack) Push(c *cell) {
	s.cells = append(s.cells, c)
}

func (s *stack) Pop() *cell {
	c := s.cells[len(s.cells)-1]
	s.cells = s.cells[:len(s.cells)-1]
	return c
}

func (s *stack) Len() int {
	return len(s.cells)
}
