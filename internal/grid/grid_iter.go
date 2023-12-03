//go:build goexperiment.rangefunc

package grid

import (
	"github.com/Merovius/AdventOfCode/internal/iter"
)

func (g *Grid[T]) Cells(yield func(Pos, T) bool) {
	g.Rect(g.Bounds())(yield)
}

func (g *Grid[T]) Rect(r Rectangle) iter.Seq2[Pos, T] {
	return iter.Lift(r.Intersect(g.Bounds()).All, g.At)
}

func (g *Grid[T]) Find(f func(Pos, T) bool) iter.Seq2[Pos, T] {
	return iter.Filter2(g.Cells, f)
}

func Find[T comparable](g *Grid[T], v T) iter.Seq[Pos] {
	return iter.Left(g.Find(func(p Pos, w T) bool {
		return w == v
	}))
}

func (r Rectangle) All(yield func(Pos) bool) {
	for p := r.Min; p.Row < r.Max.Row; p.Row++ {
		for p.Col = r.Min.Col; p.Col < r.Max.Col; p.Col++ {
			if !yield(p) {
				return
			}
		}
	}
}
