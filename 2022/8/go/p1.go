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

	type Cell struct {
		maxR    int
		maxL    int
		maxT    int
		maxB    int
		visible bool
	}
	dp := make([][]Cell, height)
	for i := range dp {
		dp[i] = make([]Cell, width)
		for j := range dp[i] {
			if i == 0 || j == 0 || i == height-1 || j == width-1 {
				dp[i][j].maxR = 0
				dp[i][j].maxL = 0
				dp[i][j].maxT = 0
				dp[i][j].maxB = 0
				dp[i][j].visible = true
			} else {
				dp[i][j].maxR = math.MaxInt
				dp[i][j].maxL = math.MaxInt
				dp[i][j].maxT = math.MaxInt
				dp[i][j].maxB = math.MaxInt
				dp[i][j].visible = false
			}
		}
	}

	update := func(x, y int) {
		l := max(data[x-1][y], dp[x-1][y].maxL)
		r := max(data[x+1][y], dp[x+1][y].maxR)
		t := max(data[x][y-1], dp[x][y-1].maxT)
		b := max(data[x][y+1], dp[x][y+1].maxB)
		dp[x][y].maxR = r
		dp[x][y].maxL = l
		dp[x][y].maxT = t
		dp[x][y].maxB = b
		if min(l, r, b, t) < data[x][y] {
			dp[x][y].visible = true
		}
	}

	for m := 0; m <= min(width, height)/2; m += 1 {
		l, r := 1+m, width-2-m
		t, b := 1+m, height-2-m

		i, j := t, l

		for j <= r {
			update(i, j)
			if j == r {
				break
			}
			j += 1
		}

		for i <= b {
			update(i, j)
			if i == b {
				break
			}
			i += 1
		}

		for j >= l {
			update(i, j)
			if j == l {
				break
			}
			j -= 1
		}

		for i >= t {
			update(i, j)
			if i == t {
				break
			}
			i -= 1
		}
	}

	for m := min(width, height)/2; m >= 0; m -= 1 {
		l, r := 1+m, width-2-m
		t, b := 1+m, height-2-m

		i, j := t, l

		for i <= b {
			update(i, j)
			if i == b {
				break
			}
			i += 1
		}

		for j <= r {
			update(i, j)
			if j == r {
				break
			}
			j += 1
		}

		for i >= t {
			update(i, j)
			if i == t {
				break
			}
			i -= 1
		}

		for j >= l {
			update(i, j)
			if j == l {
				break
			}
			j -= 1
		}
	}

	res := 0

	for _, l := range dp {
		for _, c := range l {
			if c.visible {
				res += 1
			}
		}
	}

	fmt.Println(res)
}
