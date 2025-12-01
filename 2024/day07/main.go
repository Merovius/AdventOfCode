package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/math"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

type Equation struct {
	Test   int
	Inputs []int
}

func Parse(in []byte) ([]Equation, error) {
	return parse.Lines(
		parse.Struct[Equation](
			split.On(": "),
			parse.Signed[int],
			parse.Slice(split.Fields, parse.Signed[int]),
		),
	)(string(in))
}

func Part1(in []Equation) int {
	var total int
	for _, e := range in {
		if solvable(e.Test, e.Inputs) {
			total += e.Test
		}
	}
	return total
}

func solvable(res int, ins []int) bool {
	n, ins := ins[len(ins)-1], ins[:len(ins)-1]
	if len(ins) == 0 {
		return res == n
	}
	if solvable(res-n, ins) {
		return true
	}
	return res%n == 0 && solvable(res/n, ins)
}

func Part2(in []Equation) int {
	var total int
	for _, e := range in {
		if solvable2(e.Test, e.Inputs) {
			total += e.Test
		}
	}
	return total
}

func solvable2(res int, ins []int) bool {
	n, ins := ins[len(ins)-1], ins[:len(ins)-1]
	if len(ins) == 0 {
		return res == n
	}
	if res >= n {
		if solvable2(res-n, ins) {
			return true
		}
	}
	if res%n == 0 {
		if solvable2(res/n, ins) {
			return true
		}
	}
	p := int(math.Pow10(math.Digits(n)))
	return res%p == n && solvable2(res/p, ins)
}
