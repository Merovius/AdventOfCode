package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"slices"

	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
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
	// Note: Technically, this is incorrect. The input could only specify a
	// partial order and it would have to be made into a transitive closure.
	// That is, the input could be
	//     A|B
	//     B|C
	//
	//     B,X,A
	//     C,A,B
	// Which would be two unsorted series, but the relation we are building
	// would not "see" either of them.
	// However, it turns out our input specifies a total ordering for all
	// the actually appearing sequences, so…
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
			rand.Shuffle(len(p), func(i, j int) {
				p[i], p[j] = p[j], p[i]
			})

			slices.SortFunc(p, cmp)
			m += p[len(p)/2]
		}
	}
	fmt.Println(n)
	fmt.Println(m)
}
