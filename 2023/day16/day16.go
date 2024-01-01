package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
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
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(s string) (*Grid, error) {
	return grid.Read[Cell](strings.NewReader(s), ParseCell)
}

func Part1(g *Grid) int {
	return Energize(g, Beam{grid.Pos{0, 0}, Right})
}

func Part2(g *Grid) int {
	ch := make(chan int, 2*g.W+2*g.H)
	energize := func(b Beam) {
		ch <- Energize(g, b)
	}
	for c := 0; c < g.W; c++ {
		go energize(Beam{grid.Pos{0, c}, Down})
		go energize(Beam{grid.Pos{g.H - 1, c}, Up})
	}
	for r := 0; r < g.H; r++ {
		go energize(Beam{grid.Pos{r, 0}, Right})
		go energize(Beam{grid.Pos{r, g.W - 1}, Left})
	}
	var m int
	for i := 0; i < 2*g.W+2*g.H; i++ {
		m = max(m, <-ch)
	}
	return m
}

type Grid = grid.Grid[Cell]

type Cell uint8

const (
	Empty   Cell = iota
	SplitH       // '-'
	SplitV       // '|'
	MirrorR      // '/'
	MirrorL      // '\'
)

func (c Cell) String() string {
	switch c {
	case Empty:
		return " "
	case SplitH:
		return "-"
	case SplitV:
		return "|"
	case MirrorR:
		return "/"
	case MirrorL:
		return "\\"
	}
	panic("invalid cell")
}

func ParseCell(r rune) (Cell, error) {
	switch r {
	case '.':
		return Empty, nil
	case '-':
		return SplitH, nil
	case '|':
		return SplitV, nil
	case '/':
		return MirrorR, nil
	case '\\':
		return MirrorL, nil
	}
	return 0, fmt.Errorf("invalid input %q", r)
}

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "^"
	case Down:
		return "v"
	case Left:
		return "<"
	case Right:
		return ">"
	}
	panic("invalid Direction")
}

func (d Direction) Of(p grid.Pos) grid.Pos {
	switch d {
	case Up:
		p.Row--
	case Down:
		p.Row++
	case Left:
		p.Col--
	case Right:
		p.Col++
	}
	return p
}

func (d Direction) Vert() bool {
	return d == Up || d == Down
}

func (d Direction) Horiz() bool {
	return d == Left || d == Right
}

func (d Direction) Interact(c Cell) []Direction {
	switch c {
	case SplitH: // '-'
		if d.Vert() {
			return []Direction{Left, Right}
		}
	case SplitV: // '|'
		if d.Horiz() {
			return []Direction{Up, Down}
		}
	case MirrorR: // '/'
		switch d {
		case Up:
			d = Right
		case Down:
			d = Left
		case Left:
			d = Down
		case Right:
			d = Up
		}
	case MirrorL: // '\'
		switch d {
		case Up:
			d = Left
		case Down:
			d = Right
		case Left:
			d = Up
		case Right:
			d = Down
		}
	}
	return []Direction{d}
}

type Beam struct {
	p grid.Pos
	d Direction
}

func (b Beam) String() string {
	return fmt.Sprintf("{%v, %v}", b.p, b.d)
}

func Energize(g *Grid, b Beam) int {
	var (
		q    container.LIFO[Beam]
		seen = make(set.Set[Beam])
	)
	push := func(b Beam) {
		if seen.Contains(b) || !g.Valid(b.p) {
			return
		}
		seen.Add(b)
		q.Push(b)
	}
	push(b)
	for q.Len() > 0 {
		b := q.Pop()
		for _, d := range b.d.Interact(g.At(b.p)) {
			push(Beam{d.Of(b.p), d})
		}
	}

	energized := make(set.Set[grid.Pos])
	for b := range seen {
		energized.Add(b.p)
	}
	return len(energized)
}
