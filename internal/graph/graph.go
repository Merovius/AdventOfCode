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

// BreadthFirstSearch calculates the shortest path from start to a node
// satisfying goal.
func BreadthFirstSearch[N comparable, E any](g Graph[N, E], start N, goal func(N) bool) []E {
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

// Dijkstra calculates the shortest path from start to a node satisfying goal
// using Dijkstra's algorithm.
func Dijkstra[N comparable, E any, W Weight](g Weighted[N, E, W], start N, goal func(N) bool) []E {
	return AStar(g, start, goal, func(N) W { return 0 })
}

// AStar calculates the shortest path from start to a node satisfying goal
// using the A* algorithm.
func AStar[N comparable, E any, W Weight](g Weighted[N, E, W], start N, goal func(N) bool, h func(N) W) []E {
	type el struct {
		w W
		e E
	}
	var (
		q = container.HeapFunc[el]{
			Less: func(a, b el) bool { return a.w < b.w },
		}
		prev  = make(map[N]E)
		found bool
		end   N
	)
	for _, e := range g.Edges(start) {
		q.Push(el{g.Weight(e), e})
	}
	for q.Len() > 0 {
		edge := q.Pop()
		to := g.To(edge.e)
		if _, ok := prev[to]; ok {
			continue
		}
		prev[to] = edge.e
		if goal(to) {
			found, end = true, to
			break
		}
		for _, e := range g.Edges(to) {
			q.Push(el{edge.w + g.Weight(e) + h(g.To(e)), e})
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
