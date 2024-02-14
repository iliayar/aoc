package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func signum(n int) int {
	if n < 0 {
		return -1
	} else if n > 0 {
		return 1
	} else {
		return 0
	}
}

func calcTailPosition(head Point, tail Point) (tailRes Point) {
	dx := head.x - tail.x
	dy := head.y - tail.y

	if abs(dy) > 1 || abs(dx) > 1 {
		tailRes = Point{tail.x + signum(dx), tail.y + signum(dy)}
	} else {
		tailRes = tail
	}

	// fmt.Println(head, tail, " -> ", tailRes)

	return
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := strings.Trim(string(content), "\n\r\t ")

	knots := make([]Point, 10)

	tailPoints := map[Point]bool{}
	tailPoints[knots[9]] = true

	for _, line := range strings.Split(contentStr, "\n") {
		args := strings.Split(line, " ")

		dx, dy := func() (int, int) {
			switch args[0] {
			case "R":
				return 1, 0
			case "L":
				return -1, 0
			case "D":
				return 0, 1
			case "U":
				return 0, -1
			}
			panic("Unknown direction")
		}()

		count, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
        knots[k] = calcTailPositionknots[k - 1], knots[k]

		for i := 0; i < count; i += 1 {
			knots[0].x += dx
			knots[0].y += dy

			for k := 1; k < len(knots); k += 1 {
				knots[k] = calcTailPosition(knots[k-1], knots[k])
			}

			tailPoints[knots[len(knots)-1]] = true
		}
	}

	res := len(tailPoints)
	fmt.Println(res)
}
