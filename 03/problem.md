# Lobby, December 3
[LINK](https://adventofcode.com/2025/day/3)

## Given
- Batteries have joltage rating 1-9
- Rating noted as such:
```battery
987654321111111
811111111111119
234234234234278
818181911112111
```
- Batteries are arranged into banks, each line of digits (see above) is one bank
- Within each bank, two batteries need to be turned on
- The joltage each bank produces is equal to the number formed by digits on the batteries turned on
- Example: if you have a bank like ``12345`` and you turn on batteries ``2`` and ``4``, the bank would produce ``24`` jolts

## Find (p1)
- Find maximum joltage possible for each bank and sum 

## Notes