package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/graph"
)

func main() {
	g, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	part1 := &Graph{
		Grid:      g,
		ValidEdge: func(a, b int) bool { return b-a <= 1 },
	}
	path := graph.ShortestPath[Point, Edge](part1, g.Start, func(p Point) bool { return p == g.End })
	if path == nil {
		log.Fatal("No path found")
	}
	fmt.Printf("Shortest path from %v to %v takes %d steps\n", g.Start, g.End, len(path))

	part2 := &Graph{
		Grid:      g,
		ValidEdge: func(a, b int) bool { return a-b <= 1 },
	}
	path = graph.ShortestPath[Point, Edge](part2, g.End, func(p Point) bool { return g.At(p) == 0 })
	if path == nil {
		log.Fatal("No path found")
	}
	fmt.Printf("Shortest path from %v to 'a' takes %d steps\n", g.End, len(path))
}

func ReadInput(r io.Reader) (*Grid, error) {
	var cells []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		cells = append(cells, strings.TrimSpace(s.Text()))
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(cells) == 0 {
		return nil, errors.New("empty input")
	}
	sr, sc, er, ec, w := -1, -1, -1, -1, len(cells[0])

	for r := range cells {
		if len(cells[r]) != w {
			return nil, errors.New("inconsistent line length")
		}
		for c := range cells[r] {
			if cells[r][c] == 'S' {
				sr, sc = r, c
			}
			if cells[r][c] == 'E' {
				er, ec = r, c
			}
		}
	}
	return &Grid{
		Start:  Point{sr, sc},
		End:    Point{er, ec},
		Width:  w,
		Height: len(cells),
		Cells:  cells,
	}, nil
}

type Grid struct {
	Start  Point
	End    Point
	Width  int
	Height int
	Cells  []string
}

type Point struct {
	Row int
	Col int
}

func (p Point) Add(r, c int) Point {
	return Point{p.Row + r, p.Col + c}
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row+1, p.Col+1)
}

type Edge struct {
	From Point
	To   Point
}

type Graph struct {
	*Grid
	ValidEdge func(int, int) bool
}

func (g *Graph) Edges(p Point) []Edge {
	var out []Edge
	for _, δ := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		n := p.Add(δ[0], δ[1])
		if !g.Valid(n) {
			continue
		}
		if !g.ValidEdge(g.At(p), g.At(n)) {
			continue
		}
		out = append(out, Edge{p, n})
	}
	return out
}

func (g *Graph) From(e Edge) Point { return e.From }

func (g *Graph) To(e Edge) Point { return e.To }

func (g *Grid) At(p Point) int {
	v := g.Cells[p.Row][p.Col]
	if v == 'S' {
		v = 'a'
	}
	if v == 'E' {
		v = 'z'
	}
	return int(v - 'a')
}

func (g *Grid) Valid(p Point) bool {
	return p.Row >= 0 && p.Row < g.Height && p.Col >= 0 && p.Col < g.Width
}
