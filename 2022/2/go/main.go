package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	Rock = iota
	Paper
	Scissor
)

type Sign int

const (
	Win = iota
	Lose
	Draw
)

type Result int

func play(lhs, rhs Sign) Result {
	if rhs == getWin(lhs) {
		return Lose
	} else if rhs == lhs {
		return Draw
	} else {
		return Win
	}
}

func getWin(sign Sign) Sign {
	switch sign {
	case Rock:
		return Paper
	case Paper:
		return Scissor
	case Scissor:
		return Rock
	}

	panic("Unreachable")
}

func getLose(sign Sign) Sign {
	switch sign {
	case Rock:
		return Scissor
	case Paper:
		return Rock
	case Scissor:
		return Paper
	}

	panic("Unreachable")
}

func getForResult(sign Sign, result Result) Sign {
	switch result {
	case Win:
		return getWin(sign)
	case Lose:
		return getLose(sign)
	case Draw:
		return sign
	}

	panic("Unreachable")
}

func convertSign(sign string) (Sign, error) {
	if sign == "A" {
		return Rock, nil
	}

	if sign == "B" {
		return Paper, nil
	}

	if sign == "C" {
		return Scissor, nil
	}

	return 0, fmt.Errorf("Unknown sign %s", sign)
}

func convertResult(sign string) (Result, error) {
	if sign == "X" {
		return Lose, nil
	}

	if sign == "Y" {
		return Draw, nil
	}

	if sign == "Z" {
		return Win, nil
	}

	return 0, fmt.Errorf("Unknown sign %s", sign)
}

func scoreForSign(sign Sign) int {
	switch sign {
	case Rock:
		return 1
	case Paper:
		return 2
	case Scissor:
		return 3
	}

	panic("Unreachable")
}

func scoreForResult(result Result) int {
	switch result {
	case Win:
		return 6
	case Draw:
		return 3
	case Lose:
		return 0
	}

	panic("Unreachable")
}

func main() {
	content, err := os.ReadFile("input.txt")

	if err != nil {
		panic(err)
	}

	content_str := string(content)
	lines := strings.Split(content_str, "\n")

	score := 0

	for _, line := range lines {
		if line == "" {
			break
		}

		choices := strings.Split(line, " ")

		other, err := convertSign(choices[0])

		if err != nil {
			panic(err)
		}

		result, err := convertResult(choices[1])

		if err != nil {
			panic(err)
		}

		my := getForResult(other, result)

		score += scoreForSign(my)
		score += scoreForResult(result)
	}

	fmt.Println(score)
}
