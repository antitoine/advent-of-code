package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Pulse int

const (
	High Pulse = iota
	Low  Pulse = iota
)

type Input struct {
	pulse Pulse
	from  ModuleId
}

func NewInput(pulse Pulse, from ModuleId) Input {
	return Input{pulse: pulse, from: from}
}

type Output struct {
	pulse Pulse
	from  ModuleId
	to    []ModuleId
}

func (o Output) GetNextInputs() (map[ModuleId]Input, int64, int64) {
	var countOfLowPulses int64
	var countOfHighPulses int64
	if o.pulse == Low {
		countOfLowPulses = int64(len(o.to))
	} else {
		countOfHighPulses = int64(len(o.to))
	}
	inputs := make(map[ModuleId]Input, len(o.to))
	for _, moduleId := range o.to {
		inputs[moduleId] = NewInput(o.pulse, o.from)
	}
	return inputs, countOfLowPulses, countOfHighPulses
}

type ModuleId string

type Module interface {
	Handle(input Input) Output
	GetModuleOutputIds() []ModuleId
	GetId() ModuleId
}

type FlipFlop struct {
	id        ModuleId
	isOnState bool
	to        []ModuleId
}

func NewFlipFlop(id string, linkedTo []ModuleId) *FlipFlop {
	module := &FlipFlop{
		to: linkedTo,
		id: ModuleId(id),
	}
	return module
}

func (f *FlipFlop) Handle(input Input) Output {
	if input.pulse == High {
		return Output{}
	}
	f.isOnState = !f.isOnState
	var result Pulse
	if f.isOnState {
		result = High
	} else {
		result = Low
	}
	return Output{
		pulse: result,
		from:  f.id,
		to:    f.to,
	}
}

func (f *FlipFlop) GetModuleOutputIds() []ModuleId {
	return f.to
}

func (f *FlipFlop) GetId() ModuleId {
	return f.id
}

type Conjunction struct {
	alreadyReceived map[ModuleId]Pulse
	to              []ModuleId
	id              ModuleId
}

type AttachModuleId func(ModuleId)

func NewConjunction(id string, linkedTo []ModuleId) (*Conjunction, AttachModuleId) {
	module := &Conjunction{
		alreadyReceived: make(map[ModuleId]Pulse),
		to:              linkedTo,
		id:              ModuleId(id),
	}
	return module, func(moduleId ModuleId) {
		module.alreadyReceived[moduleId] = Low
	}
}

func (c *Conjunction) Handle(input Input) Output {
	c.alreadyReceived[input.from] = input.pulse
	onlyHigh := true
	for _, pulse := range c.alreadyReceived {
		if pulse == Low {
			onlyHigh = false
			break
		}
	}
	var result Pulse
	if onlyHigh {
		result = Low
	} else {
		result = High
	}
	return Output{
		pulse: result,
		from:  c.id,
		to:    c.to,
	}
}

func (c *Conjunction) GetModuleOutputIds() []ModuleId {
	return c.to
}

func (c *Conjunction) GetId() ModuleId {
	return c.id
}

type Broadcast struct {
	to []ModuleId
	id ModuleId
}

func NewBroadcast(linkedTo []ModuleId) *Broadcast {
	module := &Broadcast{
		id: ModuleId("broadcaster"),
		to: linkedTo,
	}
	return module
}

func (b *Broadcast) Handle(input Input) Output {
	return Output{
		pulse: input.pulse,
		from:  b.id,
		to:    b.to,
	}
}

func (b *Broadcast) GetModuleOutputIds() []ModuleId {
	return b.to
}

func (b *Broadcast) GetId() ModuleId {
	return b.id
}

type Sand struct {
	id                ModuleId
	parent            ModuleId
	countOfLowPulses  int64
	countOfHighPulses int64
}

func NewSand() (*Sand, AttachModuleId) {
	module := &Sand{
		id: ModuleId("rx"),
	}
	return module, func(moduleId ModuleId) {
		module.parent = moduleId
	}
}

func (s *Sand) Handle(input Input) Output {
	if input.pulse == Low {
		s.countOfLowPulses++
	} else {
		s.countOfHighPulses++
	}
	return Output{}
}

func (s *Sand) GetModuleOutputIds() []ModuleId {
	return nil
}

func (s *Sand) GetId() ModuleId {
	return s.id
}

func parseModuleIds(modulesStr string) []ModuleId {
	modules := strings.Split(modulesStr, ", ")
	moduleIds := make([]ModuleId, len(modules))
	for i, module := range modules {
		moduleIds[i] = ModuleId(module)
	}
	return moduleIds
}

func parseModule(line string) (Module, AttachModuleId) {
	var module Module
	var declareModuleId AttachModuleId
	if nextFlipFlopStr, foundFlipFlopPrefix := strings.CutPrefix(line, "%"); foundFlipFlopPrefix {
		parts := strings.Split(nextFlipFlopStr, " -> ")
		if len(parts) != 2 {
			log.Fatalf("Unable to parse flip-flop module line: %s", line)
		}
		module = NewFlipFlop(parts[0], parseModuleIds(parts[1]))
	} else if nextConjunctionStr, foundConjunctionPrefix := strings.CutPrefix(line, "&"); foundConjunctionPrefix {
		parts := strings.Split(nextConjunctionStr, " -> ")
		if len(parts) != 2 {
			log.Fatalf("Unable to parse conjunction module line: %s", line)
		}
		module, declareModuleId = NewConjunction(parts[0], parseModuleIds(parts[1]))
	} else {
		log.Fatalf("Unable to parse module line: %s", line)
	}
	return module, declareModuleId
}

func parseBroadcast(line string) *Broadcast {
	modulesStr, foundPrefix := strings.CutPrefix(line, "broadcaster -> ")
	if !foundPrefix {
		log.Fatalf("Unable to parse line: %s", line)
	}
	return NewBroadcast(parseModuleIds(modulesStr))
}

func parseInput(input io.Reader) (*Broadcast, map[ModuleId]Module, *Sand) {
	scanner := bufio.NewScanner(input)

	var broadcast *Broadcast
	modules := make(map[ModuleId]Module)

	sand, setParentModule := NewSand()
	modules[sand.GetId()] = sand

	waitingAttachingParentModule := make(map[ModuleId]AttachModuleId)
	waitingAttachingParentModule[sand.GetId()] = setParentModule

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "broadcaster ->") {
			broadcast = parseBroadcast(line)
		} else {
			module, declareModuleId := parseModule(line)
			modules[module.GetId()] = module
			if declareModuleId != nil {
				waitingAttachingParentModule[module.GetId()] = declareModuleId
			}
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	// Setup conjunctions modules
	for moduleId, module := range modules {
		for _, outputModuleId := range module.GetModuleOutputIds() {
			if declareModuleIdFunc, found := waitingAttachingParentModule[outputModuleId]; found {
				declareModuleIdFunc(moduleId)
			}
		}
	}

	return broadcast, modules, sand
}

func TriggerOnce(broadcast *Broadcast, modules map[ModuleId]Module, moduleHighInputWatch *ModuleId) (int64, int64, []ModuleId) {
	firstInput := NewInput(Low, broadcast.GetId())
	countOfLowPulses := int64(1)
	countOfHighPulses := int64(0)

	var moduleHighInputDetected []ModuleId

	tasks := make([][]Output, 1)
	tasks = append(tasks, []Output{broadcast.Handle(firstInput)})

	for len(tasks) > 0 {
		nextTask := tasks[0]
		tasks = tasks[1:]
		var nextOutput []Output
		for _, output := range nextTask {
			nextInputs, newCountOfLowPulses, newCountOfHighPulses := output.GetNextInputs()
			countOfLowPulses += newCountOfLowPulses
			countOfHighPulses += newCountOfHighPulses
			for forModuleId, input := range nextInputs {
				if moduleHighInputWatch != nil && *moduleHighInputWatch == forModuleId && input.pulse == High {
					moduleHighInputDetected = append(moduleHighInputDetected, output.from)
				}
				nextModule, found := modules[forModuleId]
				if !found {
					continue
				}
				receivedOutput := nextModule.Handle(input)
				if len(receivedOutput.to) > 0 {
					nextOutput = append(nextOutput, receivedOutput)
				}
			}
		}
		if len(nextOutput) > 0 {
			tasks = append(tasks, nextOutput)
		}
	}

	return countOfLowPulses, countOfHighPulses, moduleHighInputDetected
}

func getResultForPart1(text io.Reader) int64 {
	broadcast, modules, _ := parseInput(text)
	var countOfLowPulses, countOfHighPulses int64
	for i := int64(0); i < 1000; i++ {
		newCountOfLowPulses, newCountOfHighPulses, _ := TriggerOnce(broadcast, modules, nil)
		countOfLowPulses += newCountOfLowPulses
		countOfHighPulses += newCountOfHighPulses
	}
	log.Printf("Count of low pulses: %d", countOfLowPulses)
	log.Printf("Count of high pulses: %d", countOfHighPulses)
	return countOfLowPulses * countOfHighPulses
}

func allModulesStepsDetected(moduleIdsLowAfter map[ModuleId]int64) (bool, []int64) {
	var result []int64
	for _, after := range moduleIdsLowAfter {
		if after == -1 {
			return false, nil
		}
		result = append(result, after)
	}
	return true, result
}

func greatestCommonDivisor(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func leastCommonMultiple(a, b int64) int64 {
	return a * b / greatestCommonDivisor(a, b)
}

func getResultForPart2(text io.Reader) int64 {
	broadcast, modules, sand := parseInput(text)

	log.Printf("Broadcast: %#v", *broadcast)
	log.Printf("Modules: %v", modules)
	log.Printf("Sand: %#v", *sand)

	parentSandModuleId := sand.parent
	parentSandModule, foundParentSandModule := modules[parentSandModuleId]
	if !foundParentSandModule {
		log.Fatalf("Unable to find parent sand module: %s", parentSandModuleId)
	}

	log.Printf("Parent Sand module: %#v", parentSandModule)

	moduleIdsWatching := make(map[ModuleId]int64)
	for moduleId := range parentSandModule.(*Conjunction).alreadyReceived {
		moduleIdsWatching[moduleId] = -1
	}

	log.Printf("Waiting modules IDs to be high: %#v", moduleIdsWatching)
	isAllDetected, afters := allModulesStepsDetected(moduleIdsWatching)
	for i := int64(1); !isAllDetected; i++ {
		_, _, detected := TriggerOnce(broadcast, modules, &sand.parent)
		for _, moduleId := range detected {
			if moduleIdsWatching[moduleId] == -1 {
				log.Printf("Module %s is high after %d steps", moduleId, i)
				moduleIdsWatching[moduleId] = i
			}
		}
		isAllDetected, afters = allModulesStepsDetected(moduleIdsWatching)
	}

	log.Printf("All modules are high after %#v steps", afters)

	minRequiredStep := int64(1)
	for _, steps := range afters {
		minRequiredStep = leastCommonMultiple(minRequiredStep, steps)
	}
	return minRequiredStep
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

	result := getResultForPart2(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
