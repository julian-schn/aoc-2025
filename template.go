package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

/*
mental model (p1)
1. this is a template
2. test
*/

func main() {
	start := time.Now()

	slice, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseInput(filename string) ([][]rune, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var inputs [][]rune
	for scanner.Scan() {
		line := scanner.Text()
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
