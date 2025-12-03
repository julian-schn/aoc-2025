package main

import (
	"fmt"
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

func main() {
	start := time.Now()

	// Your code logic here

	elapsed := time.Since(start)
	fmt.Println("Runtime:", elapsed)
}
