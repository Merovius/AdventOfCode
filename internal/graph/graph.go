package graph

import (
	"fmt"
	"iter"
	"sync"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/set"
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

// Dense is a weighted graph represented by an adjacency matrix.
// It implements Weighted[int, [2]int, W].
type Dense[W constraints.Integer | constraints.Float] struct {
	N int // Number of nodes
	W []W // Weight of edge i->j at i•N+j
}

// NewDense creates a Dense graph with n nodes.
func NewDense[W constraints.Integer | constraints.Float](n int) *Dense[W] {
	g := &Dense[W]{
		N: n,
		W: make([]W, n*n),
	}
	var _ Weighted[int, [2]int, W] = g
	return g
}

// Edges returns the non-zero edges adjacent to i.
func (g *Dense[W]) Edges(i int) [][2]int {
	e := make([][2]int, 0, g.N)
	for j, w := range g.W[g.N*i : g.N*i+g.N] {
		if w != 0 {
			e = append(e, [2]int{i, j})
		}
	}
	return e
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

// From returns e[0].
func (g *Dense[W]) From(e [2]int) int {
	return e[0]
}

// From returns e[1].
func (g *Dense[W]) To(e [2]int) int {
	return e[1]
}

// Weight returns the weight of e.
func (g *Dense[W]) Weight(e [2]int) W {
	return g.W[g.N*e[0]+e[1]]
}

// SetWeight sets the weight of i->j to w.
func (g *Dense[W]) SetWeight(i, j int, w W) {
	g.W[g.N*i+j] = w
}

// Sparse is a weighted graph represented by an adjacency list.
// It implements Weighted[int, [2]int, W].
type Sparse[W constraints.Integer | constraints.Float] struct {
	N int

	edges [][]sparseEdge[W]
}

func NewSparse[W constraints.Integer | constraints.Float](n int) *Sparse[W] {
	g := &Sparse[W]{n, make([][]sparseEdge[W], n)}
	var _ Weighted[int, [2]int, W] = g
	return g
}

type sparseEdge[W any] struct {
	i int
	w W
}

// Edges returns the edges adjacent to i. The returned slice must not be
// modified.
func (g *Sparse[W]) Edges(i int) [][2]int {
	out := make([][2]int, len(g.edges[i]))
	for j, e := range g.edges[i] {
		out[j] = [2]int{i, e.i}
	}
	return out
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

// From returns e[0].
func (g *Sparse[W]) From(e [2]int) int {
	return e[0]
}

// To returns e[1].
func (g *Sparse[W]) To(e [2]int) int {
	return e[1]
}

// Weight returns the weight of e.
func (g *Sparse[W]) Weight(e [2]int) W {
	for _, f := range g.edges[e[0]] {
		if f.i == e[1] {
			return f.w
		}
	}
	panic(fmt.Errorf("no edge from %d to %d", e[0], e[1]))
}

// SetWeight sets the weight of i->j to w.
func (g *Sparse[W]) SetWeight(i, j int, w W) {
	for k, e := range g.edges[i] {
		if e.i == j {
			g.edges[i][k].w = w
			return
		}
	}
	g.edges[i] = append(g.edges[i], sparseEdge[W]{j, w})
}
