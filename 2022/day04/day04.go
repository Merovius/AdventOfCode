package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	input, err := ReadInput(os.Stdin)
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

func ReadInput(r io.Reader) ([][2]Range, error) {
	var out [][2]Range
	for {
		var p [2]Range
		i, err := fmt.Scanf("%d-%d,%d-%d\n", &p[0].Min, &p[0].Max, &p[1].Min, &p[1].Max)
		if err == io.EOF {
			return out, nil
		}
		if i != 4 || err != nil {
			return nil, err
		}
		out = append(out, p)
	}
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
