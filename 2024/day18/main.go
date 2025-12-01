package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"gonih.org/AdventOfCode/internal/container"
	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

func main() {
	log.SetFlags(0)
	w := flag.Int("w", 0, "width")
	h := flag.Int("h", 0, "height")
	n := flag.Int("n", 0, "number of bytes")
	flag.Parse()
	if *w == 0 || *h == 0 || *n == 0 {
		log.Fatal("Both -w, -h and -n are required")
	}

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in, *w, *h, *n))
	fmt.Println(Part2(in, *w, *h))
}

func Parse(in []byte) ([][2]int, error) {
	return parse.Slice(split.Lines, parse.Array[[2]int](split.On(","), parse.Signed[int]))(string(in))
}

func Part1(in [][2]int, w, h, n int) int {
	g := grid.New[Cell](w, h)
	for _, xy := range in[:n] {
		g.Set(grid.Pos{xy[1], xy[0]}, Corrupted)
	}
	return ShortestPath(g)
}

func Part2(in [][2]int, w, h int) [2]int {
	g := grid.New[Cell](w, h)
	i := sort.Search(len(in), func(i int) bool {
		clear(g.G)
		for _, xy := range in[:i] {
			g.Set(grid.Pos{xy[1], xy[0]}, Corrupted)
		}
		return ShortestPath(g) == 0
	})
	return in[i-1]
}

type Cell uint16

const Corrupted Cell = math.MaxUint16

func ShortestPath(g *grid.Grid[Cell]) int {
	// Specialized BFS, as graph.BreathFirstSearch allocates too much.
	type el struct {
		p grid.Pos
		d Cell
	}
	start := grid.Pos{}
	end := grid.Pos{g.W - 1, g.H - 1}
	q := container.MakeFIFO[el](g.W * g.H)
	q.Push(el{start, 0})
	for q.Len() > 0 {
		e := q.Pop()
		if g.At(e.p) > 0 {
			continue
		}
		g.Set(e.p, e.d)
		for p, c := range g.Neigh4(e.p) {
			if c > 0 {
				continue
			}
			if p == end {
				return int(e.d + 1)
			}
			q.Push(el{p, e.d + 1})
		}

	}
	return 0
}
