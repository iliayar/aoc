package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type CellType int

const (
	EMPTY CellType = iota
	ROCK
	SANDSOURCE
	SAND
)

func (ct CellType) show() string {
	switch ct {
	case EMPTY:
		return "."
	case ROCK:
		return "#"
	case SANDSOURCE:
		return "+"
	case SAND:
		return "O"
	}
	panic("unreachable")
}

type Cell struct {
	x int
	y int
}

type Bounds struct {
	right, left, upper, lower int
}

type Surface struct {
	m      map[Cell]CellType
	bounds Bounds
}

func createSurface() Surface {
	b := Bounds{}
	b.right, b.left = math.MinInt, math.MaxInt
	b.upper, b.lower = math.MinInt, math.MaxInt

	return Surface{
		m:      map[Cell]CellType{},
		bounds: b,
	}
}

func (s *Surface) drawLine(from, to Cell, cellType CellType) {
	if from.x == to.x {
		x := from.x
		for y := min(from.y, to.y); y <= max(from.y, to.y); y += 1 {
			s.drawPoint(Cell{x, y}, cellType)
		}
	} else if from.y == to.y {
		y := from.y
		for x := min(from.x, to.x); x <= max(from.x, to.x); x += 1 {
			s.drawPoint(Cell{x, y}, cellType)
		}
	} else {
		panic("Can draw only straight lines")
	}
}

func (s *Surface) drawPoint(cell Cell, cellType CellType) {
	if cell.x < s.bounds.left {
		s.bounds.left = cell.x
	}
	if cell.x > s.bounds.right {
		s.bounds.right = cell.x
	}
	if cellType != SAND {
		if cell.y < s.bounds.lower {
			s.bounds.lower = cell.y
		}
		if cell.y > s.bounds.upper {
			s.bounds.upper = cell.y
		}
	}
	s.m[cell] = cellType
}

func (s *Surface) get(cell Cell) CellType {
	if cell.y == s.bounds.upper+2 {
		return ROCK
	}
	return s.m[cell]
}

func (s *Surface) getBounds() Bounds {
	res := s.bounds
	res.upper += 2
	return res
}

func (s *Surface) show() string {
	b := s.getBounds()

	res := ""
	for y := b.lower; y <= b.upper; y += 1 {
		for x := b.left; x <= b.right; x += 1 {
			res += s.get(Cell{x, y}).show()
		}
		res += "\n"
	}
	return res
}

type Simulation struct {
	surface Surface

	sandSource Cell
}

func createSimulation() Simulation {
	return Simulation{
		surface: createSurface(),
	}
}

func (s *Simulation) addSandSource(cell Cell) {
	s.surface.drawPoint(cell, SANDSOURCE)
	s.sandSource = cell
}

func (s *Simulation) step() bool {
	x, y := s.sandSource.x, s.sandSource.y

    if s.surface.get(Cell{x, y}) == SAND {
        return false
    }

	for {
		downType := s.surface.get(Cell{x, y + 1})
		leftDownType := s.surface.get(Cell{x - 1, y + 1})
		rightDownType := s.surface.get(Cell{x + 1, y + 1})

		if downType == EMPTY {
			y += 1
		} else if leftDownType == EMPTY {
			y += 1
			x -= 1
		} else if rightDownType == EMPTY {
			y += 1
			x += 1
		} else {
			break
		}
	}

	s.surface.drawPoint(Cell{x, y}, SAND)
	return true
}

func parsePath(input string) []Cell {
	input = strings.Trim(input, "\n\r\t ")
	points := strings.Split(input, " -> ")
	res := []Cell{}
	for _, point := range points {
		parts := strings.Split(point, ",")
		x, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		res = append(res, Cell{x, y})
	}
	return res
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	sim := createSimulation()
	sim.addSandSource(Cell{500, 0})

	contentStr := strings.Trim(string(content), "\n\r\t ")
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		path := parsePath(line)
		for i := 0; i < len(path)-1; i += 1 {
			sim.surface.drawLine(path[i], path[i+1], ROCK)
		}
	}

	res := 0
	for sim.step() {
		// fmt.Println(sim.surface.show())
		res += 1
	}
	fmt.Println(res)
}

func min(xs ...int) int {
	res := math.MaxInt
	for _, x := range xs {
		if x < res {
			res = x
		}
	}
	return res
}

func max(xs ...int) int {
	res := math.MinInt
	for _, x := range xs {
		if x > res {
			res = x
		}
	}
	return res
}
