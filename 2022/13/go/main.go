package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ValueInt int
type ValueArray []Value

type Ord int

const (
	LT Ord = iota
	GT
	EQ
)

type Value interface {
	compare(other Value) Ord
}

func (v ValueInt) compare(other Value) Ord {
	switch otherV := other.(type) {
	case ValueInt:
		lhs := int(v)
		rhs := int(otherV)
		if lhs < rhs {
			return LT
		} else if lhs > rhs {
			return GT
		} else {
			return EQ
		}
	case ValueArray:
		lhs := ValueArray([]Value{v})
		return lhs.compare(otherV)
	}

	panic("unreachable")
}

func (v ValueArray) compare(other Value) Ord {
	switch otherV := other.(type) {
	case ValueInt:
		rhs := ValueArray([]Value{otherV})
		return v.compare(rhs)
	case ValueArray:
		lhs := []Value(v)
		rhs := []Value(otherV)
		for i := 0; ; i += 1 {
			if i == len(lhs) && i == len(rhs) {
				return EQ
			} else if i == len(lhs) {
				return LT
			} else if i == len(rhs) {
				return GT
			}

			switch lhs[i].compare(rhs[i]) {
			case GT:
				return GT
			case LT:
				return LT
			default:
				continue
			}
		}
	}

	panic("uncreachable")
}

func parseValue(input string, i int) (Value, int) {
	if input[i] == '[' {
		i += 1

		res := []Value{}
		for i < len(input) {
			if input[i] == ']' {
				return ValueArray(res), i + 1
			} else {
				value, newI := parseValue(input, i)
				i = newI
				res = append(res, value)

				if input[i] == ',' {
					i += 1
					continue
				} else if input[i] == ']' {
					continue
				}

				panic("Unexpected token")
			}
		}

		panic("Unclosed bracket")
	} else {
		start := i
		for input[i] >= '0' && input[i] <= '9' {
			i += 1
		}
		n, err := strconv.Atoi(input[start:i])
		if err != nil {
			panic(err)
		}
		return ValueInt(n), i
	}
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := strings.Trim(string(content), "\n\r\t ")

	sigV1 := ValueArray([]Value{ValueArray([]Value{ValueInt(2)})})
	sigV2 := ValueArray([]Value{ValueArray([]Value{ValueInt(6)})})

	values := []Value{sigV1, sigV2}

	for _, pair := range strings.Split(contentStr, "\n\n") {
		valuesStr := strings.Split(pair, "\n")
		lhs, _ := parseValue(valuesStr[0], 0)
		rhs, _ := parseValue(valuesStr[1], 0)
		values = append(values, lhs)
		values = append(values, rhs)
	}

	sort.Sort(ValueList(values))

	res := 1
	for i := range values {
		if values[i].compare(sigV1) == EQ || values[i].compare(sigV2) == EQ {
			res *= i + 1
		}
	}

    fmt.Println(res)
}

type ValueList []Value

func (vl ValueList) Len() int {
	return len(vl)
}

func (vl ValueList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}

func (vl ValueList) Less(i, j int) bool {
	return vl[i].compare(vl[j]) == LT
}
