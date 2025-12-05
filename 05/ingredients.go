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
mental model (p1)
1. parse id ranges by looping over input until empty line into 2d slice
2. parse ids be looping overe input after empty line until EOF
3. loop over fresh ingredient ids, compare to all range pairs
4. if id is within any range, increment fresh counter
*/

/*
mental model (p2)
1. parse fresh ids ranges into a slice.
2. sort the slice by start id.
3. iterate through sorted ranges to merge overlaps:
    - if overlap: merge (extend current.end).
    - if gap: store current range, start new one.
4. loop through the final merged ranges:
    - add (end - start + 1) to a running total
*/

type interval struct {
	low  int64
	high int64
}

func main() {
	start := time.Now()

	ranges, err := parseRanges("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	ids, err := parseIds("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	freshIds, err := findFreshIds(ids, ranges)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	sortedRanges, err := sortRanges(ranges)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	totalIds, err := findTotalIds(sortedRanges)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("P1: ", freshIds)
	fmt.Println("P2: ", totalIds)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseRanges(filename string) ([][]rune, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var inputs [][]rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break // reached separator between ranges and ids
		}
		inputs = append(inputs, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("no inputs parsed from %s", filename)
	}

	return inputs, nil
}

func parseIds(filename string) ([][]rune, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Skip recalled ranges
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			break
		}
	}

	var inputs [][]rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		inputs = append(inputs, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("no inputs parsed from %s", filename)
	}

	return inputs, nil
}

func findFreshIds(ids [][]rune, ranges [][]rune) (int, error) {
	parsedRanges := make([]interval, 0, len(ranges))
	for idx, r := range ranges {
		parts := strings.Split(string(r), "-")
		if len(parts) != 2 {
			return 0, fmt.Errorf("range %d is malformed: %q", idx, string(r))
		}

		low, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid low bound %q in range %d: %w", parts[0], idx, err)
		}

		high, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid high bound %q in range %d: %w", parts[1], idx, err)
		}

		if low > high {
			return 0, fmt.Errorf("range %d has low > high: %d > %d", idx, low, high)
		}

		parsedRanges = append(parsedRanges, interval{low: low, high: high})
	}

	fresh := 0

	for idx, idRunes := range ids {
		idVal, err := strconv.ParseInt(string(idRunes), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid ingredient id at index %d: %q: %w", idx, string(idRunes), err)
		}

		for _, r := range parsedRanges {
			if idVal >= r.low && idVal <= r.high {
				fresh++
				break
			}
		}
	}

	return fresh, nil
}

func sortRanges(ranges [][]rune) ([]interval, error) {
	parsed := make([]interval, 0, len(ranges))
	for idx, r := range ranges {
		parts := strings.Split(string(r), "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("range %d is malformed: %q", idx, string(r))
		}

		low, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid low bound %q in range %d: %w", parts[0], idx, err)
		}

		high, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid high bound %q in range %d: %w", parts[1], idx, err)
		}

		if low > high {
			return nil, fmt.Errorf("range %d has low > high: %d > %d", idx, low, high)
		}

		parsed = append(parsed, interval{low: low, high: high})
	}

	sort.Slice(parsed, func(i, j int) bool {
		if parsed[i].low == parsed[j].low {
			return parsed[i].high < parsed[j].high
		}
		return parsed[i].low < parsed[j].low
	})

	return parsed, nil
}

func findTotalIds(sorted []interval) (int64, error) {
	if len(sorted) == 0 {
		return 0, fmt.Errorf("no ranges supplied")
	}

	total := int64(0)
	current := sorted[0]

	for i := 1; i < len(sorted); i++ {
		next := sorted[i]
		if next.low <= current.high+1 { // overlap or touch
			if next.high > current.high {
				current.high = next.high
			}
			continue
		}

		total += current.high - current.low + 1
		current = next
	}

	total += current.high - current.low + 1

	return total, nil
}
