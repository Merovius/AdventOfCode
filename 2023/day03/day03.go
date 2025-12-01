package main

import (
	"fmt"
	"iter"
	"log"
	"os"
	"strconv"
	"unicode"

	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/set"
	"gonih.org/AdventOfCode/internal/xiter"
)

func main() {
	g, err := grid.Read[rune](os.Stdin, func(r rune) (rune, error) { return r, nil })
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sum of part numbers:", Part1(g))
	fmt.Println("Sum of gear rations:", Part2(g))
}

func Part1(g *grid.Grid[rune]) int {
	var sum int
	for n := range FindPartNumbers(g) {
		sum += n.N
	}
	return sum
}

func Part2(g *grid.Grid[rune]) int {
	partNumbers := set.Collect(FindPartNumbers(g))

	var sum int
	for p := range grid.Find(g, '*') {
		var adj []Number
		for n := range partNumbers {
			if n.Neighborhood().Contains(p) {
				adj = append(adj, n)
			}
		}
		if len(adj) != 2 {
			continue
		}
		sum += adj[0].N * adj[1].N
	}
	return sum
}

type Number struct {
	N int
	R grid.Rectangle
}

func (n Number) Neighborhood() grid.Rectangle {
	return n.R.Inset(-1)
}

func FindNumbers(g *grid.Grid[rune]) iter.Seq[Number] {
	return func(yield func(Number) bool) {
		p := grid.Pos{}
		var b []rune
		for p.Row = range g.H {
			var n Number
			for p.Col = 0; p.Col < g.W; p.Col++ {
				n.R.Min, n.R.Max = p, p
				for g.Valid(n.R.Max) && unicode.IsDigit(g.At(n.R.Max)) {
					b = append(b, g.At(n.R.Max))
					n.R.Max.Col++
				}
				if len(b) > 0 {
					v, err := strconv.Atoi(string(b))
					if err != nil {
						panic(err)
					}
					p = n.R.Max
					n.N = v
					n.R.Max.Row++
					if !yield(n) {
						return
					}
					b = b[:0]
				}
			}
		}
	}
}

func FindPartNumbers(g *grid.Grid[rune]) iter.Seq[Number] {
	return xiter.Filter(FindNumbers(g), func(n Number) bool {
		for _, r := range g.Rect(n.Neighborhood()) {
			if r != '.' && !unicode.IsDigit(r) {
				return true
			}
		}
		return false
	})
}
