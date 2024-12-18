package main

import (
	"flag"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"sort"

	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	log.SetFlags(0)
	w := flag.Int("w", 0, "width")
	h := flag.Int("h", 0, "height")
	n := flag.Int("h", 0, "number of bytes")
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
	g, G := MakeGraph(w, h)
	for _, xy := range in[:n] {
		g.Set(grid.Pos{xy[1], xy[0]}, Corrupted)
	}
	return ShortestPath(G, w, h)
}

func Part2(in [][2]int, w, h int) [2]int {
	g, G := MakeGraph(w, h)
	i := sort.Search(len(in), func(i int) bool {
		clear(g.G)
		for _, xy := range in[:i] {
			g.Set(grid.Pos{xy[1], xy[0]}, Corrupted)
		}
		return ShortestPath(G, w, h) == 0
	})
	return in[i-1]
}

type Cell byte

const (
	Empty Cell = iota
	Corrupted
)

type Grid = grid.Grid[Cell]
type Graph = graph.Graph[grid.Pos, [2]grid.Pos]

func MakeGraph(w, h int) (*Grid, Graph) {
	g := grid.New[Cell](w, h)
	return g, graph.NeighborFunc(func(p grid.Pos) iter.Seq[grid.Pos] {
		return func(yield func(grid.Pos) bool) {
			for _, q := range g.Neigh4(p) {
				if g.At(q) != Empty {
					continue
				}
				if !yield(q) {
					return
				}
			}
		}
	})
}

func ShortestPath(G Graph, w, h int) int {
	return len(graph.BreadthFirstSearch(G, grid.Pos{0, 0}, func(p grid.Pos) bool {
		return p.Row == h-1 && p.Col == w-1
	}))
}
