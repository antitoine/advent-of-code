package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Opcode int8

const (
	adv Opcode = 0
	bxl Opcode = 1
	bst Opcode = 2
	jnz Opcode = 3
	bxc Opcode = 4
	out Opcode = 5
	bdv Opcode = 6
	cdv Opcode = 7
)

type Operand int8

type Register rune

const (
	regA Register = 'A'
	regB Register = 'B'
	regC Register = 'C'
)

var registerLineRegex = regexp.MustCompile(`Register ([A-C]): (\d+)`)

func parseRegisterLine(line string) (Register, int64) {
	matches := registerLineRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		log.Fatalf("Unable to parse register line: %s", line)
	}
	register := Register(rune(matches[1][0]))
	value, errValue := strconv.ParseInt(matches[2], 10, 64)
	if errValue != nil {
		log.Fatalf("Unable to parse register value: %v", errValue)
	}
	return register, value
}

type Instruction int8

func parseProgramLine(line string) []Instruction {
	lineParts := strings.Split(strings.TrimPrefix(line, "Program: "), ",")
	if len(lineParts)%2 != 0 {
		log.Fatalf("Invalid program line: %s", line)
	}
	instructions := make([]Instruction, len(lineParts))
	for i, part := range lineParts {
		value, errValue := strconv.ParseInt(part, 10, 8)
		if errValue != nil {
			log.Fatalf("Unable to parse program value: %v", errValue)
		}
		instructions[i] = Instruction(value)
	}
	return instructions
}

func parseInput(input io.Reader) (map[Register]int64, []Instruction) {
	scanner := bufio.NewScanner(input)

	registers := make(map[Register]int64)
	var instructions []Instruction
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Register") {
			register, value := parseRegisterLine(line)
			registers[register] = value
		} else if strings.HasPrefix(line, "Program") {
			instructions = parseProgramLine(line)
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return registers, instructions
}

type Program struct {
	registers    map[Register]int64
	instructions []Instruction
	pointer      int
	outs         []Instruction
}

func (p *Program) comboOperand(op Operand) int64 {
	switch op {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return p.registers[regA]
	case 5:
		return p.registers[regB]
	case 6:
		return p.registers[regC]
	default:
		log.Fatalf("Invalid combo operand: %d", op)
	}
	return 0
}

func (p *Program) process() bool {
	if p.pointer < 0 || p.pointer >= len(p.instructions) {
		return true
	}
	opcode := Opcode(p.instructions[p.pointer])
	operand := Operand(p.instructions[p.pointer+1])
	switch opcode {
	case adv:
		p.registers[regA] = p.registers[regA] / int64(math.Pow(2, float64(p.comboOperand(operand))))
	case bxl:
		p.registers[regB] = p.registers[regB] ^ int64(operand)
	case bst:
		p.registers[regB] = p.comboOperand(operand) % 8
	case jnz:
		if p.registers[regA] != 0 {
			p.pointer = int(operand)
			return false
		}
	case bxc:
		p.registers[regB] = p.registers[regB] ^ p.registers[regC]
	case out:
		p.outs = append(p.outs, Instruction(p.comboOperand(operand)%8))
	case bdv:
		p.registers[regB] = p.registers[regA] / int64(math.Pow(2, float64(p.comboOperand(operand))))
	case cdv:
		p.registers[regC] = p.registers[regA] / int64(math.Pow(2, float64(p.comboOperand(operand))))
	}
	p.pointer += 2
	return false
}

func newRegisters(registers map[Register]int64) map[Register]int64 {
	nr := make(map[Register]int64)
	for register, value := range registers {
		nr[register] = value
	}
	return nr
}

func getResult(input io.Reader) int64 {
	registers, instructions := parseInput(input)
	var a int64
	for itr := len(instructions) - 1; itr >= 0; itr-- {
		a <<= 3
		for {
			updatedRegisters := newRegisters(registers)
			updatedRegisters[regA] = a
			program := Program{updatedRegisters, instructions, 0, nil}
			for !program.process() {
			}
			if slices.Equal(program.outs, instructions[itr:]) {
				break
			}
			a++
		}
	}

	return a
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
