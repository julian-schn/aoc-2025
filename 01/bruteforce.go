package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	result, err := p2("input.txt")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(result)
}

func p1(filename string) (int, error) {
	position := 50 // starting position
	hits := 0      // how many times we land on 0

	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		first := line[0] // 'L' or 'R'
		rest := line[1:] // distance as string

		steps, err := strconv.Atoi(rest)
		if err != nil {
			return 0, fmt.Errorf("invalid distance %q in line %q: %w", rest, line, err)
		}

		rotation := 0
		switch first {
		case 'L':
			rotation = -steps // L → towards lower numbers
		case 'R':
			rotation = steps // R → towards higher numbers
		default:
			return 0, fmt.Errorf("invalid direction %q in line %q", first, line)
		}

		position += rotation

		// wrap to 0–99
		position = ((position % 100) + 100) % 100

		if position == 0 {
			hits++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return hits, nil
}

func p2(filename string) (int, error) {
	position := 50 // starting position
	hits := 0      // how many times we cross 0

	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		first := line[0] // 'L' or 'R'
		rest := line[1:] // distance as string

		steps, err := strconv.Atoi(rest)
		if err != nil {
			return 0, fmt.Errorf("invalid distance %q in line %q: %w", rest, line, err)
		}

		rotation := 0
		switch first {
		case 'L':
			rotation = -steps // L → towards lower numbers
		case 'R':
			rotation = steps // R → towards higher numbers
		default:
			return 0, fmt.Errorf("invalid direction %q in line %q", first, line)
		}

		if rotation > 0 {
			// move right (towards higher numbers)
			for i := 0; i < rotation; i++ {
				position++
				if position == 100 {
					position = 0
				}
				if position == 0 {
					hits++
				}
			}
		} else if rotation < 0 {
			// move left (towards lower numbers)
			for i := 0; i > rotation; i-- {
				position--
				if position < 0 {
					position = 99
				}
				if position == 0 {
					hits++
				}
			}
		}
		// rotation == 0 → do nothing
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return hits, nil
}
