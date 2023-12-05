package slices

import "slices"

func Clone[S ~[]E, E any](s S) S {
	return slices.Clone(s)
}
