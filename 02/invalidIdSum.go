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
Mental model (p1):
1. Parse each range into 2d array (note edge cases) parseRanges()
2. Iterate every number, cast to string and compare first half and second half
3. If match, convert to int and add to sum (invalid ID)
*/

/*
Mental Model (p2):
1. Same as above
2. Find patterns? -> Find pattern identification algorithm
3. Sum invalid IDs
*/

func main() {
	start := time.Now()

	ranges, err := parseRanges("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	sumP1, err := findInvalidSumP1(ranges)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Part 1: %d\n", sumP1)

	sumP2, err := findInvalidSumP2(ranges)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Part 2: %d\n", sumP2)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseRanges(filename string) ([][2]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var result [][2]int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Split by comma
		parts := strings.Split(line, ",")

		for _, part := range parts {
			bounds := strings.Split(strings.TrimSpace(part), "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range %q", part)
			}

			low, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid low bound %q: %w", bounds[0], err)
			}

			high, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid high bound %q: %w", bounds[1], err)
			}

			result = append(result, [2]int{low, high})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func findInvalidSumP1(ranges [][2]int) (int, error) {
	sum := 0

	for _, r := range ranges {
		low := r[0]
		high := r[1]

		for i := low; i <= high; i++ {
			s := strconv.Itoa(i)
			length := len(s)

			if length%2 != 0 {
				continue // skip odd-length ids
			}

			mid := length / 2
			firstHalf := s[:mid]
			secondHalf := s[mid:]

			if firstHalf == secondHalf {
				num, err := strconv.Atoi(s)
				if err != nil {
					return 0, fmt.Errorf("invalid number %q: %w", s, err)
				}
				sum += num
			}
		}
	}

	return sum, nil
}

func findInvalidSumP2(ranges [][2]int) (int, error) {
	sum := 0

	for _, r := range ranges {
		low := r[0]
		high := r[1]

		for i := low; i <= high; i++ {
			s := strconv.Itoa(i)
			if isInvalidP2(s) {
				sum += i
			}
		}
	}

	return sum, nil
}

func isInvalidP2(s string) bool {
	n := len(s)
	// Try all possible pattern lengths L
	// The pattern must repeat at least twice, so L <= n/2
	for L := 1; L <= n/2; L++ {
		if n%L == 0 {
			pattern := s[:L]
			repeats := n / L
			if strings.Repeat(pattern, repeats) == s {
				return true
			}
		}
	}
	return false
}
