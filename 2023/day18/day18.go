package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/math"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	part1, part2, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Enclosed(part1))
	fmt.Println("Part 2:", Enclosed(part2))
}

func Parse(s string) ([]Op, []Op, error) {
	var part1, part2 []Op

	for len(s) > 0 {
		l, rest, _ := strings.Cut(s, "\n")
		s = rest
		f := strings.Fields(l)
		if len(f) != 3 {
			return nil, nil, fmt.Errorf("invalid line %q: does not have three fields", l)
		}
		var (
			o   Op
			err error
		)
		switch f[0] {
		case "U":
			o.Direction = Up
		case "D":
			o.Direction = Down
		case "L":
			o.Direction = Left
		case "R":
			o.Direction = Right
		default:
			return nil, nil, fmt.Errorf("invalid line %q: invalid direction %q", l, f[0])
		}

		o.Distance, err = strconv.Atoi(f[1])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid line %q: %w", l, err)
		}
		part1 = append(part1, o)

		if len(f[2]) != 9 {
			return nil, nil, fmt.Errorf("invalid line %q: last field has not length 9", l)
		}
		n, err := strconv.ParseInt(f[2][2:7], 16, 32)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid line %q: %w", l, err)
		}
		o.Distance = int(n)
		switch f[2][7] {
		case '0':
			o.Direction = Right
		case '1':
			o.Direction = Down
		case '2':
			o.Direction = Left
		case '3':
			o.Direction = Up
		default:
			return nil, nil, fmt.Errorf("invalid line %q: invalid direction %q", l, f[2][8])
		}
		part2 = append(part2, o)
	}
	return part1, part2, nil
}

func Enclosed(ops []Op) int {
	// https://en.wikipedia.org/wiki/Shoelace_formula#Triangle_formula
	// The border needs extra accounting. We can imagine the corners of the
	// polygon in the middle of the cell. The border then cuts cells in half,
	// so for every border cell, we have to add half a cell to the total. Every
	// convex corner cell furthermore adds another quarter cell and every
	// concave corner cell overcounts by a quarter cell. There are always 4
	// more convex than concave corner cells to close the polygon, so we have
	// to add one more to the total.
	var (
		p grid.Pos
		A int
		b int
	)
	for _, o := range ops {
		q := o.Apply(p)
		A += p.Row*q.Col - q.Row*p.Col
		b += o.Distance
		p = q
	}
	return math.Abs(A)/2 + b/2 + 1
}

type Plan struct {
	Plan    []Op
	HexPlan []Op
}

type Op struct {
	Direction Direction
	Distance  int
}

func (o Op) Apply(p grid.Pos) grid.Pos {
	var δ grid.Pos
	switch o.Direction {
	case Up:
		δ = grid.Pos{-o.Distance, 0}
	case Down:
		δ = grid.Pos{o.Distance, 0}
	case Left:
		δ = grid.Pos{0, -o.Distance}
	case Right:
		δ = grid.Pos{0, o.Distance}
	}
	return p.Add(δ)
}

type Direction int8

const (
	_ Direction = iota
	Up
	Down
	Left
	Right
)
