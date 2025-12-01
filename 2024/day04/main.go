package main

import (
	"fmt"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/grid"
)

func main() {
	g, err := grid.Read(os.Stdin, func(r rune) (rune, error) { return r, nil })
	if err != nil {
		log.Fatal(err)
	}
	var total int
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			p := grid.Pos{r, c}
			if g.At(p) != 'X' {
				continue
			}
		δloop:
			for _, δ := range []grid.Pos{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
				q := p
				for _, want := range "MAS" {
					q = q.Add(δ)
					if !g.Valid(q) || g.At(q) != want {
						continue δloop
					}
				}
				total++
			}
		}
	}
	fmt.Println(total)
	total = 0
	for r := 1; r < g.H; r++ {
		for c := 1; c < g.W; c++ {
			p := grid.Pos{r, c}
			if g.At(p) != 'A' {
				continue
			}

			m := 0
		δloop2:
			for _, δ := range []grid.Pos{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}} {
				q := p.Sub(δ)
				for _, want := range "MAS" {
					if !g.Valid(q) || g.At(q) != want {
						continue δloop2
					}
					q = q.Add(δ)
				}
				m++
			}
			if m == 2 {
				total++
			}
		}
	}
	fmt.Println(total)

}
