# Gift Shop, December 2
[LINK](https://adventofcode.com/2025/day/2)

## Given
- Numeric IDs
- Some are invalid
- ID ranges are given, within which checks need to be done
- Ranges are seperated by comma

## Find (p1)
- Invalid IDs, invalid are any IDs that repeat the same pattern of numbers twice
- Sum of all invalid IDs

## Edge Cases (p1)
- IDs with leading zeroes are always invalid

## Find (p2)
- Any IDs where a pattern repeats at least twice are **also** invalid, even if single character