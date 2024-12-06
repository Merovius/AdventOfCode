package main

import (
	"fmt"
	"iter"
	"log"
	"os"
	"slices"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	g, err := grid.Read(os.Stdin, func(r rune) (Cell, error) {
		return Cell(r), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	var (
		start   = g.Pos(slices.Index(g.G, Position))
		d0      = grid.Pos{-1, 0}
		visited = make(set.Set[grid.Pos])
	)
	g.Set(start, Empty)
	for p, _ := range Run(g, start, d0) {
		visited.Add(p)
	}
	fmt.Println(len(visited))

	obsts := make(set.Set[grid.Pos])
	for p, _ := range Run(g, start, d0) {
		if p == start {
			continue
		}
		g.Set(p, Obstruction)
		if IsLoop(Run(g, start, d0)) {
			obsts.Add(p)
		}
		g.Set(p, Empty)
	}
	fmt.Println(len(obsts))
}

type Cell byte

const (
	Empty       = '.'
	Obstruction = '#'
	Position    = '^'
)

func Run(g *grid.Grid[Cell], p, d grid.Pos) iter.Seq2[grid.Pos, grid.Pos] {
	return func(yield func(grid.Pos, grid.Pos) bool) {
		for g.Valid(p) {
			if !yield(p, d) {
				return
			}
			q := p.Add(d)
			for g.Valid(q) && g.At(q) == Obstruction {
				d.Row, d.Col = d.Col, -d.Row
				q = p.Add(d)
			}
			p = q
		}
	}
}

func IsLoop(path iter.Seq2[grid.Pos, grid.Pos]) bool {
	seen := make(set.Set[[2]grid.Pos])
	for p, d := range path {
		pd := [2]grid.Pos{p, d}
		if seen.Contains(pd) {
			return true
		}
		seen.Add(pd)
	}
	return false
}

func Print(g *grid.Grid[Cell], p, d grid.Pos) {
	var q grid.Pos
	for q.Row = 0; q.Row < g.H; q.Row++ {
		for q.Col = 0; q.Col < g.W; q.Col++ {
			if q == p {
				switch d {
				case grid.Pos{-1, 0}:
					fmt.Print("^")
				case grid.Pos{0, 1}:
					fmt.Print(">")
				case grid.Pos{1, 0}:
					fmt.Print("V")
				case grid.Pos{0, -1}:
					fmt.Print("<")
				}
				continue
			}
			fmt.Print(string(g.At(q)))
		}
		fmt.Println()
	}
	fmt.Println()
}
