# Brainstorming / Notes

## Current Approach (Brute Force)
- Iterate through every number in the range
- Cast to string
- Try every possible pattern length `L` from 1 up to `len/2`
- Check if the string is just that pattern repeated `k` times
- `s == repeat(s[:L], count)`
- Logically simple, but O(N^2) pattern matching?

## Optimization Idea: Sliding Window / Hashing (copy paste)
- Rabin-Karp style?? (thx cs stackexchange)
- Use a **Rolling Hash**
- Precompute hash for the whole string
- Allows getting hash of *any* substring in O(1) time
- The trick:
    - To check if pattern of length `L` repeats...
    - Don't need to compare chars
    - Just compare Hash(Prefix of length N-L) vs Hash(Suffix of length N-L)
    - If they match -> the pattern repeats!
- `hash(s[0...N-L]) == hash(s[L...N])`
- Way faster, no string slicing or comparing loops
- O(1) check per pattern length
- Need to handle collisions? Probably rare enough to ignore for this or double check on match.
