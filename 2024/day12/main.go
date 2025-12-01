package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/container"
	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/set"
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
	var (
		price     int
		perimeter int
	)
	inc := func(_ Edge) { perimeter++ }
	for area := range Areas(g, inc) {
		price += area * perimeter
		perimeter = 0
	}
	return price
}

func Part2(g *grid.Grid[Cell]) int {
	perimeter := make(set.Set[Edge])
	var price int
	for area := range Areas(g, perimeter.Add) {
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

// Edge separating two points.
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

// Areas finds all regions in g, yielding their area. perimeter is called for
// each Edge in the perimeter of the currently found region.
func Areas(g *grid.Grid[Cell], perimeter func(Edge)) iter.Seq[int] {
	return func(yield func(int) bool) {
		seen := make(set.Set[grid.Pos])
		var queue container.LIFO[grid.Pos]
		for p, _ := range g.All() {
			if seen.Contains(p) {
				continue
			}
			var area int
			queue.Push(p)
			for queue.Len() > 0 {
				p := queue.Pop()
				if seen.Contains(p) {
					continue
				}
				seen.Add(p)
				area++
				for _, δ := range []grid.Pos{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
					q := p.Add(δ)
					if !g.Valid(q) || g.At(p) != g.At(q) {
						perimeter(Edge{p, q})
						continue
					}
					queue.Push(q)
				}
			}
			if !yield(area) {
				return
			}
		}
	}
}
