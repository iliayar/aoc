package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Noop struct {
}

type Addx struct {
	op int
}

type State struct {
	cycle int
	x     int
}

type Instruction interface {
	execute(state *State)
	cyclesCount() int
}

func (Noop) execute(state *State) {}

func (Noop) cyclesCount() int { return 1 }

func (ins Addx) execute(state *State) {
	state.x += ins.op
}

func (Addx) cyclesCount() int { return 2 }

func parseLine(line string) Instruction {
	args := strings.Split(line, " ")

	switch args[0] {
	case "noop":
		return Noop{}
	case "addx":
		op, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		return Addx{op}
	default:
		panic(fmt.Sprintf("Unknown instruction: %s", line))
	}
}

type CRT struct {
    lines [][]byte
    currentLine []byte
}

func (crt *CRT) draw(state *State) {
    pos := len(crt.currentLine)

    if pos == 40 {
        crt.lines = append(crt.lines, crt.currentLine)
        crt.currentLine = []byte{}
        pos = 0
    }

    if pos >= state.x - 1 && pos <= state.x + 1 {
        crt.currentLine = append(crt.currentLine, '#')
    } else {
        crt.currentLine = append(crt.currentLine, '.')
    }
}

func (crt *CRT) print() {
    crt.lines = append(crt.lines, crt.currentLine)
    for _, line := range crt.lines {
        fmt.Println(string(line))
    }
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := strings.Trim(string(content), "\n\r\t ")

	inss := []Instruction{}

	for _, line := range strings.Split(contentStr, "\n") {
		inss = append(inss, parseLine(line))
	}

	state := State{}

    crt := CRT{}

	state.cycle = 1
    state.x = 1
	for _, ins := range inss {
		for i := 0; i < ins.cyclesCount(); i += 1 {
			state.cycle += 1

            crt.draw(&state) 
		}

        ins.execute(&state)
	}

    crt.print()
}
