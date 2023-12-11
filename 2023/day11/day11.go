package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/interval"
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
	var g []grid.Pos
	for p := range grid.Find(in, Galaxy) {
		g = append(g, p)
	}
	var emptyRows []int
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
	var emptyCols []int
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
			if g1 == g2 {
				continue
			}
			ir := interval.MakeOC(g1.Row, g2.Row)
			ic := interval.MakeOC(g1.Col, g2.Col)
			δr, δc := ir.Len(), ic.Len()
			for _, r := range emptyRows {
				if ir.Contains(r) {
					δr += age
				}
			}
			for _, c := range emptyCols {
				if ic.Contains(c) {
					δc += age
				}
			}
			sum += δr + δc
		}
	}
	return sum
}
