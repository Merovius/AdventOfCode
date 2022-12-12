package graph

import (
	"github.com/Merovius/AdventOfCode/internal/container"
	"golang.org/x/exp/constraints"
)

type Edge[Node any] struct {
	From Node
	To   Node
}

type Graph[Node, Edge any] interface {
	Edges(Node) []Edge
	From(Edge) Node
	To(Edge) Node
}

type Weight interface {
	constraints.Integer | constraints.Float
}

type Weighted[Node, Edge, Weight any] interface {
	Graph[Node, Edge]
	Weight(Edge) Weight
}

// ShortestPath calculates the shortest path from start to a node satisfying
// goal. It uses Breadth-First-Search.
func ShortestPath[N comparable, E any](g Graph[N, E], start N, goal func(N) bool) []E {
	var (
		q     container.FIFO[E]
		prev  = make(map[N]E)
		found bool
		end   N
	)
	for _, ne := range g.Edges(start) {
		q.Push(ne)
	}
	for q.Len() > 0 {
		edge := q.Pop()
		to := g.To(edge)
		if _, ok := prev[to]; ok {
			continue
		}
		prev[to] = edge
		if goal(to) {
			found, end = true, to
			break
		}
		for _, e := range g.Edges(to) {
			q.Push(e)
		}
	}
	if !found {
		return nil
	}
	var out []E
	for end != start {
		e := prev[end]
		out = append(out, e)
		end = g.From(e)
	}
	reverse(out)
	return out
}

func reverse[T any](s []T) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - 1 - i
		s[i], s[j] = s[j], s[i]
	}
}
