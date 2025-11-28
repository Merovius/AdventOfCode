package set

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSparse(t *testing.T) {
	const (
		N = 128 // maximum integer to use
	)
	t.Parallel()

	s := NewSparse(N)
	for i := range N {
		if s.Contains(i) {
			t.Fatalf("s.Contains(%d) before inserting it", i)
		}
		s.Add(i)
		if !s.Contains(i) {
			t.Fatalf("!s.Contains(%d) after inserting it", i)
		}
		if n := s.Len(); n != i+1 {
			t.Fatalf("s.Len() = %d after inserting %d, want %d", n, i, i+1)
		}
	}
	t.Run("All", func(t *testing.T) {
		rnd := rand.New(rand.NewPCG(0, 0))
		for i := range 100 {
			t.Run(fmt.Sprintf("%.2d", i), func(t *testing.T) {
				var (
					s    = NewSparse(N)
					want []int
				)
				for range 10 {
					i := rnd.IntN(N)
					if !s.Contains(i) {
						s.Add(i)
						want = append(want, i)
					}
				}
				var got []int
				for i := range s.All() {
					got = append(got, i)
				}
				if d := cmp.Diff(got, want); d != "" {
					t.Errorf("All() yielded wrong sequence (-got,+want):\n%s", d)
				}
				s.Clear()
				if s.Len() != 0 {
					t.Errorf("Clear() left %d elements, want 0", s.Len())
				}
			})
		}
	})
	t.Run("Sorted", func(t *testing.T) {
		rnd := rand.New(rand.NewPCG(0, 0))
		for i := range 100 {
			t.Run(fmt.Sprintf("%.2d", i), func(t *testing.T) {
				var (
					s    = NewSparse(N)
					want []int
				)
				for range 10 {
					i := rnd.IntN(N)
					s.Add(i)
					want = append(want, i)
				}
				slices.Sort(want)
				want = slices.Compact(want)
				var got []int
				for i := range s.Sorted() {
					got = append(got, i)
				}
				if d := cmp.Diff(got, want); d != "" {
					t.Errorf("Reverse() yielded wrong sequence (-got,+want):\n%s", d)
				}
				s.Clear()
				if s.Len() != 0 {
					t.Errorf("Clear() left %d elements, want 0", s.Len())
				}
			})
		}
	})
	t.Run("Descending", func(t *testing.T) {
		rnd := rand.New(rand.NewPCG(0, 0))
		for i := range 100 {
			t.Run(fmt.Sprintf("%.2d", i), func(t *testing.T) {
				var (
					s    = NewSparse(N)
					want []int
				)
				for range 10 {
					i := rnd.IntN(N)
					s.Add(i)
					want = append(want, i)
				}
				slices.Sort(want)
				want = slices.Compact(want)
				slices.Reverse(want)
				var got []int
				for i := range s.Descending() {
					got = append(got, i)
				}
				if d := cmp.Diff(got, want); d != "" {
					t.Errorf("Reverse() yielded wrong sequence (-got,+want):\n%s", d)
				}
				s.Clear()
				if s.Len() != 0 {
					t.Errorf("Clear() left %d elements, want 0", s.Len())
				}
			})
		}
	})
}
