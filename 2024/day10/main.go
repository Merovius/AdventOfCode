package main

import (
	"fmt"
	"iter"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/grid"
)

func main() {
	g, err := grid.Read(os.Stdin, func(r rune) (int8, error) {
		return int8(r - '0'), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	G := &Graph{g}

	var score int
	var rating int
	for p, h := range G.g.All() {
		if h != 0 {
			continue
		}
		for q := range graph.WalkDepthFirst(G, p) {
			if g.At(q) == 9 {
				score++
			}
		}
		for q := range Walk(G, p) {
			if g.At(q) == 9 {
				rating++
			}
		}
	}
	fmt.Println(score)
	fmt.Println(rating)
}

type Graph struct {
	g *grid.Grid[int8]
}

func (g *Graph) Edges(p grid.Pos) [][2]grid.Pos {
	var out [][2]grid.Pos
	v := g.g.At(p)
	for _, q := range g.g.Neigh4(p) {
		w := g.g.At(q)
		if w-v == 1 {
			out = append(out, [2]grid.Pos{p, q})
		}
	}
	return out
}

func (g *Graph) From(e [2]grid.Pos) grid.Pos {
	return e[0]
}

func (g *Graph) To(e [2]grid.Pos) grid.Pos {
	return e[1]
}

// Walk the Graph in depth-first order, taking all possible paths. Differs from
// graph.WalkDepthFirst in that nodes may be visited multiple times.
func Walk(g *Graph, start grid.Pos) iter.Seq[grid.Pos] {
	return func(yield func(grid.Pos) bool) {
		var q container.FIFO[grid.Pos]
		q.Push(start)
		for q.Len() > 0 {
			p := q.Pop()
			if !yield(p) {
				return
			}
			for _, e := range g.Edges(p) {
				q.Push(g.To(e))
			}
		}
	}
}
