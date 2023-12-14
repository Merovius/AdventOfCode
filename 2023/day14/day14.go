//go:build goexperiment.rangefunc

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
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

func Parse(s string) (*Grid, error) {
	return grid.Read[Cell](strings.NewReader(s), grid.Enum[Cell]('.', 'O', '#'))
}

func Part1(g *Grid) int {
	SlideNorth(g)
	return Load(g)
}

func Part2(g *Grid) int {
	seen := make(map[string]int)

	const c = 1000000000
	i := 0
	for ; i < c; i++ {
		Cycle(g)
		if j, ok := seen[string(g.G)]; ok {
			δ := i - j
			for ; i+δ < c; i += δ {
			}
			i++
			break
		}
		seen[string(g.G)] = i
	}
	for ; i < c; i++ {
		Cycle(g)
	}
	return Load(g)
}

type Grid = grid.Grid[Cell]

type Cell uint8

const (
	Empty Cell = iota
	Round
	Cube
)

func (c Cell) String() string {
	switch c {
	case Empty:
		return "."
	case Round:
		return "O"
	case Cube:
		return "#"
	default:
		panic(fmt.Errorf("invalid cell %d", c))
	}
}

func Cycle(g *Grid) {
	SlideNorth(g)
	SlideWest(g)
	SlideSouth(g)
	SlideEast(g)
}

func SlideNorth(g *Grid) {
	for p := (grid.Pos{0, 0}); p.Col < g.W; p.Col++ {
		top := grid.Pos{0, p.Col}
		for p.Row = 0; p.Row < g.H; p.Row++ {
			switch g.At(p) {
			case Round:
				g.Set(p, Empty)
				g.Set(top, Round)
				top.Row++
			case Cube:
				top.Row = p.Row + 1
			}
		}
	}
}

func SlideWest(g *Grid) {
	for p := (grid.Pos{0, 0}); p.Row < g.H; p.Row++ {
		top := grid.Pos{p.Row, 0}
		for p.Col = 0; p.Col < g.W; p.Col++ {
			switch g.At(p) {
			case Round:
				g.Set(p, Empty)
				g.Set(top, Round)
				top.Col++
			case Cube:
				top.Col = p.Col + 1
			}
		}
	}
}

func SlideSouth(g *Grid) {
	for p := (grid.Pos{g.H - 1, 0}); p.Col < g.W; p.Col++ {
		top := grid.Pos{g.H - 1, p.Col}
		for p.Row = g.H - 1; p.Row >= 0; p.Row-- {
			switch g.At(p) {
			case Round:
				g.Set(p, Empty)
				g.Set(top, Round)
				top.Row--
			case Cube:
				top.Row = p.Row - 1
			}
		}
	}
}

func SlideEast(g *Grid) {
	for p := (grid.Pos{0, g.W - 1}); p.Row < g.H; p.Row++ {
		top := grid.Pos{p.Row, g.W - 1}
		for p.Col = g.W - 1; p.Col >= 0; p.Col-- {
			switch g.At(p) {
			case Round:
				g.Set(p, Empty)
				g.Set(top, Round)
				top.Col--
			case Cube:
				top.Col = p.Col - 1
			}
		}
	}
}

func Load(g *Grid) int {
	var n int
	for p, c := range g.Cells {
		if c == Round {
			n += g.H - p.Row
		}
	}
	return n
}

func Dump(g *Grid) {
	for p := (grid.Pos{0, 0}); p.Row < g.H; p.Row++ {
		for p.Col = 0; p.Col < g.W; p.Col++ {
			fmt.Print(g.At(p))
		}
		fmt.Println()
	}
}
