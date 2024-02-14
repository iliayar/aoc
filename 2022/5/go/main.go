package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type State int

const (
	INIT State = iota
	ACTIONS
)

type Crate byte

func parseCrateLine(line string) []*Crate {
	lineBytes := []byte(line)

	cratesCount := (len(line) + 1) / 4
	res := make([]*Crate, cratesCount)

	for i := 0; i < len(line); i += 4 {
		if lineBytes[i] == '[' && lineBytes[i+2] == ']' {
			crate := Crate(lineBytes[i+1])
			res[i/4] = &crate
		}
	}

	return res
}

type Action = struct {
	from  int
	count int
	to    int
}

func parseActionLine(line string) Action {
	parts := strings.Split(line, " ")

	count, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	from, err := strconv.Atoi(parts[3])
	if err != nil {
		panic(err)
	}

	to, err := strconv.Atoi(parts[5])
	if err != nil {
		panic(err)
	}

	return Action{
		from:  from - 1,
		count: count,
		to:    to - 1,
	}
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := string(content)

	state := INIT

	var cratesRaw [][]*Crate
	var actions []Action

	for _, line := range strings.Split(contentStr, "\n") {
		if state == INIT {
			if line == "" {
				state = ACTIONS
			} else {
				cratesRaw = append(cratesRaw, parseCrateLine(line))
			}
		} else if state == ACTIONS {
			if line == "" {
				break
			} else {
				actions = append(actions, parseActionLine(line))
			}
		}
	}

	stacksCount := len(cratesRaw[len(cratesRaw)-1])
	var crates [][]Crate = make([][]Crate, stacksCount)

	for i := len(cratesRaw) - 1; i >= 0; i -= 1 {
		for j, crate := range cratesRaw[i] {
			if crate != nil {
				crates[j] = append(crates[j], *crate)
			}
		}
	}

	for _, action := range actions {
		fromCrate := crates[action.from]
		toCrate := crates[action.to]

		// for i := 0; i < action.count; i += 1 {
		//     crate := fromCrate[len(fromCrate) - 1]
		//     fromCrate = fromCrate[:len(fromCrate) - 1]
		//
		//     toCrate = append(toCrate, crate)
		// }

		cratesChunk := fromCrate[len(fromCrate)-action.count:]
		fromCrate = fromCrate[:len(fromCrate)-action.count]

		toCrate = append(toCrate, cratesChunk...)

		crates[action.from] = fromCrate
		crates[action.to] = toCrate
	}

	for _, stack := range crates {
		fmt.Printf("%c", stack[len(stack)-1])
	}
	fmt.Printf("\n")
}
