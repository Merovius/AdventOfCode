package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWalkNeighbors(t *testing.T) {
	for dim := 2; dim < 4; dim++ {
		g := NewGrid(dim)
		var got []Cell
		zero := make(Cell, dim)
		g.WalkNeighbors(zero, func(c Cell) {
			got = append(got, c)
		})
		sort.Slice(got, func(i, j int) bool {
			return Cell.less(got[i], got[j])
		})
		t.Logf("WalkNeighbors[dim=%d] returned %v", dim, got)
		for i := 1; i < len(got); i++ {
			if Cell.cmp(got[i-1], got[i]) == 0 {
				t.Errorf("WalkNeighbors[dim=%d] returned duplicate neighbor %v", dim, got[i])
			}
		}
		for _, c := range got {
			if Cell.dist(zero, c) != 1 {
				t.Errorf("WalkNeighbors[dim=%d] returned non-neighbor %v", dim, c)
			}
		}

		if want := int(math.Pow(3, float64(dim))) - 1; len(got) != want {
			t.Errorf("WalkNeighbors[dim=%d] returned %d neighbors, want %d", dim, len(got), want)
		}
	}
}

func TestCountAliveNeighbors(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	for i := 0; i < 100; i++ {
		dim := rnd.Intn(3) + 2
		g := NewGrid(dim)
		want := rnd.Intn(int(math.Pow(3, float64(dim))))
		offs := make(Cell, dim)
		for i := range offs {
			offs[i] += rnd.Intn(11) - 5
		}
		n := 0
		for n < want {
			c := make(Cell, dim)
			for i := range c {
				c[i] = rnd.Intn(3) - 1 + offs[i]
			}
			if !(Cell.eq(offs, c) || g.Alive(c)) {
				g.SetAlive(c)
				n++
			}
		}
		got := g.CountAliveNeighbors(offs)
		if got != want {
			t.Errorf("CountAliveNeighbors(%v) = %d, want %d", g, got, want)
		}
	}
}

func TestCellFromKey(t *testing.T) {
	tcs := []struct {
		dim  int
		key  string
		want Cell
	}{
		{3, "[1 2 3]", Cell{1, 2, 3}},
		{4, "[1 2 3 4]", Cell{1, 2, 3, 4}},
	}
	for _, tc := range tcs {
		g := NewGrid(tc.dim)
		if got := g.cellFromKey(tc.key); !cmp.Equal(got, tc.want) {
			t.Errorf("Grid[dim=%d].cellFromKey(%q) = %v, want %v", tc.dim, tc.key, got, tc.want)
		}
	}
}

func TestKeyForCell(t *testing.T) {
	tcs := []struct {
		cell Cell
		want string
	}{
		{Cell{1, 2, 3}, "[1 2 3]"},
		{Cell{1, 2, 3, 4}, "[1 2 3 4]"},
	}
	for _, tc := range tcs {
		if got := fmt.Sprint(tc.cell); got != tc.want {
			t.Errorf("fmt.Sprintf(%v) = %q, want %q", tc.cell, got, tc.want)
		}
	}
}
