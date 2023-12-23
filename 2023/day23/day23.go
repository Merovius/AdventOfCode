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
	edges := MakeGraph(g)
	start := grid.Pos{0, 1}
	end := grid.Pos{g.H - 1, g.W - 2}

	return LongestPath(g, edges, start, end)
}

func Part2(g *Grid) int {
	edges := MakeGraph(g)
	for _, es := range edges {
		for e := range es {
			e.To, e.From = e.From, e.To
			edges[e.From].Add(e)
		}
	}
	start := grid.Pos{0, 1}
	end := grid.Pos{g.H - 1, g.W - 2}

	return LongestPath(g, edges, start, end)
}

func MakeGraph(g *Grid) map[grid.Pos]set.Set[Edge] {
	edges := make(map[grid.Pos]set.Set[Edge])

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
				edges[p] = make(set.Set[Edge])
				break
			}
		}
	}
	edges[grid.Pos{0, 1}] = make(set.Set[Edge])
	edges[grid.Pos{g.H - 1, g.W - 2}] = make(set.Set[Edge])

	type node struct {
		p     grid.Pos // current node
		prev  grid.Pos // node this one was reached from
		cross grid.Pos // last crossing seen
		n     int      // steps since the last crossing
	}
	var q container.FIFO[node]
	for n := range edges {
		q.Push(node{n, n, n, 0})
	}
	for q.Len() > 0 {
		n := q.Pop()
		if _, ok := edges[n.p]; ok && n.p != n.cross {
			edges[n.cross].Add(Edge{n.cross, n.p, n.n})
			continue
		}
		for c, δ := range Δ {
			if p := n.p.Add(δ); g.Valid(p) && p != n.prev && (g.At(p) == Path || g.At(p) == c) {
				q.Push(node{p, n.p, n.cross, n.n + 1})
			}
		}
	}
	return edges
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

func LongestPath(g *Grid, edges map[grid.Pos]set.Set[Edge], from, to grid.Pos) int {
	type node struct {
		l *list[grid.Pos]
		n int
	}
	var (
		q       container.LIFO[node]
		longest int = math.MinInt
	)
	q.Push(node{Push(nil, from), 0})
	for q.Len() > 0 {
		n := q.Pop()
		if n.l.head == to {
			longest = max(longest, n.n)
			continue
		}
		for e := range edges[n.l.head] {
			if n.l.Contains(e.To) {
				continue
			}
			q.Push(node{Push(n.l, e.To), n.n + e.N})
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
