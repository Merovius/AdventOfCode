//go:build goexperiment.rangefunc

package graph

import (
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