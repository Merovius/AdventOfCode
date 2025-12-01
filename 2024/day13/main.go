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
	ms, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(ms))
	fmt.Println(Part2(ms))
}

func Parse(in []byte) ([]Machine, error) {
	return parse.Slice(
		split.Blocks,
		parse.Struct[Machine](
			split.Lines,
			parse.Array[[2]int](
				split.Regexp(`Button A: X\+(\d+), Y\+(\d+)`),
				parse.Signed[int],
			),
			parse.Array[[2]int](
				split.Regexp(`Button B: X\+(\d+), Y\+(\d+)`),
				parse.Signed[int],
			),
			parse.Array[[2]int](
				split.Regexp(`Prize: X=(\d+), Y=(\d+)`),
				parse.Signed[int],
			),
		),
	)(string(in))

}

type Machine struct {
	ButtonA [2]int
	ButtonB [2]int
	Prize   [2]int
}

func Part1(ms []Machine) int {
	var total int
	for _, m := range ms {
		x, ok := math.Cramer2(m.ButtonA, m.ButtonB, m.Prize)
		if !ok {
			continue
		}
		total += x[0]*3 + x[1]
	}
	return total
}

func Part2(ms []Machine) int {
	var total int
	for _, m := range ms {
		m.Prize[0] += 10000000000000
		m.Prize[1] += 10000000000000
		x, ok := math.Cramer2(m.ButtonA, m.ButtonB, m.Prize)
		if !ok {
			continue
		}
		total += x[0]*3 + x[1]
	}
	return total
}
