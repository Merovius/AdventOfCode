package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	log.SetFlags(0)
	save := flag.Int("save", 100, "Minimum time to save")
	flag.Parse()

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in, *save))
	fmt.Println(Part2(in, *save))
}

func Parse(buf []byte) (*grid.Grid[Cell], error) {
	return grid.Read(bytes.NewReader(buf), func(r rune) (Cell, error) {
		if i := strings.IndexRune(".#SE", r); i >= 0 {
			return Cell(i), nil
		}
		return 0, fmt.Errorf("invalid character %q", r)
	})
}

type Cell byte

const (
	Empty Cell = iota
	Wall
	Start
	End
)

func Part1(g *grid.Grid[Cell], save int) int {
	return findCheats(g, save, 2)
}

func Part2(g *grid.Grid[Cell], save int) int {
	return findCheats(g, save, 20)
}

func findCheats(g *grid.Grid[Cell], save, cheat int) int {
	var (
		start = g.Pos(slices.Index(g.G, Start))
		p     = start
		d     = grid.Up
		path  = []grid.Pos{p}
		dist  = grid.New[int](g.W, g.H)
	)
	for g.At(p) != End {
		for _, dd := range []grid.Direction{d, d.RotateRight(), d.RotateLeft()} {
			if q := dd.Move(p); g.At(q) != Wall {
				path, p, d = append(path, q), q, dd
				break
			}
		}
		dist.Set(p, len(path)-1)
	}
	var N int
	for i, p := range path[:len(path)-save] {
		// Find all valid cheat targets:
		// - δr varies from [-cheat,cheat]
		// - δc varies from [-cheat+|δr|,cheat-|δr|]
		// - Thus |δr|+|δc| ≤ cheat
		// - Also restrict that p.Row+δr ∈ [0,g.H) and p.Col+δc ∈ [0,g.W)
		lo := max(-cheat, -p.Row)
		hi := min(cheat, g.H-1-p.Row)
		for δr := lo; δr <= hi; δr++ {
			aδr := math.Abs(δr)
			lo := max(-cheat+aδr, -p.Col)
			hi := min(cheat-aδr, g.W-1-p.Col)
			for δc := lo; δc <= hi; δc++ {
				aδc := math.Abs(δc)

				// i+|δr|+|δc| is the distance from start to q,
				// if cheating.
				// dist.At(q) is the normal distance from start
				// to q.
				q := p.Add(grid.Pos{δr, δc})
				if dist.At(q)-(i+aδr+aδc) >= save {
					N++
				}
			}
		}
	}
	return N
}
