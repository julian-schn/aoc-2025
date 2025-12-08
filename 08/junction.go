package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type JunctionBox struct {
	X, Y, Z int
}

type JunctionDist struct {
	IdxA int
	IdxB int
	Dist float64
}

func findRoot(id int, parent []int) int {
	if parent[id] == id {
		return id
	}
	parent[id] = findRoot(parent[id], parent)
	return parent[id]
}

func main() {
	start := time.Now()

	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var boxes []JunctionBox

	// Convert bytes to string, then split
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		coords := strings.Split(line, ",")
		if len(coords) < 3 {
			continue
		}

		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])

		boxes = append(boxes, JunctionBox{X: x, Y: y, Z: z})
	}

	var dists []JunctionDist

	for i := 0; i < len(boxes); i++ {
		for j := i + 1; j < len(boxes); j++ {
			dx := float64(boxes[i].X - boxes[j].X)
			dy := float64(boxes[i].Y - boxes[j].Y)
			dz := float64(boxes[i].Z - boxes[j].Z)

			d := math.Sqrt(dx*dx + dy*dy + dz*dz)

			dists = append(dists, JunctionDist{IdxA: i, IdxB: j, Dist: d})
		}
	}

	sort.Slice(dists, func(i, j int) bool {
		return dists[i].Dist < dists[j].Dist
	})

	parent := make([]int, len(boxes))
	for i := range parent {
		parent[i] = i
	}

	connectionsMade := 0
	components := len(boxes)
	lastA, lastB := -1, -1

	for i := 0; i < len(dists); i++ {
		pair := dists[i]

		rootA := findRoot(pair.IdxA, parent)
		rootB := findRoot(pair.IdxB, parent)

		if rootA != rootB {
			parent[rootB] = rootA
			connectionsMade++
			components--
			lastA, lastB = pair.IdxA, pair.IdxB
			if components == 1 {
				break
			}
		}
	}
	fmt.Printf("Processed %d connections.\n", connectionsMade)

	counts := make(map[int]int)

	for i := 0; i < len(boxes); i++ {
		root := findRoot(i, parent)
		counts[root]++
	}

	sizes := make([]int, 0, len(counts))
	for _, v := range counts {
		sizes = append(sizes, v)
	}

	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	if len(sizes) >= 3 {
		fmt.Printf("Top 3 circuit sizes: %d, %d, %d\n", sizes[0], sizes[1], sizes[2])
		result := sizes[0] * sizes[1] * sizes[2]
		fmt.Println("Answer:", result)
	} else {
		fmt.Println("Not enough circuits formed to calculate top 3!")
		fmt.Println("Sizes found:", sizes)
	}

	if components == 1 && lastA >= 0 && lastB >= 0 {
		product := boxes[lastA].X * boxes[lastB].X
		fmt.Printf("Final connection between boxes at (%d, %d, %d) and (%d, %d, %d)\n",
			boxes[lastA].X, boxes[lastA].Y, boxes[lastA].Z,
			boxes[lastB].X, boxes[lastB].Y, boxes[lastB].Z)
		fmt.Printf("Product of their X coordinates: %d\n", product)
	} else {
		fmt.Println("Failed to connect all junction boxes into a single circuit.")
	}

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}
