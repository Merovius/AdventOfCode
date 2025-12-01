package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"gonih.org/AdventOfCode/2020/aoc"
)

func main() {
	vs, err := aoc.SlurpNumbers(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(vs) == 0 {
		log.Fatal("no adapters given")
	}
	// Using all adapters really just means sorting them in ascending order.
	sort.Ints(vs)
	// Number of 3-differences starts at 1, because our device always adds one 3-jolt difference.
	N1, N3 := 0, 1
	switch vs[0] {
	case 1:
		N1++
	case 3:
		N3++
	}
	for i := 1; i < len(vs); i++ {
		switch vs[i] - vs[i-1] {
		case 1:
			N1++
		case 2:
		case 3:
			N3++
		default:
			log.Fatalf("Difference between %d and %d jolt adapters too large", vs[i], vs[i-1])
		}
	}
	fmt.Printf("%d 1-jolt differences, %d 3-jolt differences, product is %d\n", N1, N3, N1*N3)
	fmt.Printf("Number of valid adapter stacks: %v\n", countStacks(vs))
}

func countStacks(vs []int) int {
	// memoize answers for stacks we already counted.
	memo := make(map[string]int)
	var rec func(int, []int) int
	rec = func(jolts int, rest []int) (n int) {
		if len(rest) == 0 {
			return 1
		}
		k := fmt.Sprintf("%v%v", jolts, rest)
		if n, ok := memo[k]; ok {
			return n
		}
		for i := 0; i < len(rest); i++ {
			if rest[i] > jolts+3 {
				break
			}
			n += rec(rest[i], rest[i+1:])
		}
		memo[k] = n
		return n
	}
	return rec(0, vs)
}
