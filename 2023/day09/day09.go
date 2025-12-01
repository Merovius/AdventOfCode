package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	part1, part2 := Solve(in)
	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}

func Parse(s string) ([][]int, error) {
	return parse.TrimSpace(parse.Slice(
		split.Lines,
		parse.Slice(
			split.Fields,
			parse.Signed[int],
		),
	))(s)
}

func Solve(in [][]int) (part1, part2 int) {
	for _, h := range in {
		δ2, δ1 := Predict(h)
		part1, part2 = part1+δ1, part2+δ2
	}
	return part1, part2
}

func Predict(h []int) (left, right int) {
	if allZero(h) {
		return 0, 0
	}
	a, b := h[0], h[len(h)-1]
	next := h[:0]
	for i := 1; i < len(h); i++ {
		next = append(next, h[i]-h[i-1])
	}
	δa, δb := Predict(next)
	return a - δa, b + δb
}

func allZero(s []int) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}
