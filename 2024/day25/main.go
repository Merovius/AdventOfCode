package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in))
}

func Parse(in []byte) ([]*grid.Grid[bool], error) {
	grids, err := parse.Slice(
		split.Blocks,
		func(s string) (*grid.Grid[bool], error) {
			return grid.Read(strings.NewReader(s), func(r rune) (bool, error) {
				switch r {
				case '.':
					return false, nil
				case '#':
					return true, nil
				default:
					return false, fmt.Errorf("invalid character %q in input", r)
				}
			})
		},
	)(string(in))
	if err != nil {
		return nil, err
	}
	if len(grids) == 0 {
		return nil, errors.New("no grids in input")
	}
	w, h := grids[0].W, grids[0].H
	for _, g := range grids[1:] {
		if g.W != w || g.H != h {
			return nil, errors.New("grids do not all have the same size")
		}
	}
	return grids, nil
}

func Part1(in []*grid.Grid[bool]) int {
	var (
		locks [][]int
		keys  [][]int
	)
	for _, g := range in {
		var heights []int
		for p := (grid.Pos{}); p.Col < g.W; p.Col++ {
			var h int
			for p.Row = 0; p.Row < g.H; p.Row++ {
				if g.At(p) {
					h++
				}
			}
			heights = append(heights, h)
		}
		if g.At(grid.Pos{0, 0}) {
			locks = append(locks, heights)
		} else {
			keys = append(keys, heights)
		}
	}

	H := in[0].H
	var n int
	for _, k := range keys {
		for _, l := range locks {
			for i := range l {
				if l[i]+k[i] > H {
					n++
					break
				}
			}
		}
	}
	return len(keys)*len(locks) - n
}
