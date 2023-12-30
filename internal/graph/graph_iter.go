//go:build goexperiment.rangefunc

package graph

import (
	"sync"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/iter"
	"github.com/Merovius/AdventOfCode/internal/set"
)

// WalkDepthFirst returns an iterator over a depth-first walk of g from start.
// Every reachable node is yielded exactly once.
func WalkDepthFirst[N comparable, E any](g Graph[N, E], start N) iter.Seq[N] {
	return func(yield func(N) bool) {
		var q container.LIFO[N]
		q.Push(start)
		seen := make(set.Set[N])
		seen.Add(start)
		for q.Len() > 0 {
			n := q.Pop()
			for _, e := range g.Edges(n) {
				m := g.To(e)
				if seen.Contains(m) {
					continue
				}
				if !yield(m) {
					return
				}
				seen.Add(m)
				q.Push(m)
			}
		}
	}
}

// EdgeSeq is like Edges, but returns an iterator.
func (g *Dense[W]) EdgeSeq(i int) iter.Seq[[2]int] {
	return func(yield func([2]int) bool) {
		for j, w := range g.W[g.N*i : g.N*i+g.N] {
			if w != 0 {
				if !yield([2]int{i, j}) {
					return
				}
			}
		}
	}
}

// WeightedEdges returns an iterator over the edges adjacent to i and their
// weights.
func (g *Sparse[W]) WeightedEdges(i int) iter.Seq2[[2]int, W] {
	return func(yield func([2]int, W) bool) {
		for _, e := range g.edges[i] {
			if !yield([2]int{i, e.i}, e.w) {
				return
			}
		}
	}
}

// MaximumFlow returns the maximum flow from source to sink, an iterator over
// one maximizing flow and the residual graph of that flow.
func MaximumFlow[N, E comparable, W Weight](g UndirectedWeighted[N, E, W], source, sink N) (W, iter.Seq2[E, W], UndirectedWeighted[N, E, W]) {
	// Edmonds-Karp algorithm: https://en.wikipedia.org/wiki/Edmonds%E2%80%93Karp_algorithm
	// TODO: Use Dinitz's Algorithm. Also use dynamic trees.
	var (
		flow  W = 0
		flows   = make(map[E]W)
	)

	for {
		var (
			q    container.FIFO[N]
			pred = make(map[N]E)
		)
		q.Push(source)
		for q.Len() > 0 {
			if _, ok := pred[sink]; ok {
				break
			}
			cur := q.Pop()
			for _, e := range g.Edges(cur) {
				to := g.To(e)
				if _, ok := pred[to]; !ok && to != source && g.Weight(e) > flows[e] {
					pred[to] = e
					q.Push(to)
				}
			}
		}
		e, ok := pred[sink]
		if !ok {
			break
		}
		df := g.Weight(e) - flows[e]
		for e, ok := pred[sink]; ok; e, ok = pred[g.From(e)] {
			df = min(df, g.Weight(e)-flows[e])
		}
		for e, ok := pred[sink]; ok; e, ok = pred[g.From(e)] {
			flows[e] += df
			flows[g.Reverse(e)] -= df
		}
		flow += df
	}
	edges := func(yield func(E, W) bool) {
		for e, w := range flows {
			if w <= 0 {
				continue
			}
			if !yield(e, w) {
				return
			}
		}
	}
	res := residual[N, E, W]{g, func(e E) W { return flows[e] }}
	return flow, edges, res
}

type residual[N, E any, W Weight] struct {
	UndirectedWeighted[N, E, W]
	flow func(E) W
}

func (r residual[N, E, W]) Weight(e E) W {
	return r.UndirectedWeighted.Weight(e) - r.flow(e)
}

// MinimumCut returns a minimal cut separating a and b, an iterator over the
// edges to cut and the nodes of the sub graphs containing a.
func MinimumCut[N, E comparable, W Weight](g UndirectedWeighted[N, E, W], a, b N) (cut W, edges iter.Seq[E], reachable iter.Seq[N]) {
	flow, edges2, res := MaximumFlow(g, a, b)

	reachable = func(yield func(N) bool) {
		var q container.LIFO[N]
		seen := make(set.Set[N])

		push := func(n N) bool {
			if seen.Contains(n) {
				return true
			}
			if !yield(n) {
				return false
			}
			seen.Add(n)
			q.Push(n)
			return true
		}
		push(a)
		for q.Len() > 0 {
			n := q.Pop()
			for _, e := range res.Edges(n) {
				if res.Weight(e) == 0 {
					continue
				}
				if !push(res.To(e)) {
					return
				}
			}
		}
	}
	toSet := sync.OnceValue(func() set.Set[N] { return set.Collect(reachable) })

	edges = func(yield func(E) bool) {
		r := toSet()
		for e, w := range edges2 {
			if g.Weight(e) < w {
				continue
			}
			if r.Contains(g.From(e)) != r.Contains(g.To(e)) {
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
	return flow, edges, reachable
}

// Undirected graph. Every edge must be reversible.
type Undirected[N, E any] interface {
	Graph[N, E]
	Reverse(E) E
}

// Undirected weighted graph. It must be Weight(Revers(e)) == Weight(e) for all edges.
type UndirectedWeighted[N, E, W any] interface {
	Undirected[N, E]
	Weight(E) W
}
