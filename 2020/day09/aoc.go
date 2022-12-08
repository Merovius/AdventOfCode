package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/Merovius/AdventOfCode/2020/aoc"
)

func main() {
	vs, err := aoc.SlurpNumbers(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(vs) < bufSize {
		log.Fatal("not enough numbers")
	}
	v, ok := findInvalidNumber(vs)
	if !ok {
		log.Fatal("No invalid number found")
	}
	fmt.Println("Invalid number:", v)
	sum := findSum(vs, v)
	fmt.Println("Range:", sum)
	min, max := bounds(sum)
	fmt.Printf("%d + %d = %d\n", min, max, min+max)
}

const bufSize = 25

func findInvalidNumber(vs []int) (v int, ok bool) {
loop:
	for i := bufSize; i < len(vs); i++ {
		for j := i - bufSize; j < i; j++ {
			for k := j + 1; k < i; k++ {
				if vs[j]+vs[k] == vs[i] {
					continue loop
				}
			}
		}
		return vs[i], true
	}
	return 0, false
}

func findSum(vs []int, v int) []int {
	for i := 0; i < len(vs); i++ {
		sum := 0
		for j := i; j < len(vs); j++ {
			sum += vs[j]
			if sum == v {
				return vs[i : j+1]
			}
		}
	}
	return nil
}

func bounds(vs []int) (min, max int) {
	min, max = math.MaxInt64, math.MinInt64
	for _, v := range vs {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
