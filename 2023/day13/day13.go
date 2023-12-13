package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
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

func Parse(s string) ([]*Grid, error) {
	return parse.Blocks(func(block string) (*Grid, error) {
		return grid.Read[Cell](strings.NewReader(block), grid.Enum[Cell]('.', '#'))
	})(s)
}

func Part1(in []*Grid) int {
	var rr, cc int
	for _, g := range in {
		if r, ok := FindRowReflection(g); ok {
			rr += r
		}
		if c, ok := FindColReflection(g); ok {
			cc += c
		}
	}
	return 100*rr + cc
}

func Part2(in []*Grid) int {
	var rr, cc int
	for _, g := range in {
		if r, ok := FindRowSmudge(g); ok {
			rr += r
		}
		if c, ok := FindColSmudge(g); ok {
			cc += c
		}
	}
	return 100*rr + cc
}

type Grid = grid.Grid[Cell]

type Cell uint8

const (
	Ash Cell = iota
	Rock
)

func FindRowReflection(g *Grid) (int, bool) {
nextRow:
	for r := 1; r < g.H; r++ {
		for i := 0; r-1-i >= 0 && r+i < g.H; i++ {
			for c := 0; c < g.W; c++ {
				if g.At(grid.Pos{r - 1 - i, c}) != g.At(grid.Pos{r + i, c}) {
					continue nextRow
				}
			}
		}
		return r, true
	}
	return 0, false
}

func FindColReflection(g *Grid) (int, bool) {
nextCol:
	for c := 1; c < g.W; c++ {
		for i := 0; c-1-i >= 0 && c+i < g.W; i++ {
			for r := 0; r < g.H; r++ {
				if g.At(grid.Pos{r, c - 1 - i}) != g.At(grid.Pos{r, c + i}) {
					continue nextCol
				}
			}
		}
		return c, true
	}
	return 0, false
}

func FindRowSmudge(g *Grid) (int, bool) {
nextRow:
	for r := 1; r < g.H; r++ {
		var n int
		for i := 0; r-1-i >= 0 && r+i < g.H; i++ {
			for c := 0; c < g.W; c++ {
				if g.At(grid.Pos{r - 1 - i, c}) != g.At(grid.Pos{r + i, c}) {
					n++
					if n > 1 {
						continue nextRow
					}
				}
			}
		}
		if n == 1 {
			return r, true
		}
	}
	return 0, false
}

func FindColSmudge(g *Grid) (int, bool) {
nextCol:
	for c := 1; c < g.W; c++ {
		var n int
		for i := 0; c-1-i >= 0 && c+i < g.W; i++ {
			for r := 0; r < g.H; r++ {
				if g.At(grid.Pos{r, c - 1 - i}) != g.At(grid.Pos{r, c + i}) {
					n++
					if n > 1 {
						continue nextCol
					}
				}
			}
		}
		if n == 1 {
			return c, true
		}
	}
	return 0, false
}
