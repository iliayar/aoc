package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func min(ns ...int) int {
	res := math.MaxInt

	for _, n := range ns {
		if n < res {
			res = n
		}
	}

	return res
}

func max(ns ...int) int {
	res := math.MinInt

	for _, n := range ns {
		if n > res {
			res = n
		}
	}

	return res
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := string(content)
	lines := strings.Split(strings.Trim(contentStr, "\n\r "), "\n")

	height, width := len(lines), len(lines[0])

	data := make([][]int, height)
	for i := range lines {
		data[i] = make([]int, width)
		for j := range lines[i] {
			data[i][j] = int(lines[i][j] - '0')
		}
	}

	findDir := func(x, y int, dx, dy int) int {
		i := 0
		orig := data[x][y]

		for {
			x += dx
			y += dy
			i += 1

			if !(x >= 0 && x < width && y >= 0 && y < height) {
                i -= 1
                break
			}

			if orig <= data[x][y] {
				break
			}

		}

		return i
	}

	res := 0

	for i, l := range data {
		for j := range l {
			l := findDir(i, j, -1, 0)
			r := findDir(i, j, 1, 0)
			t := findDir(i, j, 0, -1)
			b := findDir(i, j, 0, 1)

			res = max(res, l*r*t*b)
		}
	}

	// fmt.Println(dp)
	fmt.Println(res)
}
