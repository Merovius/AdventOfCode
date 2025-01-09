package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/math"
	"github.com/Merovius/AdventOfCode/internal/xiter"
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
	fmt.Println(Part2(in))
}

func Parse(in []byte) ([]Edge, error) {
	return parse.Slice(
		split.Lines,
		parse.Struct[Edge](
			split.Regexp(`(\w+) to (\w+) = (\d+)`),
			parse.String[string],
			parse.String[string],
			parse.Signed[int],
		),
	)(string(in))
}

type Edge struct {
	From string
	To   string
	Dist int
}

func Part1(in []Edge) int {
	return xiter.FoldL(math.Min, PathLengths(in), math.MaxInt)
}

func Part2(in []Edge) int {
	return xiter.FoldL(math.Max, PathLengths(in), math.MinInt)
}

func PathLengths(in []Edge) iter.Seq[int] {
	var (
		nodes []string
		idx   = make(map[string]int)
		dist  = make(map[[2]int]int)
	)
	lookup := func(s string) int {
		i, ok := idx[s]
		if !ok {
			i = len(nodes)
			nodes = append(nodes, s)
			idx[s] = i
		}
		return i
	}

	for _, e := range in {
		i := lookup(e.From)
		j := lookup(e.To)
		dist[[2]int{i, j}] = e.Dist
		dist[[2]int{j, i}] = e.Dist
	}
	var rec func(p, q []int, yield func(int) bool) bool
	rec = func(p, q []int, yield func(int) bool) bool {
		if len(q) == 0 {
			var d int
			for i := 1; i < len(p); i++ {
				d += dist[[2]int{p[i-1], p[i]}]
			}
			return yield(d)
		}
		n, q := q[len(q)-1], q[:len(q)-1]
		if !rec(append(p, n), q, yield) {
			return false
		}
		for i := range q {
			q[i], n = n, q[i]
			if !rec(append(p, n), q, yield) {
				return false
			}
			q[i], n = n, q[i]
		}
		return true
	}
	return func(yield func(int) bool) {
		p := make([]int, 0, len(nodes))
		q := make([]int, len(nodes))
		for i := range q {
			q[i] = i
		}
		rec(p, q, yield)
	}
}
