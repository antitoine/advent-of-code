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

func (uf *UnionFind) GetCircuitSizes() []int {
	sizes := make(map[int]int)
	for i := range uf.parent {
		root := uf.Find(i)
		sizes[root] = uf.size[root]
	}
	result := make([]int, 0, len(sizes))
	for _, s := range sizes {
		result = append(result, s)
	}
	return result
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

	// Determine number of connections based on input size
	numConnections := 10
	if n > 100 {
		numConnections = 1000
	}

	// Connect closest pairs using Union-Find
	// We process numConnections pairs (even if some are already in the same circuit)
	uf := NewUnionFind(n)
	for i := 0; i < numConnections && i < len(pairs); i++ {
		uf.Union(pairs[i].i, pairs[i].j)
	}

	// Get circuit sizes and find 3 largest
	sizes := uf.GetCircuitSizes()
	slices.SortFunc(sizes, func(a, b int) int { return b - a }) // Descending

	result := int64(1)
	for i := 0; i < 3 && i < len(sizes); i++ {
		result *= int64(sizes[i])
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
