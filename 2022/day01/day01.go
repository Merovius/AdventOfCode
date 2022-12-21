package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func main() {
	c, err := input.Slice(input.Blocks(), input.Slice(input.Lines(), input.Signed[int]())).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	a := Aggregate(c)
	fmt.Printf("The most heavily loaded elf carries %d calories\n", Max(a))
	if len(a) < 3 {
		log.Fatal("Less than 3 total elves in input")
	}
	slices.Sort(a)
	var t int
	for i := 1; i <= 3; i++ {
		t += a[len(a)-i]
	}
	fmt.Printf("The three most heavily loaded elves carry %d calories\n", t)
}

func Aggregate(c [][]int) []int {
	out := make([]int, len(c))
	for i, s := range c {
		for _, n := range s {
			out[i] += n
		}
	}
	return out
}

func Max[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		panic("Max called on empty slice")
	}
	v := s[0]
	for _, w := range s[1:] {
		if w > v {
			v = w
		}
	}
	return v
}
