package op

import (
	"cmp"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Complex | constraints.Float | constraints.Integer
}

func Add[T Number | ~string](a, b T) T {
	return a + b
}

func Sub[T Number](a, b T) T {
	return a - b
}

func Mul[T Number](a, b T) T {
	return a * b
}

func Div[T Number](a, b T) T {
	return a / b
}

func Mod[T constraints.Integer](a, b T) T {
	return a % b
}

func Max[T cmp.Ordered](a, b T) T {
	return max(a, b)
}

func Min[T cmp.Ordered](a, b T) T {
	return min(a, b)
}
