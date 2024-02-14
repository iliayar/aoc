package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range = struct {
	from int
	to   int
}

func checkSubset(lhs, rhs Range) bool {
	if lhs.from >= rhs.from && lhs.to <= rhs.to {
		return true
	}

	return false
}

func checkDistinct(lhs, rhs Range) bool {
    return lhs.to < rhs.from || rhs.to < lhs.from
}

func stringToRange(s string) Range {
	fromTo := strings.Split(s, "-")

	from, err := strconv.Atoi(fromTo[0])

	if err != nil {
		panic(err)
	}

	to, err := strconv.Atoi(fromTo[1])

	if err != nil {
		panic(err)
	}

	return Range{
		from: from,
		to:   to,
	}
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := string(content)

    res := 0

	for _, line := range strings.Split(contentStr, "\n") {
        if line == "" {
            break
        }

		assignments := strings.Split(line, ",")

		lhs, rhs := stringToRange(assignments[0]), stringToRange(assignments[1])

        if !checkDistinct(lhs, rhs) {
            res += 1
        }
	}

    fmt.Println(res)
}
