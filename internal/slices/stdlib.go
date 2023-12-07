package slices

import (
	"cmp"
	"slices"
)

func BinarySearch[S ~[]E, E cmp.Ordered](x S, target E) (int, bool) {
	return slices.BinarySearch(x, target)
}

func BinarySearchFunc[S ~[]E, E, T any](x S, target T, cmp func(E, T) int) (int, bool) {
	return slices.BinarySearchFunc(x, target, cmp)
}

func Clip[S ~[]E, E any](s S) S {
	return slices.Clip(s)
}

func Clone[S ~[]E, E any](s S) S {
	return slices.Clone(s)
}

func Compact[S ~[]E, E comparable](s S) S {
	return slices.Compact(s)
}

func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	return slices.CompactFunc(s, eq)
}

func Compare[S ~[]E, E cmp.Ordered](s1, s2 S) int {
	return slices.Compare(s1, s2)
}

func CompareFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, cmp func(E1, E2) int) int {
	return slices.CompareFunc(s1, s2, cmp)
}

func Contains[S ~[]E, E comparable](s S, v E) bool {
	return slices.Contains(s, v)
}

func ContainsFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	return slices.ContainsFunc(s, f)
}

func Delete[S ~[]E, E any](s S, i, j int) S {
	return slices.Delete(s, i, j)
}

func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	return slices.DeleteFunc(s, del)
}

func Equal[S ~[]E, E comparable](s1, s2 S) bool {
	return slices.Equal(s1, s2)
}

func EqualFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, eq func(E1, E2) bool) bool {
	return slices.EqualFunc(s1, s2, eq)
}

func Grow[S ~[]E, E any](s S, n int) S {
	return slices.Grow(s, n)
}

func Index[S ~[]E, E comparable](s S, v E) int {
	return slices.Index(s, v)
}

func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	return slices.IndexFunc(s, f)
}

func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	return slices.Insert(s, i, v...)
}

func IsSorted[S ~[]E, E cmp.Ordered](x S) bool {
	return slices.IsSorted(x)
}

func IsSortedFunc[S ~[]E, E any](x S, cmp func(a, b E) int) bool {
	return slices.IsSortedFunc(x, cmp)
}

func Max[S ~[]E, E cmp.Ordered](x S) E {
	return slices.Max(x)
}

func MaxFunc[S ~[]E, E any](x S, cmp func(a, b E) int) E {
	return slices.MaxFunc(x, cmp)
}

func Min[S ~[]E, E cmp.Ordered](x S) E {
	return slices.Min(x)
}

func MinFunc[S ~[]E, E any](x S, cmp func(a, b E) int) E {
	return slices.MinFunc(x, cmp)
}

func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	return slices.Replace(s, i, j, v...)
}

func Reverse[S ~[]E, E any](s S) {
	slices.Reverse(s)
}

func Sort[S ~[]E, E cmp.Ordered](x S) {
	slices.Sort(x)
}

func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	slices.SortFunc(x, cmp)
}

func SortStableFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	slices.SortStableFunc(x, cmp)
}
