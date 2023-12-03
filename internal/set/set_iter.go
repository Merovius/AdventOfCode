//go:build goexperiment.rangefunc

package set

import "github.com/Merovius/AdventOfCode/internal/iter"

func Collect[E comparable](s iter.Seq[E]) Set[E] {
	out := make(Set[E])
	for e := range s {
		out.Add(e)
	}
	return out
}

func (s Set[E]) Elements(yield func(E) bool) {
	for e := range s {
		if !yield(e) {
			return
		}
	}
}
