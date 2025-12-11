package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

/*
mental model:
1. Represent the devices as a directed graph where nodes are device names and edges are connections.
2. The input format "aaa: bbb ccc" defines edges from 'aaa' to 'bbb' and 'aaa' to 'ccc'.
3. The goal is to finding the number of distinct paths from "you" to "out".
4. Since we need to find *all* paths, and cycles are not mentioned (but "Data can't flow backwards" implies it's a DAG or at least we treat it as directed), Simple DFS is sufficient.
5. However, there might be many overlapping subproblems (many paths converging and then splitting), so Memoization (caching results for each node) is critical to avoid exponential blowup.
   - map[string]int where key is node name, value is number of paths from that node to "out".

Part 2 Update:
1. We need paths from 'svr' to 'out' that pass through 'dac' and 'fft'.
2. Method: Break down into segments. Since data can't flow backwards, either 'dac' comes before 'fft' or vice-versa.
   - Path 1: svr -> ... -> dac -> ... -> fft -> ... -> out
     Count = paths(svr, dac) * paths(dac, fft) * paths(fft, out)
   - Path 2: svr -> ... -> fft -> ... -> dac -> ... -> out
     Count = paths(svr, fft) * paths(fft, dac) * paths(dac, out)
   - Total = Path 1 + Path 2
*/

func main() {
	start := time.Now()

	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	graph, err := parseInput(filename)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	part1Result := solvePart1(graph)
	fmt.Println("Part 1: distinct paths from you to out:", part1Result)

	part2Result := solvePart2(graph)
	fmt.Println("Part 2: paths from svr to out via dac and fft:", part2Result)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseInput(filename string) (map[string][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Format: "aaa: bbb ccc"
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		source := strings.TrimSpace(parts[0])
		targetsStr := strings.TrimSpace(parts[1])
		targets := strings.Fields(targetsStr) // splits by whitespace

		graph[source] = append(graph[source], targets...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return graph, nil
}

func solvePart1(graph map[string][]string) int {
	memo := make(map[string]int)
	return countPaths("you", "out", graph, memo)
}

func solvePart2(graph map[string][]string) int {
	// Path A: svr -> dac -> fft -> out
	p1 := countPaths("svr", "dac", graph, make(map[string]int))
	p2 := countPaths("dac", "fft", graph, make(map[string]int))
	p3 := countPaths("fft", "out", graph, make(map[string]int))

	pathA := p1 * p2 * p3

	// Path B: svr -> fft -> dac -> out
	p4 := countPaths("svr", "fft", graph, make(map[string]int))
	p5 := countPaths("fft", "dac", graph, make(map[string]int))
	p6 := countPaths("dac", "out", graph, make(map[string]int))

	pathB := p4 * p5 * p6

	return pathA + pathB
}

func countPaths(current, target string, graph map[string][]string, memo map[string]int) int {
	if current == target {
		return 1
	}

	if count, ok := memo[current]; ok {
		return count
	}

	totalPaths := 0
	neighbors := graph[current]
	for _, next := range neighbors {
		totalPaths += countPaths(next, target, graph, memo)
	}

	memo[current] = totalPaths
	return totalPaths
}
