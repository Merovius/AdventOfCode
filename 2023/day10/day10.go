package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(g))
	fmt.Println(Part2(g))
}

func Parse(s string) (*Grid, error) {
	return grid.Read(strings.NewReader(s), ParseCell)
}

const (
	North = (1 << (iota + 1))
	East
	South
	West
)

type Direction int

func Directions() []Direction {
	return []Direction{North, East, South, West}
}

func (d Direction) invert() Direction {
	switch d {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	}
	panic(fmt.Errorf("invalid direction %d", d))
}

func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return fmt.Sprintf("Direction(%#x)", d)
	}
}

type Cell uint8

const (
	Ground Cell = 0
	Start  Cell = 1
	Vert   Cell = North | South
	Horiz  Cell = East | West
	BendNE Cell = North | East
	BendSE Cell = South | East
	BendSW Cell = South | West
	BendNW Cell = North | West
)

func ParseCell(r rune) (Cell, error) {
	switch r {
	case 'S':
		return Start, nil
	case '|':
		return North | South, nil
	case '-':
		return East | West, nil
	case 'L':
		return North | East, nil
	case 'J':
		return North | West, nil
	case 'F':
		return South | East, nil
	case '7':
		return West | South, nil
	case '.':
		return Ground, nil
	}
	return 0, fmt.Errorf("invalid cell %q", r)
}

var cellStr = map[Cell]string{
	Ground:         "░",
	Start:          "•",
	Vert:           "│",
	Vert | Start:   "║",
	Horiz:          "─",
	Horiz | Start:  "═",
	BendNE:         "└",
	BendNE | Start: "╚",
	BendSE:         "┌",
	BendSE | Start: "╔",
	BendSW:         "┐",
	BendSW | Start: "╗",
	BendNW:         "┘",
	BendNW | Start: "╝",
}

func (c Cell) String() string {
	s, ok := cellStr[c]
	if ok {
		return s
	}
	return fmt.Sprintf("Cell(%#x)", c)
}

type Grid = grid.Grid[Cell]

func Part1(g *Grid) int {
	return len(findLoop(g)) / 2
}

func findLoop(g *Grid) []grid.Pos {
	var start grid.Pos
	for g := range grid.Find(g, Start) {
		start = g
		break
	}
	var dir Direction
	for _, d := range []Direction{North, East, South, West} {
		q := d.Of(start)
		if !g.Valid(q) {
			continue
		}
		if g.At(q).Connects(d) {
			dir = d
			break
		}
	}
	if dir == 0 {
		panic("no direction has a connecting pipe")
	}
	loop := []grid.Pos{start}
	for p := dir.Of(start); p != start; p = dir.Of(p) {
		loop = append(loop, p)
		dir = g.At(p).Pass(dir)
	}
	return loop
}

// Connect returns whether p can be moved into when heading into d.
func (p Cell) Connects(d Direction) bool {
	switch d {
	case North:
		return (p & South) != 0
	case South:
		return (p & North) != 0
	case East:
		return (p & West) != 0
	case West:
		return (p & East) != 0
	default:
		panic(fmt.Errorf("invalid direction %d", d))
	}
}

// Pass returns the direction p is left into when entering heading d.
func (p Cell) Pass(d Direction) Direction {
	if Direction(p)&d.invert() == 0 {
		panic(fmt.Errorf("heading %v into %v is illegal", d, p))
	}
	return Direction(p&^Start) ^ d.invert()
}

func (d Direction) Of(p grid.Pos) grid.Pos {
	switch d {
	case North:
		p.Row -= 1
	case East:
		p.Col += 1
	case South:
		p.Row += 1
	case West:
		p.Col -= 1
	default:
		panic(fmt.Errorf("invalid direction %d", d))
	}
	return p
}

func Part2(g *Grid) int {
	loop := findLoop(g)
	// Determine the pipe shape of the start position
	v := direction[loop[1].Sub(loop[0])]
	v |= direction[loop[len(loop)-1].Sub(loop[0])]

	gr := NewGraph(g, loop[0], v)

	// Do a depth-first search to determine the connected component of (0,0)
	// (outside)
	outside := make(set.Set[grid.Pos])
	for p := range graph.WalkDepthFirst(gr, grid.Pos{0, 0}) {
		outside.Add(grid.Pos{p.Row / 3, p.Col / 3})
	}
	return g.W*g.H - len(outside)
}

// maps a direction vector unto the corresponding pipe-connection
var direction = map[grid.Pos]Cell{
	grid.Pos{0, 1}:  East,
	grid.Pos{1, 0}:  South,
	grid.Pos{0, -1}: West,
	grid.Pos{-1, 0}: North,
}

type Wall bool

type Graph struct {
	g *grid.Grid[Wall]
}

// maps a pipe-value to a 3x3 square describing it as a wall
var walls = map[Cell][3][3]Wall{
	Vert:   {{false, true, false}, {false, true, false}, {false, true, false}},
	Horiz:  {{false, false, false}, {true, true, true}, {false, false, false}},
	BendNE: {{false, true, false}, {false, true, true}, {false, false, false}},
	BendNW: {{false, true, false}, {true, true, false}, {false, false, false}},
	BendSW: {{false, false, false}, {true, true, false}, {false, true, false}},
	BendSE: {{false, false, false}, {false, true, true}, {false, true, false}},
}

// NewGraph creates a 3-times scaled up version of g, where every cell is a 3x3
// square with loop segments filled in as wall.
func NewGraph(g *Grid, start grid.Pos, startVal Cell) *Graph {
	big := grid.New[Wall](3*g.W, 3*g.H)
	for p, v := range g.Cells {
		if p == start {
			v = startVal
		}
		p.Row *= 3
		p.Col *= 3
		for r, row := range walls[v] {
			for c, val := range row {
				big.Set(grid.Pos{p.Row + r, p.Col + c}, val)
			}
		}
	}
	return &Graph{big}
}

type Edge struct {
	src grid.Pos
	dst grid.Pos
}

func (g *Graph) Edges(p grid.Pos) []Edge {
	edges := make([]Edge, 0, 4)
	for _, q := range g.g.Neigh4(p) {
		if !g.g.At(q) {
			edges = append(edges, Edge{p, q})
		}
	}
	return edges
}

func (g *Graph) From(e Edge) grid.Pos {
	return e.src
}

func (g *Graph) To(e Edge) grid.Pos {
	return e.dst
}
