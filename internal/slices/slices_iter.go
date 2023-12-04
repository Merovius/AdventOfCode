//go:build goexperiment.rangefunc

package slices

import "github.com/Merovius/AdventOfCode/internal/iter"

func Elements[E any](s []E) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, e := range s {
			if !yield(i, e) {
				return
			}
		}
	}
}

func Backwards[E any](s []E) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}
