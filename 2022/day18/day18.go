package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input"
	"github.com/Merovius/AdventOfCode/internal/math"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	data, err := input.Slice(input.Lines(), input.Array[Node](input.Split(","), input.Signed[int]())).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g := NewGraph(data)
	c := g.Droplets()
	var (
		total  int
		total2 int
	)
	for _, d := range c {
		total += g.SurfaceArea(d)
		total2 += g.OutsideSurfaceArea(d)
	}
	fmt.Println("Total surface area of droplets:", total)
	fmt.Println("Total air surface area of droplets:", total2)

}

type Node [3]int

func (n Node) Coords() (x, y, z int) {
	return n[0], n[1], n[2]
}

type Edge struct {
	From Node
	To   Node
}

func Max(n ...Node) (x, y, z int) {
	for _, n := range n {
		x = math.Max(x, n[0])
		y = math.Max(y, n[1])
		z = math.Max(z, n[2])
	}
	return x, y, z
}

type Cell int

const (
	Air Cell = iota
	Water
	Steam
)

type Graph struct {
	grid  *Grid3[Cell]
	label *Grid3[int]
}

func NewGraph(data []Node) *Graph {
	x, y, z := Max(data...)
	x += 1
	y += 1
	z += 1
	g := &Graph{
		grid:  NewGrid3[Cell](x, y, z),
		label: NewGrid3[int](x, y, z),
	}
	g.grid.SetAll(Water, data...)
	g.setLabels()
	return g
}

func (g *Graph) setLabels() {
	var (
		visited = make(set.Set[Node])
		label   int
		visit   func(Node, int)
	)
	visit = func(n Node, l int) {
		if visited.Contains(n) {
			return
		}
		visited.Add(n)
		g.label.Set(n[0], n[1], n[2], l)
		for _, e := range g.grid.Edges(n) {
			visit(e.To, l)
		}
	}
	for z := 0; z < g.grid.Z; z++ {
		for y := 0; y < g.grid.Y; y++ {
			for x := 0; x < g.grid.X; x++ {
				n := Node{x, y, z}
				if visited.Contains(n) {
					continue
				}
				visit(n, label)
				label++
			}
		}
	}
}

type Grid3[T comparable] struct {
	X    int
	Y    int
	Z    int
	grid []T
}

func NewGrid3[T comparable](x, y, z int) *Grid3[T] {
	return &Grid3[T]{
		X:    x,
		Y:    y,
		Z:    z,
		grid: make([]T, x*y*z),
	}
}

func (g *Grid3[T]) At(x, y, z int) T {
	if 0 <= x && x < g.X && 0 <= y && y < g.Y && 0 <= z && z < g.Z {
		return g.grid[z*g.Y*g.X+y*g.X+x]
	}
	return *new(T)
}

func (g *Grid3[T]) Set(x, y, z int, v T) {
	if 0 <= x && x < g.X && 0 <= y && y < g.Y && 0 <= z && z < g.Z {
		g.grid[z*g.Y*g.X+y*g.X+x] = v
	}
}

func (g *Grid3[T]) Valid(x, y, z int) bool {
	return -1 <= x && x <= g.X && -1 <= y && y <= g.Y && -1 <= z && z <= g.Z
}

func (g *Grid3[T]) SetAll(v T, n ...Node) {
	for _, n := range n {
		g.Set(n[0], n[1], n[2], v)
	}
}

func (g *Grid3[T]) Edges(n Node) []Edge {
	e := g.Neighbors(n)
	v := g.At(n[0], n[1], n[2])
	out := e[:0]
	for _, e := range e {
		if g.At(e.To[0], e.To[1], e.To[2]) == v {
			out = append(out, e)
		}
	}
	return out
}

func (g *Grid3[T]) Neighbors(n Node) []Edge {
	var out []Edge
	do := func(δx, δy, δz int) {
		x, y, z := n[0]+δx, n[1]+δy, n[2]+δz
		if g.Valid(x, y, z) {
			out = append(out, Edge{n, Node{x, y, z}})
		}
	}
	do(-1, 0, 0)
	do(1, 0, 0)
	do(0, -1, 0)
	do(0, 1, 0)
	do(0, 0, -1)
	do(0, 0, 1)
	return out
}

func (g *Graph) From(e Edge) [3]int { return e.From }

func (g *Graph) To(e Edge) [3]int { return e.To }

type Droplet = set.Set[Node]

func (g *Graph) Droplets() map[int]set.Set[Node] {
	m := make(map[int]set.Set[Node])
	for z := 0; z < g.grid.Z; z++ {
		for y := 0; y < g.grid.Y; y++ {
			for x := 0; x < g.grid.X; x++ {
				if g.grid.At(x, y, z) != Water {
					continue
				}
				l := g.label.At(x, y, z)
				d, ok := m[l]
				if !ok {
					d = make(set.Set[Node])
					m[l] = d
				}
				d.Add(Node{x, y, z})
			}
		}
	}
	return m
}

func (g *Graph) SurfaceArea(d Droplet) int {
	var area int
	for n := range d {
		area += 6 - len(g.grid.Edges(n))
	}
	return area
}

func (g *Graph) OutsideSurfaceArea(d Droplet) int {
	var area int
	for n := range d {
		δ := 0
		for _, e := range g.grid.Neighbors(n) {
			if g.label.At(e.To.Coords()) == 0 {
				δ++
			}
		}
		area += δ
	}
	return area
}
