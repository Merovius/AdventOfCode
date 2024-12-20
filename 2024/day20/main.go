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
		p    = g.Pos(slices.Index(g.G, Start))
		d    = grid.Up
		path = []grid.Pos{p}
	)
	for g.At(p) != End {
		for _, dd := range []grid.Direction{d, d.RotateRight(), d.RotateLeft()} {
			if q := dd.Move(p); g.At(q) != Wall {
				path, p, d = append(path, q), q, dd
				break
			}
		}
	}
	var N int
	for i, p := range path[:len(path)-save] {
		// jump ahead save steps and see if we can reach it by
		// cheating.
		for j, q := range path[i+save:] {
			if d := q.Sub(p).Length(); d <= cheat && d <= j {
				N++
			}
		}
	}
	return N
}
