//go:build goexperiment.rangefunc

package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/container"
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

func Parse(s string) (*Grid, error) {
	return grid.Read(strings.NewReader(s), grid.Enum[Cell]('#', '.', '^', '>', 'v', '<'))
}

type Cell uint8

const (
	Forest Cell = iota
	Path
	SlopeU
	SlopeR
	SlopeD
	SlopeL
)

type Grid = grid.Grid[Cell]

func Part1(g *Grid) int {
	return LongestPath(MakeGraph(g))
}

func Part2(g *Grid) int {
	G := MakeGraph(g)
	for i := range G.N {
		for _, e := range G.Edges(i) {
			G.SetWeight(G.To(e), G.From(e), G.Weight(e))
		}
	}
	return LongestPath(G)
}

type Graph struct {
	*graph.Dense[int]
	Nodes []grid.Pos
}

// MakeGraph creates a dense graph from g, where nodes are crossings (as well
// as the start and finish node) and edge-weights are the length of path
// segments between crossings.
//
// The nodes are sorted by row first and column second.
func MakeGraph(g *Grid) *Graph {
	var nodes []grid.Pos
	idx := make(map[grid.Pos]int)
	add := func(p grid.Pos) {
		idx[p] = len(nodes)
		nodes = append(nodes, p)
	}
	add(grid.Pos{0, 1})

	// find all crossings
	for p, c := range g.Cells {
		if c == Forest {
			continue
		}
		var n int
		for _, q := range g.Neigh4(p) {
			if g.At(q) == Forest {
				continue
			}
			n++
			if n > 2 {
				add(p)
				break
			}
		}
	}
	add(grid.Pos{g.H - 1, g.W - 2})

	G := &Graph{
		Dense: graph.NewDense[int](len(nodes)),
		Nodes: nodes,
	}

	type node struct {
		p     grid.Pos // current cell
		prev  grid.Pos // node this one was reached from
		cross int      // last crossing seen
		n     int      // steps since the last crossing
	}
	var q container.FIFO[node]
	for i, n := range nodes {
		q.Push(node{n, n, i, 0})
	}
	for q.Len() > 0 {
		n := q.Pop()
		if i, ok := idx[n.p]; ok && n.p != G.Nodes[n.cross] {
			G.SetWeight(n.cross, i, n.n)
			continue
		}
		for c, δ := range Δ {
			if p := n.p.Add(δ); g.Valid(p) && p != n.prev && (g.At(p) == Path || g.At(p) == c) {
				q.Push(node{p, n.p, n.cross, n.n + 1})
			}
		}
	}
	return G
}

var Δ = map[Cell]grid.Pos{
	SlopeU: grid.Pos{-1, 0},
	SlopeR: grid.Pos{0, 1},
	SlopeD: grid.Pos{1, 0},
	SlopeL: grid.Pos{0, -1},
}

type Edge struct {
	From grid.Pos
	To   grid.Pos
	N    int
}

func LongestPath(g *Graph) int {
	type node struct {
		l *list[int]
		n int
	}
	var (
		q       container.LIFO[node]
		longest int = math.MinInt
		to          = g.N - 1
	)
	q.Push(node{Push(nil, 0), 0})
	for q.Len() > 0 {
		n := q.Pop()
		if n.l.head == to {
			longest = max(longest, n.n)
			continue
		}
		for e := range g.EdgeSeq(n.l.head) {
			if n.l.Contains(g.To(e)) {
				continue
			}
			q.Push(node{Push(n.l, g.To(e)), n.n + g.Weight(e)})
		}
	}
	return longest
}

type list[E comparable] struct {
	head E
	tail *list[E]
}

func Push[E comparable](l *list[E], e E) *list[E] {
	return &list[E]{e, l}
}

func (l *list[E]) Contains(e E) bool {
	return l != nil && (l.head == e || l.tail.Contains(e))
}
