package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	type Manual struct {
		Rules [][2]int
		Pages [][]int
	}
	manual, err := parse.Struct[Manual](
		split.Blocks,
		parse.Slice(
			split.Lines,
			parse.Array[[2]int](
				split.On("|"),
				parse.Signed[int],
			),
		),
		parse.Slice(
			split.Lines,
			parse.Slice(
				split.On(","),
				parse.Signed[int],
			),
		),
	)(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	g := grid.New[int](100, 100)
	for _, r := range manual.Rules {
		g.Set(grid.Pos{r[0], r[1]}, -1)
		g.Set(grid.Pos{r[1], r[0]}, 1)
	}
	cmp := func(a, b int) int {
		return g.At(grid.Pos{a, b})
	}
	var (
		n int
		m int
	)
	for _, p := range manual.Pages {
		if slices.IsSortedFunc(p, cmp) {
			n += p[len(p)/2]
		} else {
			slices.SortFunc(p, cmp)
			m += p[len(p)/2]
		}
	}
	fmt.Println(n)
	fmt.Println(m)
}
