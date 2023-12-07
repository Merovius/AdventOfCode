// Requires Go 1.22 (or gotip) and GOEXPERIMENT=rangefunc to be set.
//go:build goexperiment.rangefunc

package iter

type Seq[A any] func(yield func(A) bool)
type Seq2[A, B any] func(yield func(A, B) bool)

func Filter[A any](s Seq[A], f func(A) bool) Seq[A] {
	return func(yield func(A) bool) {
		for a := range s {
			if f(a) && !yield(a) {
				return
			}
		}
	}
}

func Filter2[A, B any](s Seq2[A, B], f func(A, B) bool) Seq2[A, B] {
	return func(yield func(A, B) bool) {
		for a, b := range s {
			if f(a, b) && !yield(a, b) {
				return
			}
		}
	}
}

func Left[A, B any](s Seq2[A, B]) Seq[A] {
	return Project(s, func(a A, _ B) A { return a })
}

func Right[A, B any](s Seq2[A, B]) Seq[B] {
	return Project(s, func(_ A, b B) B { return b })
}

func Lift[A, B any](s Seq[A], f func(A) B) Seq2[A, B] {
	return func(yield func(A, B) bool) {
		for a := range s {
			if !yield(a, f(a)) {
				return
			}
		}
	}
}

func Project[A, B, C any](s Seq2[A, B], f func(A, B) C) Seq[C] {
	return func(yield func(C) bool) {
		for a, b := range s {
			if !yield(f(a, b)) {
				return
			}
		}
	}
}

func Enumerate[A any](s Seq[A]) Seq2[int, A] {
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

func Len[A any](s Seq[A]) int {
	i := 0
	for _ = range s {
		i++
	}
	return i
}

func FoldR[A, B any](f func(A, B) B, s Seq[A], z B) B {
	for a := range s {
		z = f(a, z)
	}
	return z
}

func FoldL[A, B any](f func(B, A) B, s Seq[A], z B) B {
	for a := range s {
		z = f(z, a)
	}
	return z
}
