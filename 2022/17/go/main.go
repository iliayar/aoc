package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

type Point struct {
	x, y int
}

func (p Point) add(other Point) Point {
	return Point{p.x + other.x, p.y + other.y}
}

var PDown Point = Point{0, -1}
var PLeft Point = Point{-1, 0}
var PRight Point = Point{1, 0}

type Mask []Point

func (m Mask) intersects(other Mask) bool {
	for _, p := range m {
		for _, op := range other {
			if p == op {
				return true
			}
		}
	}
	return false
}

func (m *Mask) extend(other Mask) {
	for _, p := range other {
		*m = append(*m, p)
	}
}

func (m Mask) contains(op Point) bool {
	return m.intersects(Mask{op})
}

func (m Mask) print() {
	bounds, err := m.bounds()
	if err != nil {
		return
	}

	for y := bounds.t; y >= 0; y -= 1 {
		for x := 0; x <= bounds.r; x += 1 {
			if m.contains(Point{x, y}) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

type Bounds struct {
	t, b, l, r int
}

func (m Mask) bounds() (Bounds, error) {
	if len(m) == 0 {
		return Bounds{}, errors.New("No points")
	}

	t := math.MinInt
	b := math.MaxInt
	l := math.MaxInt
	r := math.MinInt

	for _, p := range m {
		if p.x < l {
			l = p.x
		}
		if p.x > r {
			r = p.x
		}
		if p.y < b {
			b = p.y
		}
		if p.y > t {
			t = p.y
		}
	}

	return Bounds{
		t: t,
		b: b,
		l: l,
		r: r,
	}, nil
}

func (m Mask) height() int {
	bounds, err := m.bounds()
	if err != nil {
		return 0
	}

	return bounds.t - bounds.b + 1
}

type RockPattern Mask

type Rock struct {
	pattern RockPattern
	pos     Point
}

func (r *Rock) getMask() Mask {
	res := []Point{}
	for _, p := range r.pattern {
		res = append(res, p.add(r.pos))
	}
	return Mask(res)
}

func (r Rock) move(pd Point) Rock {
	return Rock{
		pos:     r.pos.add(pd),
		pattern: r.pattern,
	}
}

type Jet int

const (
	RIGHT Jet = iota
	LEFT
)

func jetFromRune(r rune) Jet {
	switch r {
	case '<':
		return LEFT
	case '>':
		return RIGHT
	}
	panic("unreachable")
}

const WIDTH = 7
const OFFSET_LEFT = 2
const OFFSET_BOT = 4

type State struct {
	surface Mask
	jets    []Jet
	curJet  int
}

func makeState(jets []Jet) State {
	return State{
		jets:    jets,
		surface: Mask{},
		curJet:  0,
	}
}

func (s *State) dropRock(pattern RockPattern) {
	bounds, err := s.surface.bounds()
	start_y := OFFSET_BOT - 1
	if err == nil {
		start_y = bounds.t + OFFSET_BOT
	}
	r := Rock{
		pattern: pattern,
		pos:     Point{OFFSET_LEFT, start_y},
	}

	// fmt.Println(s.surface, r.getMask())

	for {
		curJet := s.jets[s.curJet%len(s.jets)]
		s.curJet += 1

		r1 := r
		if curJet == LEFT {
			r1 = r.move(PLeft)
		} else {
			r1 = r.move(PRight)
		}

		bounds, _ := r1.getMask().bounds()
		if r1.getMask().intersects(s.surface) || bounds.l < 0 || bounds.r >= WIDTH {
			r1 = r
		}

		r2 := r1.move(PDown)
		bounds, _ = r2.getMask().bounds()
		if r2.getMask().intersects(s.surface) || bounds.b < 0 {
			s.surface.extend(r1.getMask())
			break
		}

		r = r2

		// fmt.Println(s.surface, r.getMask())

		// r.getMask().print()
	}
}

func findDupParts(m Mask, n int) bool {
	bounds, err := m.bounds()
    if err != nil {
        return false
    }

    overallH := bounds.t - bounds.b + 1
	if overallH % n != 0 {
		return false
	}

    h := overallH / n
    cnts := make([]int, n)

	for _, p := range m {
        idx := p.y / h
        cnts[idx] += 1
	}

    for i := range cnts {
        if cnts[i] != cnts[0] {
            return false
        }
    }

    return true
}

func findDupMid(m Mask) bool {
	bounds, _ := m.bounds()
	// if (bounds.l + bounds.r) == 0 {
	// 	return false
	// }

    // mid = (0 + 11) / 2 = 5
    // 0..5 (6) | 6..11 (6)
	mid := (bounds.t + bounds.b) / 2
	bot := map[Point]bool{}
	top := map[Point]bool{}
	for _, p := range m {
		if p.y <= mid {
			bot[p] = true
		} else {
			top[p] = true
		}
	}

	for _, p := range m {
		if p.y > mid {
			_, e := bot[Point{p.x, p.y - mid - 1}]
			if !e {
				return false
			}
		} else {
			_, e := top[Point{p.x, p.y + mid + 1}]
			if !e {
                return false
            }
        }
	}

	return true
}


func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	patterns := []Mask{
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},
		{{2, 2}, {2, 1}, {2, 0}, {1, 0}, {0, 0}},
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	}

	contentStr := strings.Trim(string(content), "\n\r\t ")

	jets := []Jet{}
	for _, c := range contentStr {
		jets = append(jets, jetFromRune(c))
	}

	s := makeState(jets)

	for _, p := range patterns {
		p.print()
	}

    const PARTS = 3

	for i := 0; i < 100000; i += 1 {
		pattern := patterns[i%len(patterns)]
		s.dropRock(RockPattern(pattern))
        if findDupParts(s.surface, PARTS) {
        // if findDupMid(s.surface) {
            fmt.Println("found dup")
            // s.surface.print()
        }
		// s.surface.print()
		// fmt.Println("-------")
	}

	bounds, _ := s.surface.bounds()
	fmt.Println(bounds.t + 1)
}
