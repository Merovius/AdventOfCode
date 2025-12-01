package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"strings"

	"gonih.org/AdventOfCode/internal/graph"
	"gonih.org/AdventOfCode/internal/grid"
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

func Parse(s string) (*Graph, error) {
	g, err := grid.Read[int8](strings.NewReader(s), func(r rune) (int8, error) {
		if !('0' <= r && r <= '9') {
			return 0, fmt.Errorf("invalid input %q", r)
		}
		return int8(r - '0'), nil
	})
	if err != nil {
		return nil, err
	}
	return &Graph{Grid: g}, nil
}

func Part1(g *Graph) int {
	g.Min = 0
	g.Max = 3
	start := Node{
		P: grid.Pos{0, 0},
		Δ: grid.Pos{},
	}
	path := graph.Dijkstra[Node, Edge, int](g, start, func(n Node) bool {
		return n.P.Row == g.H-1 && n.P.Col == g.W-1
	})
	var total int
	for _, e := range path {
		total += e.Loss
	}
	return total
}

func Part2(g *Graph) int {
	g.Min = 3
	g.Max = 10
	start := Node{
		P: grid.Pos{0, 0},
		Δ: grid.Pos{},
	}
	path := graph.Dijkstra[Node, Edge, int](g, start, func(n Node) bool {
		return n.P.Row == g.H-1 && n.P.Col == g.W-1
	})
	var total int
	for _, e := range path {
		total += e.Loss
	}
	return total
}

type Graph struct {
	*grid.Grid[int8]
	Min int
	Max int
}

type Node struct {
	P grid.Pos
	Δ grid.Pos
}

type Edge struct {
	From Node
	To   Node
	Loss int
}

func (g *Graph) Edges(n Node) iter.Seq[Edge] {
	return func(yield func(Edge) bool) {
		for _, δ := range []grid.Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			if δ == n.Δ || (δ.Row == -n.Δ.Row && δ.Col == -n.Δ.Col) {
				continue
			}
			var (
				loss int
				p    = n.P
			)
			for i := range g.Max {
				p = p.Add(δ)
				if !g.Valid(p) {
					break
				}
				loss += int(g.At(p))
				if i >= g.Min {
					e := Edge{
						From: n,
						To:   Node{p, δ},
						Loss: loss,
					}
					if !yield(e) {
						return
					}
				}
			}
		}
	}
}

func (g *Graph) From(e Edge) Node {
	return e.From
}

func (g *Graph) To(e Edge) Node {
	return e.To
}

func (g *Graph) Weight(e Edge) int {
	return e.Loss
}
