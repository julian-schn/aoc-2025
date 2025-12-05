package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

/*
mental model (p1)
1. parse input into a 2d slice
2. loop through rows
3. check if paper (@), NOT empty (.)
4. if yes: check 8 neighbors, edge case: out of bounds
5. if there is fewer than 4 rolls of papers in neighbors, increase counter
*/
/*
mental model (p2)
1. start with same 2d slice as p1
2. precompute neighbor counts for every roll
3. queue all rolls with fewer than 4 neighbors (these can be removed)
4. pop queue, remove roll, decrease neighbor counts of its adjacent rolls
5. if a neighbor drops below 4, enqueue it
6. repeat until queue empty, return how many were removed in total
*/

func main() {
	start := time.Now()

	shelves, err := parseShelves("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	rolls, err := findValidRolls(shelves)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Part 1: %d\n", rolls)

	removed, err := findRemovableRolls(shelves)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Part 2: %d\n", removed)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseShelves(filename string) ([][]rune, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var shelves [][]rune
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		shelves = append(shelves, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(shelves) == 0 {
		return nil, fmt.Errorf("no shelves parsed from %s", filename)
	}

	return shelves, nil
}

func findValidRolls(shelves [][]rune) (int, error) {
	if len(shelves) == 0 {
		return 0, fmt.Errorf("no shelves supplied")
	}

	rows := len(shelves)
	cols := len(shelves[0])
	for i, row := range shelves {
		if len(row) != cols {
			return 0, fmt.Errorf("row %d has inconsistent length", i)
		}
	}

	count := 0
	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1} /*self*/, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if shelves[r][c] != '@' {
				continue // empty slot
			}

			neighbors := 0
			for _, d := range dirs {
				nr := r + d[0]
				nc := c + d[1]
				if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
					continue
				}
				if shelves[nr][nc] == '@' {
					neighbors++
				}
			}

			if neighbors < 4 {
				count++
			}
		}
	}

	return count, nil
}

func findRemovableRolls(shelves [][]rune) (int, error) {
	if len(shelves) == 0 {
		return 0, fmt.Errorf("no shelves supplied")
	}

	rows := len(shelves)
	cols := len(shelves[0])
	for i, row := range shelves {
		if len(row) != cols {
			return 0, fmt.Errorf("row %d has inconsistent length", i)
		}
	}

	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	neighborCount := make([][]int, rows)
	removed := make([][]bool, rows)
	for r := 0; r < rows; r++ {
		neighborCount[r] = make([]int, cols)
		removed[r] = make([]bool, cols)
		for c := 0; c < cols; c++ {
			if shelves[r][c] != '@' {
				continue
			}
			for _, d := range dirs {
				nr := r + d[0]
				nc := c + d[1]
				if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
					continue
				}
				if shelves[nr][nc] == '@' {
					neighborCount[r][c]++
				}
			}
		}
	}

	type pos struct{ r, c int }
	queue := make([]pos, 0)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if shelves[r][c] == '@' && neighborCount[r][c] < 4 {
				queue = append(queue, pos{r, c})
			}
		}
	}

	removedCount := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		r := current.r
		c := current.c

		if removed[r][c] || shelves[r][c] != '@' {
			continue
		}

		removed[r][c] = true
		removedCount++

		for _, d := range dirs {
			nr := r + d[0]
			nc := c + d[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				continue
			}
			if shelves[nr][nc] != '@' || removed[nr][nc] {
				continue
			}
			neighborCount[nr][nc]--
			if neighborCount[nr][nc] == 3 { // just dropped below 4
				queue = append(queue, pos{nr, nc})
			}
		}
	}

	return removedCount, nil
}
