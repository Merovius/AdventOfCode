package main

import (
	"errors"
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
	g, start, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(g, start))
	fmt.Println("Part 2:", Part2(g, start))
}

func Parse(s string) (*grid.Grid[Cell], grid.Pos, error) {
	g, err := grid.Read(strings.NewReader(s), grid.Enum[Cell]('.', '#', 'S'))
	if err != nil {
		return nil, grid.Pos{}, err
	}
	start := grid.Pos{-1, -1}
	for p := range grid.Find(g, Start) {
		if start != (grid.Pos{-1, -1}) {
			return nil, grid.Pos{}, errors.New("multiple start positions")
		}
		start = p
	}
	g.Set(start, Garden)
	return g, start, nil
}

func Part1(g *grid.Grid[Cell], start grid.Pos) int {
	return NReachable(g, start, 64)
}

func Part2(g *grid.Grid[Cell], start grid.Pos) int {
	// Due to the structure of the input, the reachable nodes form a diamond
	// every w•i+w/2 steps (where w is the width of the grid and i is a natural
	// number) - and in particular, n%w == w/2 for the interesting number of
	// steps.
	//
	// The area of reachable nodes for those diamonds is proportional to the
	// radius of the diamond (i.e. number of used steps), so it is a second
	// degree polynomial.
	//
	// The algebra is easier, if we assume the input is i in the w•i+w/2. We
	// map the input from the step number in the final polynomial.
	//
	// This logic only works for the actual input, not for the example.
	n := 26501365
	w := g.W
	if w/2 != n%w {
		panic(fmt.Errorf("%d/2 = %d but %d%%%d = %d", w, w/2, n, g.W, n%w))
	}
	x := []int{
		n%w + w*0,
		n%w + w*1,
		n%w + w*2,
	}
	y := []int{
		NReachable(&cover{g}, start, x[0]),
		NReachable(&cover{g}, start, x[1]),
		NReachable(&cover{g}, start, x[2]),
	}
	// three points are enough to interpolate a polynomial
	// we have p(i) = y[i] = a•x[i]²+b•x[i]+c

	// p(0) = a•0²+b•0+c
	c := y[0]
	// p(1) = a+b+c
	// p(2) = 4a+2b+c
	// ⇒ 4•p(1)-p(2) = 2•b+3•c
	// ⇒ b = (4•p(1)-p(2)-3•c)/2
	b := (4*y[1] - y[2] - 3*c) / 2
	// p(1) = a+b+c
	// ⇒ a = p(1)-b-c
	a := y[1] - b - c

	p := func(v int) int {
		return a*v*v + b*v + c
	}
	// map actual step numbers into the amount we are interested in.
	// This isn't correct for step numbers not of the w•i+n%w form, but
	// whatever.
	q := func(v int) int {
		return p((v - n%w) / w)
	}
	for i := 0; i < 3; i++ {
		if z := q(x[i]); z != y[i] {
			panic(fmt.Errorf("p(%d) = %d, want %d", x[i], z, y[i]))
		}
	}
	return q(n)
}

type Cell uint8

const (
	Garden Cell = iota
	Rocks
	Start
	Reachable
)

type Grid interface {
	At(grid.Pos) Cell
	Neigh4(grid.Pos) []grid.Pos
}

// universal cover of Grid - presents (r,c) as (r%g.H, c%g.C) of the underlying
// grid.
type cover struct {
	*grid.Grid[Cell]
}

func (w *cover) At(p grid.Pos) Cell {
	p.Row = p.Row % w.H
	if p.Row < 0 {
		p.Row += w.H
	}
	p.Col = p.Col % w.W
	if p.Col < 0 {
		p.Col += w.W
	}
	return w.Grid.At(p)
}

func (w *cover) Neigh4(p grid.Pos) []grid.Pos {
	out := make([]grid.Pos, 0, 4)
	for _, δ := range []grid.Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
		q := p.Add(δ)
		if w.At(q) != Rocks {
			out = append(out, q)
		}
	}
	return out
}

func NReachable(g Grid, start grid.Pos, n int) int {
	type node struct {
		pos   grid.Pos
		steps int
	}
	var (
		q     container.FIFO[node]
		seen  = make(set.Set[grid.Pos])
		total int
	)

	q.Push(node{start, 0})
	seen.Add(start)
	for q.Len() > 0 {
		st := q.Pop()
		if st.steps > n {
			continue
		}
		if (st.steps % 2) == (n % 2) {
			total++
		}
		for _, neigh := range g.Neigh4(st.pos) {
			if g.At(neigh) != Garden {
				continue
			}
			if seen.Contains(neigh) {
				continue
			}
			seen.Add(neigh)
			q.Push(node{neigh, st.steps + 1})
		}
	}
	return total
}
