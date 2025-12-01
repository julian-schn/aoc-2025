# Secret Entrance, December 1
[LINK](https://adventofcode.com/2025/day/1)

## Given
- 100-position dial labeled 0–99, wraps both directions; start at 50.
- Commands: `L` subtracts distance, `R` adds distance (mod 100); one per line.
- Password depends on how often the dial hits 0.

## Part One
- Count times the dial ends on 0 after each rotation (ignore clicks in between).

## Part Two
- Count every click landing on 0, including during rotations and at the end.
- Quirk: large moves can hit 0 many times (e.g., from 50, `R1000` hits 0 ten times).
