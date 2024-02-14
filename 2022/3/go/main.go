package main

import (
	"fmt"
	"os"
	"strings"
)

func getPriority(ch byte) int {
	if ch >= 'a' && ch <= 'z' {
		return int(ch-'a') + 1
	}

	if ch >= 'A' && ch <= 'Z' {
		return int(ch-'A') + 27
	}

	panic("Uhreachable")
}

func main() {
	content, err := os.ReadFile("input.txt")

	if err != nil {
		panic(err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	res := 0

	for i := 0; i < len(lines) - 1; i += 3 {
		sacks := [][]byte{
			[]byte(lines[i]),
			[]byte(lines[i + 1]),
			[]byte(lines[i + 2]),
		};

		ch := func() byte {
			for _, stCh := range sacks[0] {
				found := false

				for _, ndCh := range sacks[1] {
					if stCh == ndCh {
						found = true
					}
				}

				if !found {
					continue
				}


				for _, ndCh := range sacks[2] {
					if stCh == ndCh {
						return stCh
					}
				}
			}

			panic("Unreachable")
		}()

		res += getPriority(ch)
	}

	fmt.Println(res)
}
