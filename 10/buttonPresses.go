package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
mental model:
Part 1:
1. Each machine has a set of indicator lights (Target State) and a set of Buttons.
2. Initial state of lights is all OFF (false/0).
3. Pushing a button toggles specific lights.
4. Goal: Reach Target State from Initial State with minimum button presses.
5. Constraint: Max 13 buttons per machine.
6. Approach: Brute-force.

Part 2:
1. "Joltage" mode. Each button INCREMENTS specific counters.
2. Target State is a set of integer values for counters.
3. Goal: Reach Target values EXACTLY from 0.
4. Minimize total presses.
5. This is a system of linear equations:
   For each counter j: Sum(x_i * Button_i_effect_on_j) = Target_j
   Where x_i is count of presses for button i. x_i >= 0.
   Objective: Minimize Sum(x_i).
6. Approach:
   - Formulate as Ax = b.
   - Use Gaussian Elimination / Substitution to simplify.
   - Since variables must be non-negative integers, this is Integer Programming.
   - However, N is small (buttons <= 13, counters <= 10).
   - "Substitution" strategy:
     - Find an equation where a variable x_k has coeff 1.
     - Rewrite x_k = T - Sum(...)
     - Check constraint x_k >= 0.
     - Substitute x_k into other equations.
     - Reduced system has fewer variables.
     - If no coeff 1, maybe branching or just simple backtracking search on remaining variables.
*/

type Machine struct {
	Target   []bool
	Buttons  [][]int
	Joltages []int
}

func main() {
	start := time.Now()

	filename := "input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	machines, err := parseInput(filename)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	totalPresses := solvePart1(machines)
	fmt.Println("Total Fewest Button Presses (Part 1):", totalPresses)

	totalPresses2 := solvePart2(machines)
	fmt.Println("Total Fewest Button Presses (Part 2):", totalPresses2)

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}

func parseInput(filename string) ([]Machine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var machines []Machine
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Example line: [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}

		// 1. Parse Target State [...]
		startBracket := strings.Index(line, "[")
		endBracket := strings.Index(line, "]")
		if startBracket == -1 || endBracket == -1 {
			continue
		}
		lightStr := line[startBracket+1 : endBracket]
		target := make([]bool, len(lightStr))
		for i, r := range lightStr {
			if r == '#' {
				target[i] = true
			}
		}

		// 2. Parse Buttons (...) and Joltages {...}
		rest := line[endBracket+1:]

		// Extract Joltages
		braceStart := strings.Index(rest, "{")
		braceEnd := strings.Index(rest, "}")
		var joltages []int
		var buttonsStr string

		if braceStart != -1 && braceEnd != -1 {
			jStr := rest[braceStart+1 : braceEnd]
			buttonsStr = rest[:braceStart]

			jParts := strings.Split(jStr, ",")
			for _, jp := range jParts {
				val, err := strconv.Atoi(strings.TrimSpace(jp))
				if err == nil {
					joltages = append(joltages, val)
				}
			}
		} else {
			buttonsStr = rest
		}

		// Split by ')' to get each button group
		buttonParts := strings.Split(buttonsStr, ")")
		var buttons [][]int

		for _, bp := range buttonParts {
			bp = strings.TrimSpace(bp)
			if !strings.HasPrefix(bp, "(") {
				continue
			}
			// Remove '('
			// Content is like "0,2,3" or "3"
			content := strings.TrimPrefix(bp, "(")
			nums := strings.Split(content, ",")
			var btn []int
			for _, nStr := range nums {
				nStr = strings.TrimSpace(nStr)
				if nStr == "" {
					continue
				}
				val, err := strconv.Atoi(nStr)
				if err == nil {
					btn = append(btn, val)
				}
			}
			buttons = append(buttons, btn)
		}

		machines = append(machines, Machine{
			Target:   target,
			Buttons:  buttons,
			Joltages: joltages,
		})
	}

	return machines, scanner.Err()
}

func solvePart1(machines []Machine) int {
	total := 0
	for _, m := range machines {
		presses := solveMachineBruteForce(m)
		if presses != math.MaxInt32 {
			total += presses
		}
	}
	return total
}

func solveMachineBruteForce(m Machine) int {
	nButtons := len(m.Buttons)
	nLights := len(m.Target)
	minPresses := math.MaxInt32
	limit := 1 << nButtons

	for i := 0; i < limit; i++ {
		currentPop := 0
		temp := i
		for temp > 0 {
			if temp&1 == 1 {
				currentPop++
			}
			temp >>= 1
		}

		if currentPop >= minPresses {
			continue
		}

		state := make([]bool, nLights)
		for bIdx := 0; bIdx < nButtons; bIdx++ {
			if (i & (1 << bIdx)) != 0 {
				for _, lightIdx := range m.Buttons[bIdx] {
					if lightIdx < nLights {
						state[lightIdx] = !state[lightIdx]
					}
				}
			}
		}

		matches := true
		for l := 0; l < nLights; l++ {
			if state[l] != m.Target[l] {
				matches = false
				break
			}
		}

		if matches {
			if currentPop < minPresses {
				minPresses = currentPop
			}
		}
	}

	return minPresses
}

// Part 2 Logic

func solvePart2(machines []Machine) int {
	total := 0
	for i, m := range machines {
		presses := solveMachineLinear(m)
		if presses == math.MaxInt32 {
			fmt.Printf("Warning: Machine %d unsolvable in Part 2\n", i)
		} else {
			total += presses
		}
	}
	return total
}

type Equation struct {
	Coeffs   []int // Coefficients for each variable x_0 ... x_n-1
	Constant int
}

// Optimized Part 2 Solver using Gaussian Elimination

func solveMachineLinear(m Machine) int {
	nButtons := len(m.Buttons)
	nCounters := len(m.Joltages)

	// Represents an equation: sum(Coeffs[v] * v) = Constant
	type LinEq struct {
		Coeffs   map[int]int
		Constant int
	}

	var equations []LinEq
	for j := 0; j < nCounters; j++ {
		coeffs := make(map[int]int)
		for i, btn := range m.Buttons {
			for _, c := range btn {
				if c == j {
					coeffs[i] = 1
				}
			}
		}
		equations = append(equations, LinEq{Coeffs: coeffs, Constant: m.Joltages[j]})
	}

	// Map of resolved variables: varIndex -> Expression
	// Expression: Constant + sum(Coeffs[v] * v)
	type Expr struct {
		Constant int
		Coeffs   map[int]int
	}
	resolved := make(map[int]Expr)

	// Gaussian Elimination
	for {
		bestEqIdx := -1
		var pivotVar int
		found := false

		// Find a variable with coefficient 1 or -1 to pivot
		// Preference: 1
		for i, eq := range equations {
			for v, c := range eq.Coeffs {
				if c == 1 || c == -1 {
					bestEqIdx = i
					pivotVar = v
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			break
		}

		// Pivot found
		eq := equations[bestEqIdx]
		pivotCoeff := eq.Coeffs[pivotVar] // 1 or -1

		// Express pivotVar = (Constant - sum(other coeffs)) / pivotCoeff
		// If pivotCoeff is 1: pivotVar = Constant - sum(other)
		// If pivotCoeff is -1: pivotVar = -Constant + sum(other)

		expr := Expr{Constant: eq.Constant * pivotCoeff, Coeffs: make(map[int]int)}
		for v, c := range eq.Coeffs {
			if v == pivotVar {
				continue
			}
			expr.Coeffs[v] = -c * pivotCoeff
		}

		// Save resolved expression
		resolved[pivotVar] = expr

		// Remove equation
		equations[bestEqIdx] = equations[len(equations)-1]
		equations = equations[:len(equations)-1]

		// Substitute pivotVar in remaining equations
		for i := range equations {
			if c, ok := equations[i].Coeffs[pivotVar]; ok {
				delete(equations[i].Coeffs, pivotVar)
				// new += c * (expr)
				// LHS term c * (expr.Constant + ...)
				// Move c * expr.Constant to RHS => - (c * expr.Constant)
				equations[i].Constant -= c * expr.Constant
				for v, k := range expr.Coeffs {
					equations[i].Coeffs[v] += c * k
					if equations[i].Coeffs[v] == 0 {
						delete(equations[i].Coeffs, v)
					}
				}
			}
		}

		// Substitute pivotVar in all previously resolved expressions
		// Resolved: v = Constant + ...
		// Substitute into RHS, so logic is additive
		for rV, rExpr := range resolved {
			if c, ok := rExpr.Coeffs[pivotVar]; ok {
				delete(rExpr.Coeffs, pivotVar)
				rExpr.Constant += c * expr.Constant
				for v, k := range expr.Coeffs {
					rExpr.Coeffs[v] += c * k
					if rExpr.Coeffs[v] == 0 {
						delete(rExpr.Coeffs, v)
					}
				}
				resolved[rV] = rExpr
			}
		}
	}

	// Identify free variables
	activeVars := make(map[int]bool)
	for i := 0; i < nButtons; i++ {
		if _, ok := resolved[i]; !ok {
			activeVars[i] = true
		}
	}
	var freeVars []int
	for v := range activeVars {
		freeVars = append(freeVars, v)
	}

	minSum := math.MaxInt32

	// Backtracking search
	var backtrack func(idx int, currentVals map[int]int)
	backtrack = func(idx int, currentVals map[int]int) {
		if idx == len(freeVars) {
			// All free vars set. Validation.

			// 1. Check consistency of remaining equations (should be 0=0 etc)
			for _, eq := range equations {
				lhs := 0
				for v, c := range eq.Coeffs {
					lhs += c * currentVals[v]
				}
				if lhs != eq.Constant {
					return
				}
			}

			totalPresses := 0
			// Sum free vars
			for _, v := range freeVars {
				totalPresses += currentVals[v]
			}

			// Compute and Sum resolved vars
			for _, expr := range resolved {
				val := expr.Constant
				for fv, c := range expr.Coeffs {
					val += c * currentVals[fv]
				}
				if val < 0 {
					return // Invalid non-negativity
				}
				totalPresses += val
			}

			if totalPresses < minSum {
				minSum = totalPresses
			}
			return
		}

		// Determine bounds for freeVars[idx]
		fv := freeVars[idx]

		// Simple bounds: 0 to some reasonable max.
		// Since total presses won't likely exceed strict upper bound of joltage sum or similar.
		// For safety, let's pick 300 (covers example and likely inputs if distinct).
		// Just being reduced to ~3 vars makes 300^3 feasible? Slower but finite.

		minVal := 0
		maxVal := 1000 // generous upper bound

		for val := minVal; val <= maxVal; val++ {
			currentVals[fv] = val

			// Optimization: Check intermediate feasibility?
			// Calculate resolved vars that are now fully determined?
			// Skipped for code conciseness, relying on problem smallness + reduction.

			// Pruning: if we already exceed minSum (partial sum)
			currentSum := 0
			for i := 0; i <= idx; i++ {
				currentSum += currentVals[freeVars[i]]
			}
			if minSum != math.MaxInt32 && currentSum >= minSum {
				break
			}

			backtrack(idx+1, currentVals)
		}
	}

	backtrack(0, make(map[int]int))
	return minSum
}
