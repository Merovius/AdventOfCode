package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/grid"
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
	start := Node{
		P: grid.Pos{0, 0},
		Δ: grid.Pos{},
		N: 0,
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
	g.Ultra = true
	start := Node{
		P: grid.Pos{0, 0},
		Δ: grid.Pos{},
		N: 0,
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
	Ultra bool
}

type Node struct {
	P grid.Pos
	Δ grid.Pos
	N int
}

type Edge struct {
	From Node
	To   Node
	Loss int
}

func (g *Graph) Edges(n Node) []Edge {
	if g.Ultra {
		return g.UltraEdges(n)
	}
	var out []Edge
	for _, δ := range []grid.Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
		q := n.P.Add(δ)
		if !g.Valid(q) {
			continue
		}
		e := Edge{
			From: n,
			To:   Node{P: q, Δ: δ},
			Loss: int(g.At(q)),
		}
		if δ.Row == -n.Δ.Row && δ.Col == -n.Δ.Col {
			continue
		} else if δ == n.Δ {
			if n.N == 3 {
				continue
			}
			e.To.N = n.N + 1
		} else {
			e.To.N = 1
		}
		out = append(out, e)
	}
	return out
}

func (g *Graph) UltraEdges(n Node) []Edge {
	var out []Edge
	if n.N > 0 && n.N < 10 {
		q := n.P.Add(n.Δ)
		if g.Valid(q) {
			out = append(out, Edge{
				From: n,
				To:   Node{q, n.Δ, n.N + 1},
				Loss: int(g.At(q)),
			})
		}
	}
directions:
	for _, δ := range []grid.Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
		if δ == n.Δ || δ.Row == -n.Δ.Row && δ.Col == -n.Δ.Col {
			continue
		}
		q := n.P
		var loss int
		for i := 0; i < 4; i++ {
			q = q.Add(δ)
			if !g.Valid(q) {
				continue directions
			}
			loss += int(g.At(q))
		}
		out = append(out, Edge{
			From: n,
			To:   Node{q, δ, 4},
			Loss: loss,
		})
	}
	return out
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
