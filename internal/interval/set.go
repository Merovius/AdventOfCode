package interval

import (
	"fmt"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/math"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Set is a set of T, built out of intervals. The zero value is an empty set.
type Set[I Interval[T], T constraints.Signed] struct {
	set set[T]
}

// Contains returns whether s contains v.
func (s *Set[I, T]) Contains(v T) bool {
	return s.set.Contains(v)
}

// Len returns the number of elements in s.
func (s *Set[I, T]) Len() int {
	return s.set.Len()
}

// Add adds the interval i to s.
func (s *Set[I, T]) Add(i I) {
	s.set.Add(i.toCO())
}

// Intersect intersects i into s.
func (s *Set[I, T]) Intersect(i I) {
	s.set.Intersect(i.toCO())
}

// Continuous returns whether s is representable by a single interval. The
// empty set is considered continuous.
func (s *Set[I, T]) Continuous() bool {
	return s.set.Continuous()
}

// Intervals returns a copy of the intervals in s.
func (s *Set[I, T]) Intervals() []I {
	out := make([]I, 0, len(s.set.s))
	for _, i := range s.set.s {
		out = append(out, Convert[I](i))
	}
	return out
}

// String implements fmt.Stringer.
func (s Set[I, T]) String() string {
	var (
		parts = make([]string, 0, len(s.set.s))
	)
	for _, i := range s.set.s {
		parts = append(parts, any(Convert[I](i)).(fmt.Stringer).String())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// set is a set of T, built out of intervals.
type set[T constraints.Signed] struct {
	// the slice is in canonical form, that is all intervals are not empty and
	// s[i+1].Min > s[i].Max
	s []CO[T]
}

// Contains returns whether s contains v.
func (s *set[T]) Contains(v T) bool {
	n, ok := slices.BinarySearchFunc(s.s, v, func(i CO[T], v T) int {
		return math.Cmp(i.Min, v)
	})
	if ok {
		return true
	}
	if n == 0 {
		return false
	}
	return v < s.s[n-1].Max
}

// Len returns the number of elements in s.
func (s *set[T]) Len() int {
	var n int
	for _, i := range s.s {
		n += i.Len()
	}
	return n
}

// Add adds the interval i to s.
func (s *set[T]) Add(i CO[T]) {
	if i.Empty() {
		return
	}

	lo, lok := slices.BinarySearchFunc(s.s, i.Min, func(j CO[T], v T) int {
		if v > j.Max {
			return -1
		}
		if v+1 < j.Min {
			return 1
		}
		return 0
	})
	hi, hok := slices.BinarySearchFunc(s.s, i.Max, func(j CO[T], v T) int {
		if v > j.Max {
			return -1
		}
		if v+1 < j.Min {
			return 1
		}
		return 0
	})
	if lok {
		i.Min = math.Min(i.Min, s.s[lo].Min)
	}
	if hok {
		i.Max = math.Max(i.Max, s.s[hi].Max)
		hi++
	}
	if hi < lo {
		panic("hi < lo")
	}
	if hi == lo {
		s.s = append(s.s, CO[T]{})
		copy(s.s[hi+1:], s.s[hi:])
		hi++
	}
	s.s = append(append(s.s[:lo], i), s.s[hi:]...)
}

// Intersect intersects i into s.
func (s *set[T]) Intersect(i CO[T]) {
	lo, lok := slices.BinarySearchFunc(s.s, i.Min, func(j CO[T], v T) int {
		if v >= j.Max {
			return -1
		}
		if v < j.Min {
			return 1
		}
		return 0
	})
	if lok {
		s.s[lo].Min = math.Max(s.s[lo].Min, i.Min)
		if s.s[lo].Empty() {
			lo++
		}
	}
	s.s = s.s[lo:]
	hi, hok := slices.BinarySearchFunc(s.s, i.Max, func(j CO[T], v T) int {
		if v >= j.Max {
			return -1
		}
		if v < j.Min {
			return 1
		}
		return 0
	})
	if hok {
		s.s[hi].Max = math.Min(s.s[hi].Max, i.Max)
		if !s.s[hi].Empty() {
			hi++
		}
	}
	s.s = s.s[:hi]
}

// Continuous returns whether s is representable by a single interval. The
// empty set is considered continuous.
func (s *set[T]) Continuous() bool {
	return len(s.s) <= 1
}

// Intervals returns a copy of the intervals in s.
func (s *set[T]) Intervals() []CO[T] {
	return slices.Clone(s.s)
}

func (s *set[T]) check() {
	if len(s.s) == 0 {
		return
	}
	i := s.s[0]
	if !i.valid() {
		panic(fmt.Errorf("invalid interval %v in set", i))
	}
	for _, j := range s.s[1:] {
		if !j.valid() {
			panic(fmt.Errorf("invalid interval %v in set", i))
		}
		if j.Min <= i.Max || j.Max <= i.Min {
			panic(fmt.Errorf("invalid successive pair %v and %v in set", i, j))
		}
		i = j
	}
}
