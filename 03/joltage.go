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
1. Parse input into slice, each line = one bank
2. Find highest number made up of combination of digits in bank (cannot rearrange)
3. Create new slice using only unique numbers from original slice
4. Use highest number in new slice to find index in original slice
5. Use that index and check if it is the last index of original slice
6. If yes: Go with next highest number in new slice, repeat from step 4
7. If no: Iterate through indexes after found index in original slice, find highest number
8. Add combination to sum
*/

/*
Mental model (p2):
1. Parse input into slice, each line = one bank
2. "Window approach", create a loop that runs 12 times
3. Find upper bound so we dont run out of numbers toward the end (window size) with current string length
4. Find best digit within window and append to current string
5. Increment window start index, decrement needed numbers
6. Repeat
*/

func main() {
	start := time.Now()

	banks, err := parseBanks("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	occuringDigits, err := findOccuringDigits(banks)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	sum, err := findMaxJoltageSum(banks, occuringDigits)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Part 1: %d\n", sum)

	sumP2, err := findMaxJoltageSumP2(banks)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Part 2: %d\n", sumP2)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseBanks(filename string) ([][]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var banks [][]int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var bank []int
		for _, r := range line {
			digit, err := strconv.Atoi(string(r))
			if err != nil {
				return nil, fmt.Errorf("invalid digit %q: %w", r, err)
			}
			bank = append(bank, digit)
		}
		banks = append(banks, bank)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

func findOccuringDigits(banks [][]int) (map[int]bool, error) {
	digits := make(map[int]bool)
	for _, bank := range banks {
		for _, digit := range bank {
			digits[digit] = true
		}
	}
	return digits, nil
}

func findMaxJoltageSum(banks [][]int, occuringDigits map[int]bool) (int, error) {
	totalSum := 0

	for _, bank := range banks {
		maxJoltage := 0
		// We need to pick two batteries (digits) at indices i and j where i < j
		// to form a number d1*10 + d2.
		// We want to maximize this number.
		for i := 0; i < len(bank); i++ {
			for j := i + 1; j < len(bank); j++ {
				d1 := bank[i]
				d2 := bank[j]
				joltage := d1*10 + d2
				if joltage > maxJoltage {
					maxJoltage = joltage
				}
			}
		}
		totalSum += maxJoltage
	}

	return totalSum, nil
}

func findMaxJoltageSumP2(banks [][]int) (int, error) {
	totalSum := 0

	for _, bank := range banks {

		currentIdx := 0
		var resultDigits []int
		needed := 12

		for needed > 0 {

			endSearch := len(bank) - needed
			if currentIdx > endSearch {
				return 0, fmt.Errorf("not enough digits to form 12-digit number")
			}

			bestDigit := -1
			bestIdx := -1

			for i := currentIdx; i <= endSearch; i++ {
				if bank[i] > bestDigit {
					bestDigit = bank[i]
					bestIdx = i
				}
				if bank[i] == 9 {
					break
				}
			}

			resultDigits = append(resultDigits, bestDigit)
			currentIdx = bestIdx + 1
			needed--
		}

		valStr := ""
		for _, d := range resultDigits {
			valStr += strconv.Itoa(d)
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return 0, fmt.Errorf("failed to convert digits to int: %w", err)
		}
		totalSum += val
	}

	return totalSum, nil
}
