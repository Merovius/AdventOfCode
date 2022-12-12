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
		q    container.FIFO[E]
		prev = make(map[N]E)
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
			return makePath(start, to, g.From, prev)
		}
		for _, e := range g.Edges(to) {
			q.Push(e)
		}
	}
	return nil
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
		prio W
		dist W
		edge E
		to   N
	}
	var (
		q = container.HeapFunc[el]{
			Less: func(a, b el) bool { return a.prio < b.prio },
		}
		prev = make(map[N]E)
		dist = make(map[N]W)
	)
	q.Push(el{0, 0, *new(E), start})
	for q.Len() > 0 {
		e := q.Pop()
		if d, ok := dist[e.to]; ok && e.dist >= d {
			continue
		}
		prev[e.to], dist[e.to] = e.edge, e.dist
		if goal(e.to) {
			return makePath(start, e.to, g.From, prev)
		}
		for _, n := range g.Edges(e.to) {
			to, w := g.To(n), g.Weight(n)
			q.Push(el{e.dist + w + h(to), e.dist + w, n, to})
		}
	}
	return nil
}

func makePath[N comparable, E any](start, end N, from func(E) N, prev map[N]E) []E {
	var out []E
	for end != start {
		e := prev[end]
		out = append(out, e)
		end = from(e)
	}
	for i := 0; i < len(out)/2; i++ {
		j := len(out) - 1 - i
		out[i], out[j] = out[j], out[i]
	}
	return out
}
