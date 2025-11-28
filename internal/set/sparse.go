package set

import (
	"iter"
	"slices"
)

// Sparse represents a set of integers in the range {0,…,N-1} for some bound N.
type Sparse struct {
	i []int
	e []int
}

func NewSparse(n int) *Sparse {
	a := make([]int, 2*n)
	return &Sparse{
		i: a[0:n],
		e: a[n:n],
	}
}

// Bound returns the (exclusively) upper bound for elements in s.
func (s *Sparse) Bound() int {
	return len(s.i)
}

// Len returns the number of elements in s.
func (s *Sparse) Len() int {
	return len(s.e)
}

// Add e to s. It is O(1).
func (s *Sparse) Add(e int) {
	if i := s.i[e]; i < len(s.e) && s.e[i] == e {
		return
	}
	s.i[e] = len(s.e)
	s.e = append(s.e, e)
}

// Contains returns whether e is in s. It is O(1).
func (s *Sparse) Contains(e int) bool {
	return s.i[e] < len(s.e) && s.e[s.i[e]] == e
}

// Clear s. Clear is O(1).
func (s *Sparse) Clear() {
	s.e = s.e[:0]
}

// All yields the elements of s in insertion order. It is O(s.Len()).
func (s *Sparse) All() iter.Seq[int] {
	return slices.Values(s.e)
}

// Sorted yields the elements of s in ascending order. It is O(s.Bound()).
func (s *Sparse) Sorted() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, j := range s.i {
			if j < len(s.e) && s.e[j] == i && !yield(i) {
				return
			}
		}
	}
}

// Sorted yields the elements of s in ascending order. It is O(s.Bound()).
func (s *Sparse) Descending() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := len(s.i) - 1; i >= 0; i-- {
			j := s.i[i]
			if j < len(s.e) && s.e[j] == i && !yield(i) {
				return
			}
		}
	}
}
