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
	// m is a set, with fmt.Sprint(c) as a key, where c is a Cell.
	m map[string]bool
	// dim is the dimension of the grid
	dim int
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

func (g *Grid) Alive(c Cell) bool {
	return g.m[fmt.Sprint(c)]
}

func (g *Grid) cellFromKey(k string) Cell {
	k = k[1 : len(k)-1]
	sp := strings.Split(k, " ")
	c := make([]int, g.dim)
	for i, s := range sp {
		v, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Errorf("invalid key %q: %w", k, err))
		}
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
			if g.Alive(c) {
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
		var c2 Cell
		if isZero(idx) {
			goto inc
		}
		c2 = make(Cell, g.dim)
		copy(c2, c)
		for i := range c {
			c2[i] += idx[i]
		}
		f(c2)
	inc:
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
		if g.Alive(c) {
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
		fmt.Println(c)
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

func (g *Grid) Dump() {
	min, max := g.Bounds()
	c := append(Cell(nil), min...)
planeLoop:
	for {
		if g.dim > 2 {
			// Assemble the plane-identifier ("x2=23, x3=42:")
			var plane []string
			for i := 2; i < len(c); i++ {
				plane = append(plane, fmt.Sprintf("x%d=%d", i, c[i]))
			}
			fmt.Printf("%s:\n", strings.Join(plane, ", "))
		}
		for c[1] = min[1]; c[1] <= max[1]; c[1]++ {
			for c[0] = min[0]; c[0] <= max[0]; c[0]++ {
				if g.Alive(c) {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}

		for i := 2; i < len(c); i++ {
			c[i]++
			if c[i] <= max[i] {
				continue planeLoop
			}
			c[i] = min[i]
		}
		break
	}
}

func (c Cell) cmp(c2 Cell) int {
	if len(c) != len(c2) {
		panic("comparing cells of different dimension")
	}
	for i := range c {
		if c[i] < c2[i] {
			return -1
		} else if c[i] > c2[i] {
			return 1
		}
	}
	return 0
}

func (c Cell) less(c2 Cell) bool {
	return c.cmp(c2) < 0
}

// dist returns the distance of c and c2 in the maximum norm.
func (c Cell) dist(c2 Cell) int {
	if len(c) != len(c2) {
		panic("comparing cells of different dimension")
	}
	var d int
	for i := range c {
		δ := c[i] - c2[i]
		if δ < 0 {
			δ = -δ
		}
		if δ > d {
			d = δ
		}
	}
	return d
}

func (c Cell) eq(c2 Cell) bool {
	return c.cmp(c2) == 0
}

func (c Cell) zero() bool {
	for _, v := range c {
		if v != 0 {
			return false
		}
	}
	return true
}
