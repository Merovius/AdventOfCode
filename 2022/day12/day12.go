package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"golang.org/x/exp/slices"
)

func main() {
	g, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	part1 := &Graph{
		Grid:      g,
		ValidEdge: func(a, b int) bool { return b-a <= 1 },
	}
	path := graph.BreadthFirstSearch[grid.Pos, Edge](part1, g.Start, func(p grid.Pos) bool { return p == g.End })
	if path == nil {
		log.Fatal("No path found")
	}
	fmt.Printf("Shortest path from %v to %v takes %d steps\n", g.Start, g.End, len(path))

	part2 := &Graph{
		Grid:      g,
		ValidEdge: func(a, b int) bool { return a-b <= 1 },
	}
	path = graph.BreadthFirstSearch[grid.Pos, Edge](part2, g.End, func(p grid.Pos) bool { return g.At(p) == 0 })
	if path == nil {
		log.Fatal("No path found")
	}
	fmt.Printf("Shortest path from %v to 'a' takes %d steps\n", g.End, len(path))
}

func ReadInput(r io.Reader) (*Grid, error) {
	const (
		StartMarker = -iota - 1
		EndMarker
	)
	gg, err := grid.Read(os.Stdin, func(c rune) (int, error) {
		switch {
		case c == 'S':
			return StartMarker, nil
		case c == 'E':
			return EndMarker, nil
		case 'a' <= c && c <= 'z':
			return int(c - 'a'), nil
		default:
			return 0, fmt.Errorf("invalid codepoint %q", c)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	g := &Grid{
		Grid: gg,
	}
	if i := slices.Index(g.G, StartMarker); i >= 0 {
		g.Start = g.Pos(i)
		g.Set(g.Start, 0)
	} else {
		log.Fatal("No start marker")
	}
	if i := slices.Index(g.G, EndMarker); i >= 0 {
		g.End = g.Pos(i)
		g.Set(g.End, 25)
	} else {
		log.Fatal("No start marker")
	}
	return g, nil
}

type Grid struct {
	*grid.Grid[int]
	Start grid.Pos
	End   grid.Pos
}

type Edge struct {
	From grid.Pos
	To   grid.Pos
}

type Graph struct {
	*Grid
	ValidEdge func(int, int) bool
}

func (g *Graph) Edges(p grid.Pos) []Edge {
	var out []Edge
	for _, n := range g.Neigh4(p) {
		if !g.ValidEdge(g.At(p), g.At(n)) {
			continue
		}
		out = append(out, Edge{p, n})
	}
	return out
}

func (g *Graph) From(e Edge) grid.Pos { return e.From }

func (g *Graph) To(e Edge) grid.Pos { return e.To }
