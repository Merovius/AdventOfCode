package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"log"
	"math"
	"os"
	"slices"
	"strings"

	"gonih.org/AdventOfCode/internal/container"
	"gonih.org/AdventOfCode/internal/graph"
	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/set"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(in []byte) (Graph, error) {
	g, err := grid.Read(bytes.NewReader(in), func(r rune) (Cell, error) {
		if i := strings.IndexRune(".#SE", r); i >= 0 {
			return Cell(i), nil
		}
		return 0, fmt.Errorf("invalid cell %q", r)
	})
	return Graph{g}, err
}

func Part1(g Graph) int {
	start := g.g.Pos(slices.Index(g.g.G, Start))
	end := g.g.Pos(slices.Index(g.g.G, End))
	path := graph.Dijkstra(g, Node{start, grid.Right}, func(n Node) bool {
		return n.p == end
	})
	var score int
	for _, e := range path {
		score += g.Weight(e)
	}
	return score
}

func Part2(g Graph) int {
	start := g.g.Pos(slices.Index(g.g.G, Start))
	end := g.g.Pos(slices.Index(g.g.G, End))
	type el struct {
		dist int
		from Node
		to   Node
	}
	var (
		q = container.HeapFunc[el]{
			Less: func(a, b el) bool { return a.dist < b.dist },
		}
		prev = make(map[Node][]Node)
		dist = make(map[Node]int)
	)
	q.Push(el{0, Node{}, Node{start, grid.Right}})
	for q.Len() > 0 {
		e := q.Pop()
		if d, ok := dist[e.to]; ok && e.dist > d {
			continue
		} else if e.dist < d {
			prev[e.to] = prev[e.to][:0]
		} else if slices.Contains(prev[e.to], e.from) {
			continue
		}
		if e.from.p != start {
			prev[e.to] = append(prev[e.to], e.from)
		}
		dist[e.to] = e.dist
		for edge := range g.Edges(e.to) {
			q.Push(el{e.dist + g.Weight(edge), e.to, g.To(edge)})
		}
	}
	spots := make(set.Set[grid.Pos])
	spots.Add(start)
	var markAll func(n Node)
	markAll = func(n Node) {
		spots.Add(n.p)
		for _, m := range prev[n] {
			markAll(m)
		}
	}
	best := math.MaxInt
	for d := range grid.Directions() {
		if D, ok := dist[Node{end, d}]; ok {
			best = min(best, D)
		}
	}
	for d := range grid.Directions() {
		if dist[Node{end, d}] == best {
			markAll(Node{end, d})
		}
	}
	return len(spots)
}

func total(g Graph, path []Edge) int {
	var score int
	for _, e := range path {
		score += g.Weight(e)
	}
	return score
}

type Graph struct {
	g *grid.Grid[Cell]
}

type Cell byte

const (
	Empty Cell = iota
	Wall
	Start
	End
	Pos
)

func (c Cell) String() string {
	return string([]rune{' ', '▉', 'S', 'E', '•'}[c])

}

type Node struct {
	p grid.Pos
	d grid.Direction
}

type Edge = [2]Node

func dir(p, q grid.Pos) grid.Direction {
	δ := q.Sub(p)
	if δ.Row > 0 {
		return grid.Down
	} else if δ.Row < 0 {
		return grid.Up
	}
	if δ.Col > 0 {
		return grid.Right
	} else if δ.Col < 0 {
		return grid.Left
	}
	panic("invalid edge")
}

func (g Graph) Edges(n Node) iter.Seq[Edge] {
	return func(yield func(Edge) bool) {
		if g.g.At(n.p) == End {
			return
		}
		for _, d := range []grid.Direction{n.d, n.d.RotateLeft(), n.d.RotateRight()} {
			q := d.Move(n.p)
			if g.g.At(q) == Wall {
				continue
			}
			if !yield(Edge{n, Node{q, d}}) {
				return
			}
		}
	}
}

func (g Graph) From(e Edge) Node {
	return e[0]
}

func (g Graph) To(e Edge) Node {
	return e[1]
}

func (g Graph) Weight(e Edge) int {
	if e[0].d == e[1].d {
		return 1
	}
	return 1001
}
