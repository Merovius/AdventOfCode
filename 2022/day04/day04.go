package main

import (
	"fmt"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

func main() {
	input, err := parse.Lines(
		parse.Array[[2]Range](
			split.On(","),
			parse.Struct[Range](
				split.On("-"),
				parse.Signed[int],
				parse.Signed[int],
			),
		),
	).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There are %d pairs where one contains the other\n", len(Filter(input, EitherContains)))
	fmt.Printf("There are %d overlapping pairs\n", len(Filter(input, Overlap)))
}

type Range struct {
	Min int
	Max int
}

func Contains(a, b Range) bool {
	return a.Min <= b.Min && a.Max >= b.Max
}

func EitherContains(a, b Range) bool {
	return Contains(a, b) || Contains(b, a)
}

func Overlap(a, b Range) bool {
	return !(a.Max < b.Min || b.Max < a.Min)
}

func Filter(p [][2]Range, include func(a, b Range) bool) [][2]Range {
	var out [][2]Range
	for _, p := range p {
		if include(p[0], p[1]) {
			out = append(out, p)
		}
	}
	return out
}
