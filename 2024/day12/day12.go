package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	g, err := Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(g))
	fmt.Println(Part2(g))
}

type Cell byte

func Parse(r io.Reader) (*grid.Grid[Cell], error) {
	return grid.Read(r, func(r rune) (Cell, error) {
		return Cell(r), nil
	})
}

func Part1(g *grid.Grid[Cell]) int {
	var price int
	for area, perimeter := range Areas(g) {
		price += area * len(perimeter)
	}
	return price
}

func Part2(g *grid.Grid[Cell]) int {
	var price int
	for area, perimeter := range Areas(g) {
		var sides int
		for e := range perimeter {
			sides++
			δ := e.Δ()
			for f := e.Add(δ); perimeter.Contains(f); f = f.Add(δ) {
				delete(perimeter, f)
			}
			for f := e.Sub(δ); perimeter.Contains(f); f = f.Sub(δ) {
				delete(perimeter, f)
			}
			delete(perimeter, e)
		}
		price += area * sides
	}
	return price
}

type Edge [2]grid.Pos

// Δ returns an offset to add to e, to continue along its side.
func (e Edge) Δ() grid.Pos {
	δ := e[1].Sub(e[0])
	δ.Row, δ.Col = -δ.Col, δ.Row
	return δ
}

func (e Edge) Add(δ grid.Pos) Edge {
	e[0] = e[0].Add(δ)
	e[1] = e[1].Add(δ)
	return e
}

func (e Edge) Sub(δ grid.Pos) Edge {
	e[0] = e[0].Sub(δ)
	e[1] = e[1].Sub(δ)
	return e
}

func Areas(g *grid.Grid[Cell]) iter.Seq2[int, set.Set[Edge]] {
	return func(yield func(int, set.Set[Edge]) bool) {
		seen := make(set.Set[grid.Pos])
		for p, _ := range g.All() {
			if seen.Contains(p) {
				continue
			}
			var (
				area      int
				perimeter = make(set.Set[Edge])
			)
			var rec func(grid.Pos)
			rec = func(p grid.Pos) {
				if seen.Contains(p) {
					return
				}
				seen.Add(p)
				area++
				for _, δ := range []grid.Pos{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
					q := p.Add(δ)
					if !g.Valid(q) || g.At(p) != g.At(q) {
						perimeter.Add(Edge{p, q})
						continue
					}
					rec(q)
				}
			}
			rec(p)
			if !yield(area, perimeter) {
				return
			}
		}
	}
}
