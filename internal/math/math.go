// Package math provides generic arithmetic functions
//
// It intentionally does not support floating point types, as that would
// require reflect in generic code, for special handling of NaN, ±∞ and -0. Use
// the standard library math package instead.
package math

import "golang.org/x/exp/constraints"

type TotallyOrdered interface {
	constraints.Integer | ~string
}

// Max returns the maximum of a and b.
func Max[T TotallyOrdered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Max returns the minimum of a and b.
func Min[T TotallyOrdered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// Cmp returns -1, 0 and 1, if a < b, a == b and a > b, respectively.
func Cmp[T TotallyOrdered](a, b T) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// Sgn returns -1, 0, 1 if v < 0, v == 0, v > 0, respectively.
func Sgn[T constraints.Signed](v T) T {
	return T(Cmp(v, 0))
}

// Abs returns the absolute value of v.
func Abs[T constraints.Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
