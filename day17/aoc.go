package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	dim := flag.Int("dim", 3, "Number of dimensions the grid should use")
	flag.Parse()

	grid, err := ReadInput(os.Stdin, *dim)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		grid.Simulate()
	}
	fmt.Println("Number of alive cells:", grid.AliveCells())
}

func ReadInput(r io.Reader, dim int) (*Grid, error) {
	s := bufio.NewScanner(r)
	y := 0
	grid := NewGrid(dim)
	for s.Scan() {
		for x, b := range s.Bytes() {
			if b == '.' {
				continue
			}
			c := make([]int, dim)
			c[0], c[1] = x, y
			grid.SetAlive(c)
		}
		y++
	}
	return grid, s.Err()
}

type Cell []int

// Grid is a grid of a fixed dimension.
type Grid struct {
	// m is a set, with fmt.Sprintf("%v",c) as a key, where c is a []int.
	m map[string]bool
	// dim is the dimension of the grid
	dim int
	min []int
	max []int
}

// NewGrid returns a new Grid of the given dimension.
func NewGrid(dim int) *Grid {
	return &Grid{
		m:   make(map[string]bool),
		dim: dim,
	}
}

func (g *Grid) SetAlive(c Cell) {
	g.m[fmt.Sprint(c)] = true
}

func (g *Grid) cellFromKey(k string) Cell {
	k = k[1 : len(k)-1]
	sp := strings.Split(k, ",")
	c := make([]int, g.dim)
	for i, s := range sp {
		v, _ := strconv.Atoi(s)
		c[i] = v
	}
	return c
}

func (g *Grid) Simulate() {
	next := make(map[string]bool)
	for k := range g.m {
		c := g.cellFromKey(k)
		n := g.CountAliveNeighbors(c)
		if n == 2 || n == 3 {
			next[k] = true
		}
		g.WalkNeighbors(c, func(c Cell) {
			if g.m[fmt.Sprint(c)] {
				return
			}
			if n := g.CountAliveNeighbors(c); n == 3 {
				next[fmt.Sprint(c)] = true
			}
		})
	}
	g.m = next
}

func (g *Grid) WalkNeighbors(c Cell, f func(Cell)) {
	idx := make([]int, g.dim)
	for i := range idx {
		idx[i] = -1
	}

idxLoop:
	for {
		if isZero(idx) {
			continue
		}
		c2 := make(Cell, g.dim)
		copy(c2, c)
		for i := range c {
			c2[i] += idx[i]
		}
		f(c2)
		for i := 0; i < len(idx); i++ {
			idx[i]++
			if idx[i] <= 1 {
				continue idxLoop
			}
			idx[i] = -1
		}
		break
	}
}

func (g *Grid) CountAliveNeighbors(c Cell) int {
	var n int
	g.WalkNeighbors(c, func(c Cell) {
		if g.m[fmt.Sprint(c)] {
			n++
		}
	})
	return n
}

func isZero(vs []int) bool {
	for _, v := range vs {
		if v != 0 {
			return false
		}
	}
	return true
}

func (g *Grid) Bounds() (min, max Cell) {
	min, max = make(Cell, g.dim), make(Cell, g.dim)
	for i := range min {
		min[i], max[i] = math.MaxInt64, math.MinInt64
	}
	for k := range g.m {
		c := g.cellFromKey(k)
		for i := range c {
			if c[i] < min[i] {
				min[i] = c[i]
			}
			if c[i] > max[i] {
				max[i] = c[i]
			}
		}
	}
	return min, max
}

func (g *Grid) AliveCells() int {
	return len(g.m)
}
