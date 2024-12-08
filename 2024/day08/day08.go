package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
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
				for g.Valid(q) {
					nodes.Add(q)
					q = q.Add(δ)
				}
			}
		}
	}
	fmt.Println(len(nodes))
}

type Cell byte

const Empty Cell = '.'
