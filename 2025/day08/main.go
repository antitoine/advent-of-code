package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x, y, z int
}

type Pair struct {
	i, j     int
	distance float64
}

type UnionFind struct {
	parent []int
	rank   []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent: parent, rank: rank, size: size}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX == rootY {
		return false // Already in same circuit
	}
	// Union by rank
	if uf.rank[rootX] < uf.rank[rootY] {
		rootX, rootY = rootY, rootX
	}
	uf.parent[rootY] = rootX
	uf.size[rootX] += uf.size[rootY]
	if uf.rank[rootX] == uf.rank[rootY] {
		uf.rank[rootX]++
	}
	return true
}

func (uf *UnionFind) IsSingleCircuit() bool {
	// Check if all nodes are in the same circuit
	// We can check if the root's size equals the total number of nodes
	for i := range uf.parent {
		if uf.Find(i) == i && uf.size[i] == len(uf.parent) {
			return true
		}
	}
	return false
}

func distance(a, b Point) float64 {
	dx := float64(a.x - b.x)
	dy := float64(a.y - b.y)
	dz := float64(a.z - b.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func getResult(input io.Reader) int64 {
	scanner := bufio.NewScanner(input)
	var points []Point

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		points = append(points, Point{x, y, z})
	}

	n := len(points)

	// Generate all pairs with distances
	var pairs []Pair
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pairs = append(pairs, Pair{i, j, distance(points[i], points[j])})
		}
	}

	// Sort pairs by distance
	slices.SortFunc(pairs, func(a, b Pair) int {
		if a.distance < b.distance {
			return -1
		}
		if a.distance > b.distance {
			return 1
		}
		return 0
	})

	// Connect closest pairs until all boxes are in a single circuit
	uf := NewUnionFind(n)
	var lastPair Pair

	for _, pair := range pairs {
		// Try to union - if successful, this merged two circuits
		if uf.Union(pair.i, pair.j) {
			lastPair = pair
			// Check if all boxes are now in a single circuit
			if uf.IsSingleCircuit() {
				break
			}
		}
	}

	// Return the product of X coordinates of the last two boxes connected
	return int64(points[lastPair.i].x) * int64(points[lastPair.j].x)
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
