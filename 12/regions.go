package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
mental model:
1.  This is a tiling/packing problem. We need to fit a set of polyominoes (presents) into a rectangular region.
2.  Input has two parts:
    - Present definitions: Shape ID and visual grid.
    - Region definitions: Dimensions (WxH) and counts of each present ID required.
3.  We need to answer: How many regions can validly fit all their required presents?
4.  Algorithm: Backtracking (DFS).
    - Board is WxH boolean grid (false=empty, true=filled).
    - Find first empty command (r, c) in reading order.
    - If no empty cells:
      - Check if all required presents are used. If so, Success.
    - If empty cells exist but no presents left: Failure (backtrack).
    - Determine which present variants can cover (r, c).
      - A "variant" is a specific rotation/flip of a present shape.
      - "Covering (r, c)" means placing the variant such that its *first* solid cell lands exactly on (r, c).
      - This prevents trying the same piece in the same spot multiple times (canonical placement).
    - For each shape ID with count > 0:
      - For each unique variant of this shape:
        - Check if placing variant at (r, c) - (variant's first cell offset) fits on board and doesn't overlap.
        - If fits: Place it, decrement count, Recurse.
        - If recursion returns true: Return true (we only need one valid packing).
        - Backtrack: Remove piece, increment count.
    - Return false if no option works.

5.  Optimization:
    - Pre-calculate all unique variants for each shape ID (rotations 0,90,180,270 + flips).
    - Store variants as a list of relative coordinates (offsets from top-left).
    - "First solid cell" logic crucial for pruning.
*/

type Point struct {
	r, c int
}

type Shape struct {
	id       int
	cells    []Point
	variants [][]Point
}

type RegionRequest struct {
	width, height int
	counts        map[int]int
}

func main() {
	start := time.Now()

	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	shapes, requests, err := parseInput(filename)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	// Pre-process shapes to generate variants
	processedShapes := make(map[int][][]Point)
	for id, s := range shapes {
		processedShapes[id] = generateVariants(s)
	}

	solvableCount := 0
	for _, req := range requests {
		if solve(req, processedShapes) {
			solvableCount++
		}
	}

	fmt.Println("Part 1: Solvable regions:", solvableCount)
	fmt.Println("Runtime:", time.Since(start))
}

func solve(req RegionRequest, shapes map[int][][]Point) bool {
	// 1. Basic area check pruning
	totalArea := 0
	for id, count := range req.counts {
		if count > 0 {
			if len(shapes[id]) == 0 {
				panic(fmt.Sprintf("Shape %d has no variants?", id))
			}
			// Area of a shape is number of cells in any variant
			area := len(shapes[id][0])
			totalArea += area * count
		}
	}
	if totalArea > req.width*req.height {
		return false
	}

	board := make([][]bool, req.height)
	for i := range board {
		board[i] = make([]bool, req.width)
	}

	var todo []int
	for id, count := range req.counts {
		for k := 0; k < count; k++ {
			todo = append(todo, id)
		}
	}

	sort.Slice(todo, func(i, j int) bool {
		// Larger area first
		areaI := len(shapes[todo[i]][0])
		areaJ := len(shapes[todo[j]][0])
		return areaI > areaJ
	})

	return placePresents(board, todo, shapes, req.width, req.height)
}

func placePresents(board [][]bool, todo []int, shapes map[int][][]Point, W, H int) bool {
	if len(todo) == 0 {
		return true
	}

	currentID := todo[0]
	remaining := todo[1:]

	variants := shapes[currentID]

	for _, variant := range variants {
		maxR, maxC := 0, 0
		for _, p := range variant {
			if p.r > maxR {
				maxR = p.r
			}
			if p.c > maxC {
				maxC = p.c
			}
		}

		for r := 0; r < H-maxR; r++ {
			for c := 0; c < W-maxC; c++ {
				// Check if fits
				fits := true
				for _, p := range variant {
					if board[r+p.r][c+p.c] {
						fits = false
						break
					}
				}

				if fits {
					// Place
					for _, p := range variant {
						board[r+p.r][c+p.c] = true
					}

					if placePresents(board, remaining, shapes, W, H) {
						return true
					}

					// Remove
					for _, p := range variant {
						board[r+p.r][c+p.c] = false
					}
				}
			}
		}
	}

	return false
}

func parseInput(filename string) (map[int][]Point, []RegionRequest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	shapes := make(map[int][]Point)
	var requests []RegionRequest

	var currentID int = -1
	var currentRows []string

	processShape := func() {
		if currentID != -1 && len(currentRows) > 0 {
			var points []Point
			for r, line := range currentRows {
				for c, char := range line {
					if char == '#' {
						points = append(points, Point{r, c})
					}
				}
			}
			// Normalize points to top-left (0,0)
			minR, minC := 1000, 1000
			for _, p := range points {
				if p.r < minR {
					minR = p.r
				}
				if p.c < minC {
					minC = p.c
				}
			}
			for i := range points {
				points[i].r -= minR
				points[i].c -= minC
			}
			shapes[currentID] = points
		}
		currentID = -1
		currentRows = nil
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			processShape()
			continue
		}

		// Check if region line "12x5: ..."
		if strings.Contains(line, "x") && strings.Contains(line, ":") {
			processShape() // Finish any pending shape

			parts := strings.Split(line, ":")
			dimParts := strings.Split(parts[0], "x")
			w, _ := strconv.Atoi(dimParts[0])
			h, _ := strconv.Atoi(dimParts[1])

			countStr := strings.Fields(parts[1])
			counts := make(map[int]int)
			for i, s := range countStr {
				n, _ := strconv.Atoi(s)
				counts[i] = n
			}
			requests = append(requests, RegionRequest{w, h, counts})
			continue
		}

		// Check if shape header "0:"
		if strings.HasSuffix(line, ":") {
			processShape()
			val, _ := strconv.Atoi(strings.TrimSuffix(line, ":"))
			currentID = val
		} else {
			// It's a shape row
			currentRows = append(currentRows, line)
		}
	}
	processShape() // Last one

	return shapes, requests, scanner.Err()
}

func generateVariants(base []Point) [][]Point {
	unique := make(map[string][]Point)

	add := func(pts []Point) {
		// Normalize
		minR, minC := 1000, 1000
		for _, p := range pts {
			if p.r < minR {
				minR = p.r
			}
			if p.c < minC {
				minC = p.c
			}
		}

		normalized := make([]Point, len(pts))
		for i, p := range pts {
			normalized[i] = Point{p.r - minR, p.c - minC}
		}

		// Sort for consistent key
		sort.Slice(normalized, func(i, j int) bool {
			if normalized[i].r == normalized[j].r {
				return normalized[i].c < normalized[j].c
			}
			return normalized[i].r < normalized[j].r
		})

		keyBuilder := strings.Builder{}
		for _, p := range normalized {
			keyBuilder.WriteString(fmt.Sprintf("%d,%d|", p.r, p.c))
		}
		key := keyBuilder.String()
		unique[key] = normalized
	}

	current := base
	for i := 0; i < 4; i++ {
		add(current)       // Normal
		add(flip(current)) // Flipped
		current = rotate(current)
	}

	result := make([][]Point, 0, len(unique))
	for _, v := range unique {
		result = append(result, v)
	}
	return result
}

func rotate(pts []Point) []Point {
	// Rotate 90 deg clockwise: (r, c) -> (c, -r)
	newPts := make([]Point, len(pts))
	for i, p := range pts {
		newPts[i] = Point{p.c, -p.r}
	}
	return newPts
}

func flip(pts []Point) []Point {
	// Flip over x-axis (rows): (r, c) -> (-r, c)
	newPts := make([]Point, len(pts))
	for i, p := range pts {
		newPts[i] = Point{-p.r, p.c}
	}
	return newPts
}
