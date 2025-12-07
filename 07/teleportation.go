package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

/*
mental model (p1)
1. parse: read the input file into a 2d slice of runes (rows)
2. locate start: find the S to get the starting column index of the single beam
3. init: splitCount = 0; currentBeams slice/map for current row, set S index = true
4. simulation loop (rows from below S):
   - make nextBeams empty for next row
   - scan current beams; if beam at i:
       * '.' => falls straight down, mark i in nextBeams
       * '^' => splitCount++; if i > 0 mark i-1, if i+1 < width mark i+1
   - advance: currentBeams = nextBeams
5. result: return splitCount
*/

func main() {
	start := time.Now()

	grid, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	splits, err := countSplits(grid)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Split count:", splits)

	timelines, err := countTimelines(grid)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Timeline count:", timelines)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseInput(path string) ([][]rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var rows [][]rune
	expectedWidth := -1

	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)

		if expectedWidth == -1 {
			expectedWidth = len(row)
		} else if len(row) != expectedWidth {
			return nil, fmt.Errorf("inconsistent row width: got %d expected %d", len(row), expectedWidth)
		}

		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty grid")
	}

	return rows, nil
}

func countSplits(grid [][]rune) (int, error) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0, fmt.Errorf("grid is empty")
	}

	rows := len(grid)
	cols := len(grid[0])

	startRow, startCol := -1, -1
	for r, row := range grid {
		if len(row) != cols {
			return 0, fmt.Errorf("row %d has inconsistent width", r)
		}
		for c, ch := range row {
			if ch == 'S' {
				startRow, startCol = r, c
				break
			}
		}
		if startRow != -1 {
			break
		}
	}

	if startRow == -1 {
		return 0, fmt.Errorf("start position 'S' not found")
	}

	current := make([]bool, cols)
	current[startCol] = true
	splitCount := 0

	for r := startRow + 1; r < rows; r++ {
		next := make([]bool, cols)
		for c, active := range current {
			if !active {
				continue
			}

			switch grid[r][c] {
			case '.':
				next[c] = true
			case '^':
				if c-1 >= 0 {
					next[c-1] = true
				}
				if c+1 < cols {
					next[c+1] = true
				}
				splitCount++
			}
		}
		current = next
	}

	return splitCount, nil
}

/*
mental model (p2)
1. parse: same as p1
2. locate start: same as p1
3. use memoized recursion to find all unique end positions reachable from each position
4. for each position (row, col):
   - if at or past bottom row, return set containing this column (particle has exited grid)
   - if memoized, return cached result (copy to avoid sharing reference)
   - if '.', recursively get end positions from (row+1, col)
   - if '^', union end positions from (row+1, col-1) and (row+1, col+1)
5. result: return size of set of unique end positions from start
*/

func countTimelines(grid [][]rune) (int, error) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0, fmt.Errorf("grid is empty")
	}

	rows := len(grid)
	cols := len(grid[0])

	startRow, startCol := -1, -1
	for r, row := range grid {
		if len(row) != cols {
			return 0, fmt.Errorf("row %d has inconsistent width", r)
		}
		for c, ch := range row {
			if ch == 'S' {
				startRow, startCol = r, c
				break
			}
		}
		if startRow != -1 {
			break
		}
	}

	if startRow == -1 {
		return 0, fmt.Errorf("start position 'S' not found")
	}

	// memo: key is [row, col], value is count of paths from that point to bottom
	memo := make(map[[2]int]int)
	count := traverse(startRow+1, startCol, grid, rows, cols, memo)
	return count, nil
}

func traverse(row, col int, grid [][]rune, rows, cols int, memo map[[2]int]int) int {
	// If we've reached past the last row, we found 1 valid timeline end
	if row >= rows {
		return 1
	}

	key := [2]int{row, col}
	if cached, exists := memo[key]; exists {
		return cached
	}

	count := 0
	switch grid[row][col] {
	case '.':
		// Continue straight down
		count = traverse(row+1, col, grid, rows, cols, memo)
	case '^':
		// Split left and right
		if col-1 >= 0 {
			count += traverse(row+1, col-1, grid, rows, cols, memo)
		}
		if col+1 < cols {
			count += traverse(row+1, col+1, grid, rows, cols, memo)
		}
	default:
		// Assumes pass-through for any other characters (like '|') used in visualization but acting as path
		count = traverse(row+1, col, grid, rows, cols, memo)
	}
	memo[key] = count
	return count
}
