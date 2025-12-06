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
mental model (p1)
1. parse input as 2d slice, one dimnesion that contains all numbers (subslice?), other that contains operator
1.1. check if input is integer, save if yes, if no, save to operator slive
2. iterate over slice and compleete operations
3. sum up results
*/

func main() {
	start := time.Now()

	numbers, operators, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	sum, err := computeSum(numbers, operators)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Sum P1: ", sum)

	sumP2, err := solvePart2("input.txt")
	if err != nil {
		fmt.Println("error p2:", err)
		return
	}
	fmt.Println("Sum P2: ", sumP2)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseInput(filename string) ([][]int, []string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var numbers [][]int
	var operators []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		// Check if this line is an operator line (only "+" or "*")
		isOp := true
		for _, p := range parts {
			if p != "+" && p != "*" {
				isOp = false
				break
			}
		}

		if isOp {
			operators = parts
			continue
		}

		row := make([]int, len(parts))
		for i, p := range parts {
			val, err := strconv.Atoi(p)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid number %q: %w", p, err)
			}
			row[i] = val
		}
		numbers = append(numbers, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	if len(numbers) == 0 || len(operators) == 0 {
		return nil, nil, fmt.Errorf("no inputs parsed from %s", filename)
	}

	// Validate consistent row length and operator count matches column count
	expectedCols := len(numbers[0])
	for i, row := range numbers {
		if len(row) != expectedCols {
			return nil, nil, fmt.Errorf("row %d has inconsistent length: %d (expected %d)", i, len(row), expectedCols)
		}
	}

	if len(operators) != expectedCols {
		return nil, nil, fmt.Errorf("operator count %d does not match expected %d", len(operators), expectedCols)
	}

	return numbers, operators, nil
}

func computeSum(numbers [][]int, operators []string) (int, error) {
	if len(numbers) == 0 {
		return 0, fmt.Errorf("no numbers provided")
	}
	if len(operators) == 0 {
		return 0, fmt.Errorf("no operators provided")
	}

	total := 0

	cols := len(operators)
	for i := 0; i < cols; i++ {
		var res int
		if operators[i] == "+" {
			// sum column
			for _, row := range numbers {
				if len(row) <= i {
					return 0, fmt.Errorf("row too short at column %d", i)
				}
				res += row[i]
			}
		} else if operators[i] == "*" {
			res = 1
			for _, row := range numbers {
				if len(row) <= i {
					return 0, fmt.Errorf("row too short at column %d", i)
				}
				res *= row[i]
			}
		} else {
			return 0, fmt.Errorf("invalid operator %q at column %d", operators[i], i)
		}
		total += res
	}

	return total, nil
}

func solvePart2(filename string) (int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if len(lines) == 0 {
		return 0, fmt.Errorf("empty file")
	}

	// Pad lines to max length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	for i, line := range lines {
		if len(line) < maxLen {
			lines[i] = line + strings.Repeat(" ", maxLen-len(line))
		}
	}

	rows := len(lines)
	cols := maxLen
	var total int64

	// Helper to check if a column is empty (only spaces)
	isEmptyCol := func(c int) bool {
		for r := 0; r < rows; r++ {
			if lines[r][c] != ' ' {
				return false
			}
		}
		return true
	}

	// Process a block of columns
	processBlock := func(blockCols []int) {
		if len(blockCols) == 0 {
			return
		}

		var operands []int64
		var operator string

		// blockCols are ordered Right to Left
		for _, c := range blockCols {
			// Parse number from rows 0 to rows-2
			numStr := ""
			for r := 0; r < rows-1; r++ {
				char := lines[r][c]
				if char != ' ' {
					numStr += string(char)
				}
			}
			if numStr != "" {
				val, err := strconv.ParseInt(numStr, 10, 64)
				if err == nil {
					operands = append(operands, val)
				}
			}

			// Check operator in last row
			opChar := lines[rows-1][c]
			if opChar != ' ' {
				operator = string(opChar)
			}
		}

		if len(operands) == 0 {
			return
		}

		res := operands[0]
		for i := 1; i < len(operands); i++ {
			if operator == "+" {
				res += operands[i]
			} else if operator == "*" {
				res *= operands[i]
			}
		}
		total += res
	}

	var currentBlock []int
	// Scan columns right to left
	for c := cols - 1; c >= 0; c-- {
		if isEmptyCol(c) {
			processBlock(currentBlock)
			currentBlock = nil
		} else {
			currentBlock = append(currentBlock, c)
		}
	}
	processBlock(currentBlock)

	return total, nil
}
