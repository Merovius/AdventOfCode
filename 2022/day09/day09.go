package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	moves, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With 2 pieces, the tail visited %d positions\n", Simulate(moves, 2))
	fmt.Printf("With 10 pieces, the tail visited %d positions\n", Simulate(moves, 10))
}

func ReadInput(r io.Reader) ([]Move, error) {
	var out []Move
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		first, second, ok := strings.Cut(l, " ")
		if !ok {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		n, err := strconv.Atoi(second)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		if len(first) != 1 || strings.IndexByte("UDLR", first[0]) < 0 {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		out = append(out, Move{Direction(first[0]), n})
	}
	return out, s.Err()
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
	if abs(δx) > 1 || abs(δy) > 1 {
		t.X += sgn(δx)
		t.Y += sgn(δy)
	}
	return t
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sgn(v int) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	default:
		return 0
	}
}
