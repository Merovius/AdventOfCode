// Package slices adds extra helpers to the standard library slices package.
package slices

import (
	"cmp"
	"sort"
)

func SortBy[A any, B cmp.Ordered](a []A, b []B) {
	if len(a) != len(b) {
		panic("len(a) != len(b)")
	}
	sort.Sort(sortBy[A, B]{a, b})
}

type sortBy[A any, B cmp.Ordered] struct {
	a []A
	b []B
}

func (s sortBy[A, B]) Len() int {
	return len(s.a)
}

func (s sortBy[A, B]) Swap(i, j int) {
	s.a[i], s.a[j] = s.a[j], s.a[i]
	s.b[i], s.b[j] = s.b[j], s.b[i]
}

func (s sortBy[A, B]) Less(i, j int) bool {
	return s.b[i] < s.b[j]
}

func SortByFunc[A, B any](a []A, b []B, cmp func(B, B) int) {
	if len(a) != len(b) {
		panic("len(a) != len(b)")
	}
	sort.Sort(sortByFunc[A, B]{a, b, cmp})
}

type sortByFunc[A, B any] struct {
	a   []A
	b   []B
	cmp func(B, B) int
}

func (s sortByFunc[A, B]) Len() int {
	return len(s.a)
}

func (s sortByFunc[A, B]) Swap(i, j int) {
	s.a[i], s.a[j] = s.a[j], s.a[i]
	s.b[i], s.b[j] = s.b[j], s.b[i]
}

func (s sortByFunc[A, B]) Less(i, j int) bool {
	return s.cmp(s.b[i], s.b[j]) < 0
}
