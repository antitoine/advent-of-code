package main

import (
	"bufio"
	"gonum.org/v1/gonum/mat"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Set map[string]struct{}

func (s Set) Add(value string) {
	s[value] = struct{}{}
}

func (s Set) Contains(value string) bool {
	_, ok := s[value]
	return ok
}

type Graph map[string]Set

func (g Graph) AddEdge(from, to string) {
	if _, ok := g[from]; !ok {
		g[from] = make(Set)
	}
	g[from].Add(to)
	if _, ok := g[to]; !ok {
		g[to] = make(Set)
	}
	g[to].Add(from)
}

type Keys []string

type Matrix [][]int

func (g Graph) GetMatrixForKeys(keys Keys) (*mat.SymDense, *mat.SymDense, *mat.SymDense) {
	degreeMatrix := mat.NewSymDense(len(keys), nil)
	adjacencyMatrix := mat.NewSymDense(len(keys), nil)
	laplacianMatrix := mat.NewSymDense(len(keys), nil)
	for i, key := range keys {
		degreeMatrix.SetSym(i, i, float64(len(g[key])))
		for j, otherKey := range keys {
			if g[key].Contains(otherKey) {
				adjacencyMatrix.SetSym(i, j, 1.0)
			}
		}
	}
	for i := range keys {
		for j := range keys {
			if i != j {
				laplacianMatrix.SetSym(i, j, degreeMatrix.At(i, j)-adjacencyMatrix.At(i, j))
			}
		}
	}
	return degreeMatrix, adjacencyMatrix, laplacianMatrix
}

func parseLine(line string) (string, []string) {
	parts := strings.Split(line, ": ")
	if len(parts) != 2 {
		log.Fatalf("Unable to parse line: %s", line)
	}
	otherParts := strings.Split(parts[1], " ")
	return parts[0], otherParts
}

func parseInput(input io.Reader) (Keys, Graph) {
	scanner := bufio.NewScanner(input)

	graph := make(Graph)
	for scanner.Scan() {
		key, connectedKeys := parseLine(scanner.Text())
		for _, connectedKey := range connectedKeys {
			graph.AddEdge(key, connectedKey)
		}
	}

	if errScanningFile := scanner.Err(); errScanningFile != nil {
		log.Fatalf("Unable to scan the input file correctly: %v", errScanningFile)
	}

	keys := make(Keys, len(graph))
	for key := range graph {
		keys = append(keys, key)
	}

	return keys, graph
}

func getResult(input io.Reader) int {
	keys, graph := parseInput(input)

	_, _, laplacianMatrix := graph.GetMatrixForKeys(keys)

	var eigen mat.EigenSym
	ok := eigen.Factorize(laplacianMatrix, true)
	if !ok {
		log.Fatal("Symmetric eigen decomposition failed")
	}
	var eigenVectors mat.Dense
	eigen.VectorsTo(&eigenVectors)

	_, eigenVectorsColLen := eigenVectors.Dims()
	fiedlerVector := eigenVectors.ColView(eigenVectorsColLen - 2)

	firstGroupSize := 0
	secondGroupSize := 0
	for i := 0; i < fiedlerVector.Len(); i++ {
		if fiedlerVector.AtVec(i) > 0 {
			firstGroupSize++
		} else if fiedlerVector.AtVec(i) < 0 {
			secondGroupSize++
		}
	}

	return firstGroupSize * secondGroupSize
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
