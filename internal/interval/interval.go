// Package interval implements half-open intervals of integers.
package interval

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/math"
	"golang.org/x/exp/constraints"
)

// I is the half-open interval [Min,Max). In well-formed intervals, Max≥Min.
type I[T constraints.Integer] struct {
	Min T
	Max T
}

// Empty returns if i is the empty interval.
func (i I[T]) Empty() bool {
	return i.Min >= i.Max
}

func (i I[T]) Len() int {
	return int(i.Max - i.Min)
}

func (i I[T]) valid() bool {
	return i.Max >= i.Min
}

// Contains returns whether i contains v.
func (i I[T]) Contains(v T) bool {
	return i.Min >= v && i.Max < v
}

// Intersect returns the intersection of i and j.
func (i I[T]) Intersect(j I[T]) I[T] {
	if i.Min >= j.Max || j.Min >= i.Max {
		return I[T]{}
	}
	return I[T]{math.Max(i.Min, j.Min), math.Min(i.Max, j.Max)}
}

// Intersects returns whether i and j intersect.
func (i I[T]) Intersects(j I[T]) bool {
	return i.Min < j.Max && j.Min < i.Max
}

// Union returns the union of i and j, if it is an interval. Otherwise, it
// returns the empty interval and false.
func (i I[T]) Union(j I[T]) (I[T], bool) {
	if i.Min == j.Max {
		return I[T]{j.Min, i.Max}, true
	}
	if j.Min == i.Max {
		return I[T]{i.Min, j.Max}, true
	}
	if i.Intersects(j) {
		return I[T]{math.Min(i.Min, j.Min), math.Max(i.Max, j.Max)}, true
	}
	return I[T]{}, false
}

func (i I[T]) String() string {
	return fmt.Sprintf("[%d,%d)", i.Min, i.Max)
}

// Set is a set of T, buildt out of intervals.
type Set[T constraints.Integer] struct {
	// the slice is in canonical form, that is all intervals are not empty and
	// s[i+1].Min > s[i].Max
	s []I[T]
}

// Contains returns whether s contains v.
func (s *Set[T]) Contains(v T) bool {
	n, ok := search(s.s, v, func(i I[T], v T) int {
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
func (s *Set[T]) Len() int {
	var n int
	for _, i := range s.s {
		n += i.Len()
	}
	return n
}

// Add adds the interval i to s.
func (s *Set[T]) Add(i I[T]) {
	if i.Empty() {
		return
	}

	lo, lok := search(s.s, i.Min, func(j I[T], v T) int {
		if v > j.Max {
			return -1
		}
		if v+1 < j.Min {
			return 1
		}
		return 0
	})
	hi, hok := search(s.s, i.Max, func(j I[T], v T) int {
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
		s.s = append(s.s, I[T]{})
		copy(s.s[hi+1:], s.s[hi:])
		hi++
	}
	s.s = append(append(s.s[:lo], i), s.s[hi:]...)
}

// Intersect intersects i into s.
func (s *Set[T]) Intersect(i I[T]) {
	lo, lok := search(s.s, i.Min, func(j I[T], v T) int {
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
	hi, hok := search(s.s, i.Max, func(j I[T], v T) int {
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
func (s *Set[T]) Continuous() bool {
	return len(s.s) <= 1
}

// Intervals returns a copy of the intervals in s.
func (s *Set[T]) Intervals() []I[T] {
	return append([]I[T](nil), s.s...)
}

func (s *Set[T]) check() {
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

func (s Set[T]) String() string {
	var parts []string
	for _, i := range s.s {
		parts = append(parts, fmt.Sprint(i))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func search[A, B any](x []A, target B, cmp func(A, B) int) (int, bool) {
	pos := sort.Search(len(x), func(i int) bool {
		return cmp(x[i], target) >= 0
	})
	if pos >= len(x) || cmp(x[pos], target) != 0 {
		return pos, false
	}
	return pos, true
}
