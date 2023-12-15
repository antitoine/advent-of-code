package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ascii(char rune) int64 {
	return int64(char)
}

func hash(input string) int64 {
	var result int64
	for _, char := range input {
		if char == '\n' {
			continue
		}
		if char == '=' || char == '-' {
			break
		}
		result += ascii(char)
		result *= 17
		result = result % 256
	}
	return result
}

type Action int

const (
	OperationSet    Action = iota
	OperationRemove Action = iota
)

type Operation struct {
	Key     string
	KeyHash int64
	Value   int64
	Action  Action
}

func parseInput(input io.Reader) []Operation {
	scanner := bufio.NewScanner(input)

	var result []Operation
	for scanner.Scan() {
		for _, step := range strings.Split(scanner.Text(), ",") {
			if step == "\n" {
				continue
			}
			if strings.HasSuffix(step, "-") {
				key := step[:len(step)-1]
				result = append(result, Operation{
					Key:     key,
					KeyHash: hash(key),
					Action:  OperationRemove,
				})
			} else {
				parts := strings.Split(step, "=")
				if len(parts) != 2 {
					log.Fatalf("Unable to parse step: %s", step)
					return nil
				}
				key := parts[0]
				value, errParsingInt := strconv.ParseInt(parts[1], 10, 64)
				if errParsingInt != nil {
					log.Fatalf("Unable to parse value: %s", parts[1])
					return nil
				}
				result = append(result, Operation{
					Key:     key,
					KeyHash: hash(key),
					Value:   value,
					Action:  OperationSet,
				})
			}
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	return result
}

type Slot struct {
	Idx   int
	Value int64
}

type Box struct {
	SlotByKey map[string]*Slot
	SlotByIdx []*Slot
}

func (b *Box) Set(key string, value int64) {
	if slot, ok := b.SlotByKey[key]; ok {
		slot.Value = value
	} else {
		lastSlotIdx := len(b.SlotByIdx)
		newSlot := &Slot{
			Idx:   lastSlotIdx,
			Value: value,
		}
		b.SlotByKey[key] = newSlot
		b.SlotByIdx = append(b.SlotByIdx, newSlot)
	}
}

func (b *Box) Remove(key string) {
	if slot, ok := b.SlotByKey[key]; ok {
		delete(b.SlotByKey, key)
		newSlotByIdx := make([]*Slot, len(b.SlotByIdx)-1)
		for i := 0; i < slot.Idx; i++ {
			newSlotByIdx[i] = b.SlotByIdx[i]
		}
		for i := slot.Idx + 1; i < len(b.SlotByIdx); i++ {
			b.SlotByIdx[i].Idx--
			newSlotByIdx[i-1] = b.SlotByIdx[i]
		}
		b.SlotByIdx = newSlotByIdx
	}
}

func (b *Box) FocusPower(boxIdx int64) int64 {
	var result int64
	for _, slot := range b.SlotByIdx {
		result += (boxIdx + 1) * int64(slot.Idx+1) * slot.Value
	}
	return result
}

func getResult(input io.Reader) int64 {
	box := make([]Box, 256)
	for i := 0; i < 256; i++ {
		box[i] = Box{
			SlotByKey: make(map[string]*Slot),
		}
	}
	operations := parseInput(input)
	for _, operation := range operations {
		if operation.Action == OperationSet {
			box[operation.KeyHash].Set(operation.Key, operation.Value)
		} else {
			box[operation.KeyHash].Remove(operation.Key)
		}
	}

	var focusPower int64
	for i := int64(0); i < 256; i++ {
		focusPower += box[i].FocusPower(i)
	}
	return focusPower
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
