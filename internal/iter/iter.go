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
	return func(yield func(A) bool) {
		for a, _ := range s {
			if !yield(a) {
				return
			}
		}
	}
}

func Right[A, B any](s Seq2[A, B]) Seq[B] {
	return func(yield func(B) bool) {
		for _, b := range s {
			if !yield(b) {
				return
			}
		}
	}
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
