package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operation interface {
	perform(worryLevel int) int
}

type OperationOperand struct {
	value int
	isVar bool
}

func (op OperationOperand) getValue(worryLevel int) int {
	if op.isVar {
		return worryLevel
	}
	return op.value
}

type OperationOperands struct {
	lhs OperationOperand
	rhs OperationOperand
}

func (ops OperationOperands) apply(worryLevel int, op func(int, int) int) int {
	return op(ops.lhs.getValue(worryLevel), ops.rhs.getValue(worryLevel))
}

type MultiplyOperation struct {
	OperationOperands
}

func (op MultiplyOperation) perform(worryLevel int) int {
	return op.apply(worryLevel, func(lhs, rhs int) int { return lhs * rhs })
}

type AddOperation struct {
	OperationOperands
}

func (op AddOperation) perform(worryLevel int) int {
	return op.apply(worryLevel, func(lhs, rhs int) int { return lhs + rhs })
}

func parseOperand(opStr string) OperationOperand {
	if opStr == "old" {
		return OperationOperand{
			isVar: true,
		}
	} else {
		value, err := strconv.Atoi(opStr)
		if err != nil {
			panic(err)
		}
		return OperationOperand{
			value: value,
		}
	}
}

func parseOperation(opStr string) Operation {
	opStr, found := strings.CutPrefix(opStr, "new = ")
	if !found {
		panic(fmt.Sprintf("Invalid operation: %s", opStr))
	}

	tokens := strings.Split(opStr, " ")
	lhs := parseOperand(tokens[0])
	rhs := parseOperand(tokens[2])

	operands := OperationOperands{
		lhs: lhs,
		rhs: rhs,
	}

	if tokens[1] == "*" {
		return MultiplyOperation{
			operands,
		}
	} else if tokens[1] == "+" {
		return AddOperation{
			operands,
		}
	}

	panic(fmt.Sprintf("Unkown operation: %s", tokens[1]))
}

type TestResult bool

type Test interface {
	test(worryLevel int) TestResult
    getDivisor() int
}

type DivisibleTest struct {
	divisor int
}

func (test DivisibleTest) test(worryLevel int) TestResult {
	return worryLevel%test.divisor == 0
}

func (test DivisibleTest) getDivisor() int {
    return test.divisor
}

func parseTest(testStr string) Test {
	testStr, found := strings.CutPrefix(testStr, "divisible by ")
	if !found {
		panic(fmt.Sprintf("Unknown test: %s", testStr))
	}

	divisor, err := strconv.Atoi(testStr)
	if err != nil {
		panic(err)
	}

	return DivisibleTest{
		divisor: divisor,
	}
}

type MonkeyId int

type Action interface {
	getThrowsTo(TestResult) MonkeyId
}

type IfAction struct {
	ifTrue  MonkeyId
	ifFalse MonkeyId
}

func (action IfAction) getThrowsTo(result TestResult) MonkeyId {
	if result {
		return action.ifTrue
	} else {
		return action.ifFalse
	}
}

func parseAction(actStrs []string) Action {
	ifTrueStr, found := strings.CutPrefix(actStrs[0], "If true: throw to monkey ")
	if !found {
		panic(fmt.Sprintf("Unknown action case: %s", actStrs[0]))
	}

	ifFalseStr, found := strings.CutPrefix(actStrs[1], "If false: throw to monkey ")
	if !found {
		panic(fmt.Sprintf("Unknown action case: %s", actStrs[1]))
	}

	ifTrue, err := strconv.Atoi(ifTrueStr)
	if err != nil {
		panic(err)
	}

	ifFalse, err := strconv.Atoi(ifFalseStr)
	if err != nil {
		panic(err)
	}

	return IfAction{
		ifTrue:  MonkeyId(ifTrue),
		ifFalse: MonkeyId(ifFalse),
	}
}

type Monkey struct {
	operation Operation
	test      Test
	action    Action

	items []int
}

func parseMonkey(lines []string) Monkey {
	startingItemsStr, found := strings.CutPrefix(lines[0], "Starting items: ")
	if !found {
		panic(fmt.Sprintf("Unknown format: %s", startingItemsStr))
	}

	items := []int{}

	for _, itemStr := range strings.Split(startingItemsStr, ", ") {
		item, err := strconv.Atoi(itemStr)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
	}

	operationStr, found := strings.CutPrefix(lines[1], "Operation: ")
	if !found {
		panic(fmt.Sprintf("Unknown format: %s", operationStr))
	}

	operation := parseOperation(operationStr)

	testStr, found := strings.CutPrefix(lines[2], "Test: ")
	if !found {
		panic(fmt.Sprintf("Unknown format: %s", testStr))
	}

	test := parseTest(testStr)

	action := parseAction(lines[3:])

	return Monkey{
		operation: operation,
		test:      test,
		action:    action,
		items:     items,
	}
}

func (m *Monkey) hasItems() bool {
	return len(m.items) != 0
}

func (m *Monkey) inspectNext(bound int) (int, MonkeyId) {
	worryLevel := m.items[0]
	m.items = m.items[1:]

    // fmt.Printf("Inspecting item with worry level %d\n", worryLevel)

	worryLevel = m.operation.perform(worryLevel)
    // fmt.Printf("Worry level increasese to %d\n", worryLevel)
	// worryLevel = worryLevel / 3
    worryLevel = worryLevel % bound
    // fmt.Printf("Worry level decreases to %d\n", worryLevel)

	testResult := m.test.test(worryLevel)
    // fmt.Printf("Test result: %v\n", testResult)
	throwTo := m.action.getThrowsTo(testResult)
    // fmt.Printf("Throwing item to monkey %v\n", throwTo)

	return worryLevel, throwTo
}

func (m *Monkey) receivesItem(worryLevel int) {
	m.items = append(m.items, worryLevel)
}

func (m *Monkey) print() string {
	items := []string{}
	for _, item := range m.items {
		items = append(items, fmt.Sprint(item))
	}
	return strings.Join(items, ", ")
}

type Situation struct {
	monkeys  []*Monkey
	inspects map[MonkeyId]int
}

func makeSituation() Situation {
	return Situation{
		monkeys:  []*Monkey{},
		inspects: map[MonkeyId]int{},
	}
}

func (s *Situation) round() {
    bound := 1

    for _, monkey := range s.monkeys {
        bound *= monkey.test.getDivisor()
    }

	for id, monkey := range s.monkeys {
		for monkey.hasItems() {
			s.inspects[MonkeyId(id)] += 1
			worryLevel, throwTo := monkey.inspectNext(bound)
			s.monkeys[throwTo].receivesItem(worryLevel)
		}
	}
}

func (s *Situation) addMonkey(lines []string) {
	id := parseMonkeyId(lines[0])
    if int(id) != len(s.monkeys) {
        panic("Invalid monkey id")
    }
	monkey := parseMonkey(lines[1:])
	s.monkeys = append(s.monkeys, &monkey)
}

func (s *Situation) print() string {
	lines := []string{}
	for mId, monkey := range s.monkeys {
		lines = append(lines, fmt.Sprintf("Monkey %v: %s", mId, monkey.print()))
	}
	return strings.Join(lines, "\n")
}

func parseMonkeyId(line string) MonkeyId {
	idStr, found := strings.CutPrefix(line, "Monkey ")
	if !found {
		panic(fmt.Sprintf("Unknown format: %s", idStr))
	}

	idStr = idStr[:len(idStr)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}

	return MonkeyId(id)
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := strings.Trim(string(content), "\n\r\t ")

	currentLines := []string{}

	s := makeSituation()

	for _, line := range strings.Split(contentStr, "\n") {
		if line == "" {
			s.addMonkey(currentLines)
			currentLines = []string{}
			continue
		}

		line = strings.Trim(line, " ")
		currentLines = append(currentLines, line)
	}
	s.addMonkey(currentLines)

	// fmt.Println(s.print())

	for i := 0; i < 10000; i += 1 {
		s.round()
        // fmt.Printf("After round %d:\n%s\n\n", i, s.print())
	}

	fmt.Println(s.inspects)
}
