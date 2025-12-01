package main

import (
	"fmt"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/math"
	"gonih.org/AdventOfCode/internal/set"
)

func main() {
	g, err := grid.Read(os.Stdin, func(r rune) (Cell, error) {
		return Cell(r), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	freqs := make(map[Cell][]grid.Pos)
	for p, c := range g.All() {
		if c != Empty {
			freqs[c] = append(freqs[c], p)
		}
	}
	nodes := make(set.Set[grid.Pos])
	for _, ps := range freqs {
		for _, p := range ps {
			for _, q := range ps {
				if p == q {
					continue
				}
				if n := q.Add(q.Sub(p)); g.Valid(n) {
					nodes.Add(n)
				}
			}
		}
	}
	fmt.Println(len(nodes))

	for _, ps := range freqs {
		for _, p := range ps {
			for _, q := range ps {
				if p == q {
					continue
				}
				δ := q.Sub(p)
				z, _, _ := math.GCD(δ.Row, δ.Col)
				δ.Row, δ.Col = δ.Row/z, δ.Col/z
				for x := p; g.Valid(x); x = x.Add(δ) {
					nodes.Add(x)
				}
			}
		}
	}
	fmt.Println(len(nodes))
}

type Cell byte

const Empty Cell = '.'
