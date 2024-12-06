package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/set"
	"github.com/Merovius/AdventOfCode/internal/slices"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(s string) ([]Brick, error) {
	var out []Brick
	for len(s) > 0 {
		l, rest, ok := strings.Cut(s, "\n")
		if l == "" {
			break
		}
		s = rest
		before, after, ok := strings.Cut(l, "~")
		if !ok {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		var (
			b   Brick
			err error
		)
		b[0], err = ParseVec(before)
		if err != nil {
			return nil, err
		}
		b[1], err = ParseVec(after)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func ParseVec(s string) (Vec3, error) {
	var vec Vec3
	for i, p := range strings.Split(s, ",") {
		if i > len(vec) {
			return vec, fmt.Errorf("invalid vec %q", s)
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return vec, err
		}
		vec[i] = v
	}
	return vec, nil
}

func Part1(in []Brick) int {
	above, below := Fall(in)

	var total int
	for _, js := range above {
		can := true
		for _, j := range js {
			if len(below[j]) < 2 {
				can = false
				break
			}
		}
		if can {
			total++
		}
	}
	return total
}

func Part2(in []Brick) int {
	above, below := Fall(in)

	var total int
	for i := range in {
		var q container.Heap[int]
		seen := make(set.Set[int])
		push := func(i int) {
			if seen.Contains(i) {
				return
			}
			seen.Add(i)
			q.Push(i)
		}
		gone := set.Make(i)
		for _, j := range above[i] {
			push(j)
		}
	loop:
		for q.Len() > 0 {
			i := q.Pop()
			for _, j := range below[i] {
				if !gone.Contains(j) {
					continue loop
				}
			}
			gone.Add(i)
			for _, j := range above[i] {
				push(j)
			}
		}
		// the initial node is gone, but has not "fallen"
		total += len(gone) - 1
	}
	return total
}

func Fall(in []Brick) (above, below [][]int) {
	slices.SortFunc(in, func(a, b Brick) int {
		return cmp.Compare(a.Min(Z), b.Min(Z))
	})

	top := MakeGrid(in, -1)
	above = make([][]int, len(in))
	below = make([][]int, len(in))
	for i, b := range in {
		var h int
		for c := range b.Base {
			if j := top.At(c); j >= 0 {
				if z := in[j].Max(Z); z > h {
					below[i] = append(below[i][:0], j)
					h = z
				} else if z == h {
					below[i] = append(below[i], j)
				}
			}
			top.Set(c, i)
		}
		slices.Sort(below[i])
		below[i] = slices.Compact(below[i])

		δ := b.Min(Z) - (h + 1)
		in[i] = b.Down(δ)
	}
	for i, s := range below {
		for _, j := range s {
			above[j] = append(above[j], i)
		}
	}
	return above, below
}

type Brick [2]Vec3

func (b Brick) Base(yield func(Vec2) bool) {
	for y := min(b[0][Y], b[1][Y]); y <= max(b[0][Y], b[1][Y]); y++ {
		for x := min(b[0][X], b[1][X]); x <= max(b[0][X], b[1][X]); x++ {
			if !yield(Vec2{x, y}) {
				return
			}
		}
	}
}

func (b Brick) Down(n int) Brick {
	b[0] = b[0].Down(n)
	b[1] = b[1].Down(n)
	return b
}

func (b Brick) Min(d Dim) int {
	return min(b[0][d], b[1][d])
}

func (b Brick) Max(d Dim) int {
	return max(b[0][d], b[1][d])
}

type Vec2 [2]int

func (v Vec2) Add(w Vec2) Vec2 {
	v[X] += w[X]
	v[Y] += w[Y]
	return v
}

type Vec3 [3]int

func (v Vec3) Down(n int) Vec3 {
	v[Z] -= n
	return v
}

type Dim uint8

const (
	X Dim = iota
	Y
	Z
)

type Grid[T any] struct {
	min    Vec2
	max    Vec2
	stride int
	b      []T
}

func MakeGrid[T any](bricks []Brick, zero T) *Grid[T] {
	g := new(Grid[T])
	for _, b := range bricks {
		g.min[X] = min(g.min[X], b[0][X], b[1][X])
		g.max[X] = max(g.max[X], b[0][X], b[1][X])
		g.min[Y] = min(g.min[Y], b[0][Y], b[1][Y])
		g.max[Y] = max(g.max[Y], b[0][Y], b[1][Y])
	}
	dx, dy := g.max[X]-g.min[X]+1, g.max[Y]-g.min[Y]+1
	g.stride = dx
	g.b = make([]T, dx*dy)
	for i := range g.b {
		g.b[i] = zero
	}
	return g
}

func (g *Grid[T]) At(v Vec2) T {
	v = v.Add(g.min)
	return g.b[g.stride*v[Y]+v[X]]
}

func (g *Grid[T]) Set(v Vec2, val T) {
	v = v.Add(g.min)
	g.b[g.stride*v[Y]+v[X]] = val
}
