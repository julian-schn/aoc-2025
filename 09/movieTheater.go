package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
mental model:
1. Parse the input which consists of coordinates "x,y".
2. Store these coordinates as a list of Points.
3. The goal is to find the largest rectangle where two of the red tiles (the input points) are opposite corners.
4. Iterate through all unique pairs of points.
5. For each pair, calculate the area of the rectangle formed by these two points as opposite corners.
   The width is abs(x1 - x2) + 1.
   The height is abs(y1 - y2) + 1.
   Area = width * height.
6. Track the maximum area found.
*/

type Point struct {
	X, Y int
}

func main() {
	start := time.Now()

	points, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	// Revert to input.txt for final run or use logic to switch
	// For now, let's keep input.txt as the default

	maxArea := solvePart1(points)
	fmt.Println("Max Area Part 1:", maxArea)

	maxArea2 := solvePart2(points)
	fmt.Println("Max Area Part 2:", maxArea2)

	elapsed := time.Since(start)

	fmt.Println("Runtime:", elapsed)
}

func parseInput(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid x coordinate: %v", err)
		}
		y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid y coordinate: %v", err)
		}
		points = append(points, Point{X: x, Y: y})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return points, nil
}

func solvePart1(points []Point) int {
	maxArea := 0
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1 := points[i]
			p2 := points[j]

			width := abs(p1.X-p2.X) + 1
			height := abs(p1.Y-p2.Y) + 1
			area := width * height

			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}
func solvePart2(points []Point) int {
	maxArea := 0
	n := len(points)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := points[i]
			p2 := points[j]

			// Form rectangle
			minX := min(p1.X, p2.X)
			maxX := max(p1.X, p2.X)
			minY := min(p1.Y, p2.Y)
			maxY := max(p1.Y, p2.Y)

			width := maxX - minX + 1
			height := maxY - minY + 1
			area := width * height

			// Optimization: only check if this area is better than current max
			if area <= maxArea {
				continue
			}

			if isRectInsidePolygon(minX, maxX, minY, maxY, points) {
				maxArea = area
			}
		}
	}
	return maxArea
}

func isRectInsidePolygon(minX, maxX, minY, maxY int, polygon []Point) bool {
	// 1. Check if any polygon edge intersects the interior or crosses the rectangle boundaries in a prohibiting way
	// A strictly simpler check: No polygon edge should intersect the OPEN rectangle (minX, maxX) x (minY, maxY).
	// If a polygon edge passes through the interior, the rectangle is split (part in, part out).
	// Note: The problem might allow edges to be collinear with the rectangle boundary (shared boundary/green tiles).
	// So we specifically check for intersection with the open intervals.

	n := len(polygon)
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]

		if p1.X == p2.X {
			// Vertical edge at X = p1.X
			edgeYMin := min(p1.Y, p2.Y)
			edgeYMax := max(p1.Y, p2.Y)

			// Check if X is strictly inside (minX, maxX)
			if p1.X > minX && p1.X < maxX {
				// Check if Y intervals overlap (strictly)
				// Overlap of (minY, maxY) and [edgeYMin, edgeYMax]
				// if max(minY, edgeYMin) < min(maxY, edgeYMax)
				if max(minY, edgeYMin) < min(maxY, edgeYMax) {
					return false
				}
			}
		} else {
			// Horizontal edge at Y = p1.Y
			edgeXMin := min(p1.X, p2.X)
			edgeXMax := max(p1.X, p2.X)

			// Check if Y is strictly inside (minY, maxY)
			if p1.Y > minY && p1.Y < maxY {
				// Check if X intervals overlap (strictly)
				if max(minX, edgeXMin) < min(maxX, edgeXMax) {
					return false
				}
			}
		}
	}

	// 2. Check if the rectangle center is inside the polygon
	// Since no edges intersect the interior, checking one point is sufficient.
	// Use Ray Casting.
	cx := float64(minX+maxX) / 2.0
	cy := float64(minY+maxY) / 2.0

	inside := false
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]

		x1, y1 := float64(p1.X), float64(p1.Y)
		x2, y2 := float64(p2.X), float64(p2.Y)

		// Check if ray from (cx, cy) to (infinity, cy) crosses edge (p1, p2)
		// Ray is horizontal to the right.
		// Edge must cross the horizontal line y = cy.

		// Condition: one endpoint above cy, one below
		if (y1 > cy) != (y2 > cy) {
			// Compute X coordinate of intersection
			// x = x1 + (cy - y1) * (x2 - x1) / (y2 - y1)
			intersectX := x1 + (cy-y1)*(x2-x1)/(y2-y1)
			if cx < intersectX {
				inside = !inside
			}
		}
	}

	return inside
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
