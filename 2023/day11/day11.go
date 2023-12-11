package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/math"
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

func Parse(s string) (*Image, error) {
	return grid.Read[Cell](strings.NewReader(s), ParseCell)
}

type Image = grid.Grid[Cell]

type Cell bool

const (
	Empty  Cell = false
	Galaxy Cell = true
)

func ParseCell(r rune) (Cell, error) {
	switch r {
	case '.':
		return Empty, nil
	case '#':
		return Galaxy, nil
	default:
		return false, fmt.Errorf("invalid codepoint %q", r)
	}
}

func Part1(in *Image) int {
	return SumDistances(in, 1)
}

func Part2(in *Image) int {
	return SumDistances(in, 999999)
}

func SumDistances(in *Image, age int) int {
	g := make([]grid.Pos, 0, 1024)
	for r := range in.H {
		for c := range in.W {
			if p := (grid.Pos{r, c}); in.At(p) == Galaxy {
				g = append(g, p)
			}
		}
	}
	emptyRows := make([]int, 0, in.H)
	for r := range in.H {
		empty := true
		for c := range in.W {
			if in.At(grid.Pos{r, c}) == Galaxy {
				empty = false
				break
			}
		}
		if empty {
			emptyRows = append(emptyRows, r)
		}
	}
	emptyCols := make([]int, 0, in.W)
	for c := range in.W {
		empty := true
		for r := range in.H {
			if in.At(grid.Pos{r, c}) == Galaxy {
				empty = false
				break
			}
		}
		if empty {
			emptyCols = append(emptyCols, c)
		}
	}

	var sum int
	for i, g1 := range g {
		for _, g2 := range g[i+1:] {
			// Galaxies are ordered by row first and column second.
			// So g1.Row <= g2.Row
			δr, δc := g2.Row-g1.Row, math.Abs(g2.Col-g1.Col)
			i := search(emptyRows, g1.Row)
			i = search(emptyRows[i:], g2.Row)
			δr += age * i
			i = search(emptyCols, min(g1.Col, g2.Col))
			i = search(emptyCols[i:], max(g1.Col, g2.Col))
			δc += age * i
			sum += δr + δc
		}
	}
	return sum
}

// search returns the smallest index i of s where s[i] > v.
func search[E cmp.Ordered](s []E, v E) int {
	// our slices are short, so linear search is faster than binary search.
	for i, w := range s {
		if w > v {
			return i
		}
	}
	return len(s)
}
