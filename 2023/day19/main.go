package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
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

func (c Condition) Apply(currentRange ApprovedInstructionRange) (ApprovedInstructionRange, ApprovedInstructionRange) {
	ifTrueRange := currentRange
	ifFalseRange := currentRange
	switch c.key {
	case "x":
		switch c.operator {
		case "<":
			// x<value
			ifTrueRange.x.max = c.value - 1
			// x>=value
			ifFalseRange.x.min = c.value
		case ">":
			ifTrueRange.x.min = c.value + 1
			ifFalseRange.x.max = c.value
		}
	case "m":
		switch c.operator {
		case "<":
			ifTrueRange.m.max = c.value - 1
			ifFalseRange.m.min = c.value
		case ">":
			ifTrueRange.m.min = c.value + 1
			ifFalseRange.m.max = c.value
		}
	case "a":
		switch c.operator {
		case "<":
			ifTrueRange.a.max = c.value - 1
			ifFalseRange.a.min = c.value
		case ">":
			ifTrueRange.a.min = c.value + 1
			ifFalseRange.a.max = c.value
		}
	case "s":
		switch c.operator {
		case "<":
			ifTrueRange.s.max = c.value - 1
			ifFalseRange.s.min = c.value
		case ">":
			ifTrueRange.s.min = c.value + 1
			ifFalseRange.s.max = c.value
		}
	}
	return ifTrueRange, ifFalseRange
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

func (r *Rule) String(indent string) string {
	return indent + "if " + r.condition.String() + " {\n" +
		r.ifTrue.String(indent+"  ") + "\n" +
		indent + "} else {\n" +
		r.ifFalse.String(indent+"  ") + "\n" +
		indent + "}"
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

func (s *Step) String(indent string) string {
	if s.approved != nil {
		if *s.approved {
			return indent + "A"
		}
		return indent + "R"
	}
	return s.rule.String(indent)
}

type Range struct {
	min int64
	max int64
}

type ApprovedInstructionRange struct {
	x Range
	m Range
	a Range
	s Range
}

const approvedInstructionRangeMin = 1
const approvedInstructionRangeMax = 4000

func NewApprovedInstructionRange() ApprovedInstructionRange {
	return ApprovedInstructionRange{
		x: Range{
			min: approvedInstructionRangeMin,
			max: approvedInstructionRangeMax,
		},
		m: Range{
			min: approvedInstructionRangeMin,
			max: approvedInstructionRangeMax,
		},
		a: Range{
			min: approvedInstructionRangeMin,
			max: approvedInstructionRangeMax,
		},
		s: Range{
			min: approvedInstructionRangeMin,
			max: approvedInstructionRangeMax,
		},
	}
}

func (r ApprovedInstructionRange) String() string {
	return "{x=" + strconv.FormatInt(r.x.min, 10) + "," + strconv.FormatInt(r.x.max, 10) +
		",m=" + strconv.FormatInt(r.m.min, 10) + "," + strconv.FormatInt(r.m.max, 10) +
		",a=" + strconv.FormatInt(r.a.min, 10) + "," + strconv.FormatInt(r.a.max, 10) +
		",s=" + strconv.FormatInt(r.s.min, 10) + "," + strconv.FormatInt(r.s.max, 10) + "}"
}

func (s *Step) ComputeListOfApprovedInstructionRange(currentRange ApprovedInstructionRange) []ApprovedInstructionRange {
	if s.approved != nil {
		if *s.approved {
			return []ApprovedInstructionRange{currentRange}
		}
		return nil
	}
	if s.rule == nil {
		log.Fatalf("Unable to compute list of approved instruction range, approved and rule are nil: %s", s.String(""))
	}
	ifTrueRange, ifFalseRange := s.rule.condition.Apply(currentRange)
	var result []ApprovedInstructionRange
	result = append(result, s.rule.ifTrue.ComputeListOfApprovedInstructionRange(ifTrueRange)...)
	result = append(result, s.rule.ifFalse.ComputeListOfApprovedInstructionRange(ifFalseRange)...)
	return result
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

func parseInput(input io.Reader) RawWorkflows {
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

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return rawWorkflows
}

const firstWorkflowKey = "in"

func ComputeBreakPointForRanges(ranges []Range) []int64 {
	result := make(map[int64]bool)
	result[approvedInstructionRangeMin] = true
	result[approvedInstructionRangeMax+1] = true
	for _, r := range ranges {
		result[r.min] = true
		result[r.max+1] = true
	}
	var resultSlice []int64
	for k := range result {
		resultSlice = append(resultSlice, k)
	}
	slices.Sort(resultSlice)
	return resultSlice
}

func getResult(input io.Reader) int64 {
	rawWorkflows := parseInput(input)

	firstWorkflowContent, firstWorkflowFound := rawWorkflows[firstWorkflowKey]
	if !firstWorkflowFound {
		log.Fatalf("Unable to find the first workflow: %s", firstWorkflowKey)
	}

	firstStep := NewStep(firstWorkflowContent, rawWorkflows)
	//log.Printf("First step:\n%s", firstStep.String("  "))

	approvedInstructionRangeList := firstStep.ComputeListOfApprovedInstructionRange(NewApprovedInstructionRange())
	//log.Printf("Found %d approved instruction ranges", len(approvedInstructionRangeList))
	//for _, approvedInstructionRange := range approvedInstructionRangeList {
	//	log.Printf("Approved instruction range: %s", approvedInstructionRange.String())
	//}

	var xRanges []Range
	var mRanges []Range
	var aRanges []Range
	var sRanges []Range
	for _, approvedInstructionRange := range approvedInstructionRangeList {
		xRanges = append(xRanges, approvedInstructionRange.x)
		mRanges = append(mRanges, approvedInstructionRange.m)
		aRanges = append(aRanges, approvedInstructionRange.a)
		sRanges = append(sRanges, approvedInstructionRange.s)
	}

	xBreakPoints := ComputeBreakPointForRanges(xRanges)
	//log.Printf("Found %d x break points: %v", len(xBreakPoints), xBreakPoints)
	mBreakPoints := ComputeBreakPointForRanges(mRanges)
	//log.Printf("Found %d m break points: %v", len(mBreakPoints), mBreakPoints)
	aBreakPoints := ComputeBreakPointForRanges(aRanges)
	//log.Printf("Found %d a break points: %v", len(aBreakPoints), aBreakPoints)
	sBreakPoints := ComputeBreakPointForRanges(sRanges)
	//log.Printf("Found %d s break points: %v", len(sBreakPoints), sBreakPoints)

	var result int64
	for xIdx := 0; xIdx < len(xBreakPoints)-1; xIdx++ {
		xMin := xBreakPoints[xIdx]
		xMax := xBreakPoints[xIdx+1]

		var stillValidInstructionsAfterX []ApprovedInstructionRange
		for _, approvedInstructionRange := range approvedInstructionRangeList {
			if approvedInstructionRange.x.min <= xMin && xMax <= approvedInstructionRange.x.max+1 {
				stillValidInstructionsAfterX = append(stillValidInstructionsAfterX, approvedInstructionRange)
			}
		}
		if len(stillValidInstructionsAfterX) == 0 {
			continue
		}

		for mIdx := 0; mIdx < len(mBreakPoints)-1; mIdx++ {
			mMin := mBreakPoints[mIdx]
			mMax := mBreakPoints[mIdx+1]

			var stillValidInstructionsAfterM []ApprovedInstructionRange
			for _, approvedInstructionRange := range stillValidInstructionsAfterX {
				if approvedInstructionRange.m.min <= mMin && mMax <= approvedInstructionRange.m.max+1 {
					stillValidInstructionsAfterM = append(stillValidInstructionsAfterM, approvedInstructionRange)
				}
			}
			if len(stillValidInstructionsAfterM) == 0 {
				continue
			}

			for aIdx := 0; aIdx < len(aBreakPoints)-1; aIdx++ {
				aMin := aBreakPoints[aIdx]
				aMax := aBreakPoints[aIdx+1]

				var stillValidInstructionsAfterA []ApprovedInstructionRange
				for _, approvedInstructionRange := range stillValidInstructionsAfterM {
					if approvedInstructionRange.a.min <= aMin && aMax <= approvedInstructionRange.a.max+1 {
						stillValidInstructionsAfterA = append(stillValidInstructionsAfterA, approvedInstructionRange)
					}
				}
				if len(stillValidInstructionsAfterA) == 0 {
					continue
				}

				for sIdx := 0; sIdx < len(sBreakPoints)-1; sIdx++ {
					sMin := sBreakPoints[sIdx]
					sMax := sBreakPoints[sIdx+1]

					foundS := false
					for _, approvedInstructionRange := range stillValidInstructionsAfterA {
						if approvedInstructionRange.s.min <= sMin && sMax <= approvedInstructionRange.s.max+1 {
							foundS = true
							break
						}
					}
					if foundS {
						result += (xMax - xMin) * (mMax - mMin) * (aMax - aMin) * (sMax - sMin)
					}
				}
			}
		}
	}

	return result
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
