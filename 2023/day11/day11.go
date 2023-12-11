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

func Parse(in string) (Input, error) {
	g := make([]grid.Pos, 0, 1024)
	var nrows, ncols int
	for r := 0; len(in) > 0; r++ {
		i := strings.IndexByte(in, '\n')
		if i < 0 {
			i = len(in)
		}
		l := in[:i]
		ncols = len(l)
		nrows++
		var lastC int
		for len(l) > 0 {
			c := strings.IndexByte(l, '#')
			if c < 0 {
				break
			}
			g = append(g, grid.Pos{r, c + lastC})
			lastC += c + 1
			l = l[c+1:]
		}
		in = in[i+1:]
	}
	return Input{g, nrows, ncols}, nil
}

type Input struct {
	Galaxies []grid.Pos
	NRows    int
	NCols    int
}

func Part1(in Input) int {
	return SumDistances(in, 1)
}

func Part2(in Input) int {
	return SumDistances(in, 999999)
}

func SumDistances(in Input, age int) int {
	emptyRows := make([]int, 0, in.NRows)
rows:
	for r := range in.NRows {
		for _, g := range in.Galaxies {
			if g.Row == r {
				continue rows
			}
		}
		emptyRows = append(emptyRows, r)
	}
	emptyCols := make([]int, 0, in.NCols)
cols:
	for c := range in.NCols {
		for _, g := range in.Galaxies {
			if g.Col == c {
				continue cols
			}
		}
		emptyCols = append(emptyCols, c)
	}

	var sum int
	for i, g1 := range in.Galaxies {
		for _, g2 := range in.Galaxies[i+1:] {
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
