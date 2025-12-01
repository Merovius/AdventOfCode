package main

import (
	"fmt"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	. "gonih.org/AdventOfCode/internal/math"
	"gonih.org/AdventOfCode/internal/set"
)

func main() {
	moves, err := parse.Lines(
		parse.Struct[Move](
			split.Fields,
			parse.Enum(Up, Down, Left, Right),
			parse.Signed[int],
		),
	).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With 2 pieces, the tail visited %d positions\n", Simulate(moves, 2))
	fmt.Printf("With 10 pieces, the tail visited %d positions\n", Simulate(moves, 10))
}

type Direction byte

const (
	Up    Direction = 'U'
	Down  Direction = 'D'
	Left  Direction = 'L'
	Right Direction = 'R'
)

type Move struct {
	Direction
	N int
}

type Pos struct {
	X int
	Y int
}

func Delta(a, b Pos) (δx, δy int) {
	return b.X - a.X, b.Y - a.Y
}

func Simulate(moves []Move, n int) (visited int) {
	knots := make([]Pos, n)
	v := make(set.Set[Pos])
	for _, m := range moves {
		for i := 0; i < m.N; i++ {
			knots[0] = Step(knots[0], m.Direction)
			for j := 1; j < n; j++ {
				knots[j] = Pull(knots[j-1], knots[j])
			}
			v.Add(knots[n-1])
		}
	}
	return len(v)
}

func Step(p Pos, d Direction) Pos {
	switch d {
	case Up:
		p.Y -= 1
	case Down:
		p.Y += 1
	case Left:
		p.X -= 1
	case Right:
		p.X += 1
	}
	return p
}

func Pull(h, t Pos) Pos {
	δx, δy := Delta(t, h)
	if Abs(δx) > 1 || Abs(δy) > 1 {
		t.X += Sgn(δx)
		t.Y += Sgn(δy)
	}
	return t
}
