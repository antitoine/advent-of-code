package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Condition struct {
	key      string
	operator string
	value    int64
}

func (c Condition) String() string {
	return c.key + c.operator + strconv.FormatInt(c.value, 10)
}

func (c Condition) Test(instruction Instruction) bool {
	switch c.key {
	case "x":
		switch c.operator {
		case "<":
			return instruction.x < c.value
		case ">":
			return instruction.x > c.value
		}
	case "m":
		switch c.operator {
		case "<":
			return instruction.m < c.value
		case ">":
			return instruction.m > c.value
		}
	case "a":
		switch c.operator {
		case "<":
			return instruction.a < c.value
		case ">":
			return instruction.a > c.value
		}
	case "s":
		switch c.operator {
		case "<":
			return instruction.s < c.value
		case ">":
			return instruction.s > c.value
		}
	}
	log.Fatalf("Unable to check if instruction is approved: %s", c.String())
	return false
}

type Rule struct {
	condition Condition
	ifTrue    *Step
	ifFalse   *Step
}

var conditionRegex = regexp.MustCompile(`^([xmas])([<>])(\d+)$`)

func NewRule(conditionStr string, ifTrueStr string, ifFalseStr string, otherWorkflows RawWorkflows) *Rule {
	conditionParts := conditionRegex.FindStringSubmatch(conditionStr)
	if len(conditionParts) != 4 {
		log.Fatalf("Unable to parse condition: %s", conditionStr)
	}
	value, errParsingValue := strconv.ParseInt(conditionParts[3], 10, 64)
	if errParsingValue != nil {
		log.Fatalf("Unable to parse condition value in: %s (%v)", conditionStr, errParsingValue)
	}
	condition := Condition{
		key:      conditionParts[1],
		operator: conditionParts[2],
		value:    value,
	}
	return &Rule{
		condition: condition,
		ifTrue:    NewStep(ifTrueStr, otherWorkflows),
		ifFalse:   NewStep(ifFalseStr, otherWorkflows),
	}
}

func (r *Rule) String() string {
	return "if " + r.condition.String() + " { " + r.ifTrue.String() + " } else { " + r.ifFalse.String() + " }"
}

type Step struct {
	approved *bool
	rule     *Rule
}

func NewStep(workflowContent string, otherRawWorkflows RawWorkflows) *Step {
	stepList := strings.Split(workflowContent, ",")
	stepStr := stepList[0]
	switch stepStr {
	case "R":
		value := false
		return &Step{
			approved: &value,
		}
	case "A":
		value := true
		return &Step{
			approved: &value,
		}
	default:
		if nextWorkflowContent, ok := otherRawWorkflows[stepStr]; ok {
			return NewStep(nextWorkflowContent, otherRawWorkflows)
		}
		conditionParts := strings.Split(stepStr, ":")
		if len(conditionParts) != 2 {
			log.Fatalf("Unable to parse step for condition: %s", stepStr)
		}
		if len(stepList) == 1 {
			log.Fatalf("Unable to parse step for next step: %s", stepStr)
		}
		nextStepStr := strings.Join(stepList[1:], ",")
		return &Step{
			rule: NewRule(conditionParts[0], conditionParts[1], nextStepStr, otherRawWorkflows),
		}
	}
	return nil
}

func (s *Step) String() string {
	if s.approved != nil {
		if *s.approved {
			return "A"
		}
		return "R"
	}
	return s.rule.String()
}

func (s *Step) IsApproved(instruction Instruction) bool {
	if s.approved != nil {
		return *s.approved
	}
	if s.rule == nil {
		log.Fatalf("Unable to check step, approved and rule are nil: %s", s.String())
	}
	if s.rule.condition.Test(instruction) {
		return s.rule.ifTrue.IsApproved(instruction)
	}
	return s.rule.ifFalse.IsApproved(instruction)
}

type RawWorkflows map[string]string

var ruleRegex = regexp.MustCompile(`^([a-zA-Z0-9]+)\{(.*)}$`)

func parseRawWorkflow(line string) (string, string) {
	ruleMatches := ruleRegex.FindStringSubmatch(line)
	if len(ruleMatches) != 3 {
		log.Fatalf("Unable to parse rule for finding the key: %s", line)
	}
	return ruleMatches[1], ruleMatches[2]
}

type Instruction struct {
	x int64
	m int64
	a int64
	s int64
}

func (i Instruction) Sum() int64 {
	return i.x + i.m + i.a + i.s
}

var instructionRegex = regexp.MustCompile(`^\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}$`)

func parseInstruction(instructionStr string) Instruction {
	instructionParts := instructionRegex.FindStringSubmatch(instructionStr)
	if len(instructionParts) != 5 {
		log.Fatalf("Unable to parse instruction: %s", instructionStr)
	}

	x, errParsingX := strconv.ParseInt(instructionParts[1], 10, 64)
	if errParsingX != nil {
		log.Fatalf("Unable to parse instruction x: %s (%v)", instructionStr, errParsingX)
	}

	m, errParsingM := strconv.ParseInt(instructionParts[2], 10, 64)
	if errParsingM != nil {
		log.Fatalf("Unable to parse instruction m: %s (%v)", instructionStr, errParsingM)
	}

	a, errParsingA := strconv.ParseInt(instructionParts[3], 10, 64)
	if errParsingA != nil {
		log.Fatalf("Unable to parse instruction a: %s (%v)", instructionStr, errParsingA)
	}

	s, errParsingS := strconv.ParseInt(instructionParts[4], 10, 64)
	if errParsingS != nil {
		log.Fatalf("Unable to parse instruction s: %s (%v)", instructionStr, errParsingS)
	}

	return Instruction{
		x: x,
		m: m,
		a: a,
		s: s,
	}
}

func parseInput(input io.Reader) (RawWorkflows, []Instruction) {
	scanner := bufio.NewScanner(input)

	rawWorkflows := make(RawWorkflows)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		workflowKey, workflowContent := parseRawWorkflow(line)
		rawWorkflows[workflowKey] = workflowContent
	}

	var instructions []Instruction
	for scanner.Scan() {
		line := scanner.Text()
		instructions = append(instructions, parseInstruction(line))
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return rawWorkflows, instructions
}

const firstWorkflowKey = "in"

func getResult(input io.Reader) int64 {
	rawWorkflows, instructions := parseInput(input)

	firstWorkflowContent, firstWorkflowFound := rawWorkflows[firstWorkflowKey]
	if !firstWorkflowFound {
		log.Fatalf("Unable to find the first workflow: %s", firstWorkflowKey)
	}

	firstStep := NewStep(firstWorkflowContent, rawWorkflows)
	//log.Printf("First step: \n   %s", firstStep.String())

	log.Printf("Count of instructions: %d", len(instructions))
	count := int64(0)
	for _, instruction := range instructions {
		if firstStep.IsApproved(instruction) {
			count += instruction.Sum()
		}
	}

	return count
}

func loadFile() *os.File {
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	return inputFile
}

func main() {
	start := time.Now()
	inputFile := loadFile()
	defer inputFile.Close()

	result := getResult(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
