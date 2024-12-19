package xiter

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Filter[A any](s iter.Seq[A], f func(A) bool) iter.Seq[A] {
	return func(yield func(A) bool) {
		for a := range s {
			if f(a) && !yield(a) {
				return
			}
		}
	}
}

func Filter2[A, B any](s iter.Seq2[A, B], f func(A, B) bool) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		for a, b := range s {
			if f(a, b) && !yield(a, b) {
				return
			}
		}
	}
}

func Left[A, B any](s iter.Seq2[A, B]) iter.Seq[A] {
	return Project(s, func(a A, _ B) A { return a })
}

func Right[A, B any](s iter.Seq2[A, B]) iter.Seq[B] {
	return Project(s, func(_ A, b B) B { return b })
}

func Lift[A, B any](s iter.Seq[A], f func(A) B) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		for a := range s {
			if !yield(a, f(a)) {
				return
			}
		}
	}
}

func Project[A, B, C any](s iter.Seq2[A, B], f func(A, B) C) iter.Seq[C] {
	return func(yield func(C) bool) {
		for a, b := range s {
			if !yield(f(a, b)) {
				return
			}
		}
	}
}

func Enumerate[A any](s iter.Seq[A]) iter.Seq2[int, A] {
	return func(yield func(int, A) bool) {
		i := 0
		for a := range s {
			if !yield(i, a) {
				return
			}
			i++
		}
	}
}

func Len[A any](s iter.Seq[A]) int {
	i := 0
	for _ = range s {
		i++
	}
	return i
}

func FoldR[A, B any](f func(A, B) B, s iter.Seq[A], z B) B {
	for a := range s {
		z = f(a, z)
	}
	return z
}

func FoldL[A, B any](f func(B, A) B, s iter.Seq[A], z B) B {
	for a := range s {
		z = f(z, a)
	}
	return z
}

func Range[T constraints.Integer | constraints.Float](min, max, step T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := min; i < max; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

func Map[A, B any](s iter.Seq[A], f func(A) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range s {
			if !yield(f(a)) {
				return
			}
		}
	}
}
