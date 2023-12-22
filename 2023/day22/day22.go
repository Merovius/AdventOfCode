//go:build goexperiment.rangefunc

package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
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
	return parse.Lines(
		parse.Array[Brick](
			split.On("~"),
			parse.Array[Vec3](
				split.On(","),
				parse.Signed[int],
			),
		),
	)(s)
}

func Part1(in []Brick) int {
	supports, supported := Fall(in)

	var total int
	for _, js := range supports {
		can := true
		for j := range js {
			if len(supported[j]) < 2 {
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
	supports, supported := Fall(in)
	var total int
	for i := range in {
		var q = container.HeapFunc[int]{
			Less: func(a, b int) bool {
				return in[a].Min(Z) < in[b].Min(Z)
			},
		}
		seen := make(set.Set[int])
		gone := make(set.Set[int])
		push := func(i int) {
			if seen.Contains(i) {
				return
			}
			seen.Add(i)
			q.Push(i)
		}
		push(i)
		gone.Add(i)
		for j := range supports[i] {
			push(j)
		}
	loop:
		for q.Len() > 0 {
			i := q.Pop()
			for j := range supported[i] {
				if !gone.Contains(j) {
					continue loop
				}
			}
			gone.Add(i)
			for j := range supports[i] {
				push(j)
			}
		}
		// the initial node is gone, but has not "fallen"
		total += len(gone) - 1
	}
	return total
}

func Fall(in []Brick) (supports, supported []set.Set[int]) {
	slices.SortFunc(in, func(a, b Brick) int {
		return cmp.Compare(a.Min(Z), b.Min(Z))
	})

	vol := make(map[Vec3]int)
	free := func(b Brick) bool {
		for c := range b.Cells {
			if c[Z] <= 0 {
				return false
			}
			if _, ok := vol[c]; ok {
				return false
			}
		}
		return true
	}
	add := func(i int, b Brick) {
		for c := range b.Cells {
			if _, ok := vol[c]; ok {
				panic("block is added to non-empty space")
			}
			vol[c] = i
		}
	}
	remove := func(i int, b Brick) {
		for c := range b.Cells {
			if vol[c] != i {
				panic("block is not where it's supposed to be")
			}
			delete(vol, c)
		}
	}

	for i, b := range in {
		add(i, b)
	}

	// move all blocks down as far as possible
	for i, b := range in {
		remove(i, b)
		for {
			b2 := b.Down(1)
			if !free(b2) {
				break
			}
			b = b2
		}
		in[i] = b
		add(i, b)
	}

	// figure out which block supports which
	supports = make([]set.Set[int], len(in))
	supported = make([]set.Set[int], len(in))
	for i := range supports {
		supports[i] = make(set.Set[int])
		supported[i] = make(set.Set[int])
	}
	for i, b := range in {
		for c := range b.Cells {
			if j, ok := vol[c.Down(1)]; ok && j != i {
				supports[j].Add(i)
				supported[i].Add(j)
			}
		}
	}
	return supports, supported
}

type Brick [2]Vec3

func (b Brick) Cells(yield func(Vec3) bool) {
	for x := min(b[0][X], b[1][X]); x <= max(b[0][X], b[1][X]); x++ {
		for y := min(b[0][Y], b[1][Y]); y <= max(b[0][Y], b[1][Y]); y++ {
			for z := min(b[0][Z], b[1][Z]); z <= max(b[0][Z], b[1][Z]); z++ {
				if !yield(Vec3{x, y, z}) {
					return
				}
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
	return min(b[0][d], b[1][d])
}

type Vec2 [2]int

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
