package main

import (
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
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
		x, ok := Solve(m.ButtonA, m.ButtonB, m.Prize)
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
		x, ok := Solve(m.ButtonA, m.ButtonB, m.Prize)
		if !ok {
			continue
		}
		total += x[0]*3 + x[1]
	}
	return total
}

func Solve(a, b, c [2]int) (x [2]int, ok bool) {
	// Cramer's rule
	d := a[0]*b[1] - b[0]*a[1]
	if d == 0 {
		return [2]int{}, false
	}
	x[0] = c[0]*b[1] - b[0]*c[1]
	if x[0]%d != 0 {
		return [2]int{}, false
	}
	x[1] = a[0]*c[1] - c[0]*a[1]
	if x[1]%d != 0 {
		return [2]int{}, false
	}
	x[0], x[1] = x[0]/d, x[1]/d
	return x, true
}
