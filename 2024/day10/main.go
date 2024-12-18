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
	G := graph.NeighborFunc(func(p grid.Pos) iter.Seq[grid.Pos] {
		return func(yield func(grid.Pos) bool) {
			v := g.At(p)
			for q, w := range g.Neigh4(p) {
				if w-v == 1 && !yield(q) {
					return
				}
			}
		}
	})

	var score int
	var rating int
	for p, h := range g.All() {
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

// Walk the Graph in depth-first order, taking all possible paths. Differs from
// graph.WalkDepthFirst in that nodes may be visited multiple times.
func Walk(g graph.Graph[grid.Pos, [2]grid.Pos], start grid.Pos) iter.Seq[grid.Pos] {
	return func(yield func(grid.Pos) bool) {
		var q container.FIFO[grid.Pos]
		q.Push(start)
		for q.Len() > 0 {
			p := q.Pop()
			if !yield(p) {
				return
			}
			for e := range g.Edges(p) {
				q.Push(g.To(e))
			}
		}
	}
}
