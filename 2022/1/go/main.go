package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	cur := 0

	cs := []int{}

	for _, line := range lines {
		if line == "" {
			cs = append(cs, cur)
			cur = 0
			continue
		}

		n, err := strconv.Atoi(line)

		if err != nil {
			panic(err)
		}

		cur += n
	}

	res := 0

	sort.Sort(sort.IntSlice(cs))
	for i := 0; i < 3; i++ {
		res += cs[len(cs) - i - 1]
	}
	fmt.Println(res)
}
