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

type Cell byte

const (
	Empty       = '.'
	Obstruction = '#'
	Position    = '^'
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
		d0      = grid.Up
		visited = make(set.Set[grid.Pos])
	)
	g.Set(start, Empty)
	for p := range Run(g, start, d0) {
		visited.Add(p)
	}
	fmt.Println(len(visited))

	obsts := make(set.Set[grid.Pos])
	for p := range Run(g, start, d0) {
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

func Run(g *grid.Grid[Cell], p grid.Pos, d grid.Direction) iter.Seq2[grid.Pos, grid.Direction] {
	return func(yield func(grid.Pos, grid.Direction) bool) {
		for g.Valid(p) {
			if !yield(p, d) {
				return
			}
			q := d.Move(p)
			for g.Valid(q) && g.At(q) == Obstruction {
				d = d.RotateRight()
				q = d.Move(p)
			}
			p = q
		}
	}
}

func IsLoop(path iter.Seq2[grid.Pos, grid.Direction]) bool {
	type node struct {
		grid.Pos
		grid.Direction
	}
	seen := make(set.Set[node])
	for p, d := range path {
		n := node{p, d}
		if seen.Contains(n) {
			return true
		}
		seen.Add(n)
	}
	return false
}

func Print(g *grid.Grid[Cell], p grid.Pos, d grid.Direction) {
	var q grid.Pos
	for q.Row = 0; q.Row < g.H; q.Row++ {
		for q.Col = 0; q.Col < g.W; q.Col++ {
			if q == p {
				fmt.Print(d)
				continue
			}
			fmt.Print(string(g.At(q)))
		}
		fmt.Println()
	}
	fmt.Println()
}
