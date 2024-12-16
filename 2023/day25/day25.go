package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/set"
	"github.com/Merovius/AdventOfCode/internal/xiter"
)

// want example:
// hfx/pzl
// bvb/cmg
// nvd/jqt

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part1:", Part1(in))
}

func Parse(s string) (map[string][]string, error) {
	return parse.Map(
		split.Lines,
		split.On(": "),
		parse.String[string],
		parse.Fields(parse.String[string]),
	)(s)
}

func Part1(in map[string][]string) int {
	g := MakeGraph(in)
	for {
		// Choose two random nodes for source/sink. We know that the minimum
		// cut between them should have size 3, so try until that's what we
		// get.
		source, sink := g.ChooseNode(), g.ChooseNode()
		cut, _, reachable := graph.MinimumCut(g, source, sink)
		if cut != 3 {
			continue
		}
		r := xiter.Len(reachable)
		return r * (len(g) - r)
	}
}

type Graph map[string]set.Set[string]

func MakeGraph(in map[string][]string) Graph {
	g := make(Graph)
	for from, tos := range in {
		g[from] = make(set.Set[string])
		for _, to := range tos {
			g[to] = make(set.Set[string])
		}
	}
	for from, tos := range in {
		for _, to := range tos {
			g[from].Add(to)
			g[to].Add(from)
		}
	}
	return g
}

func (g Graph) Edges(n string) iter.Seq[[2]string] {
	return func(yield func([2]string) bool) {
		for to := range g[n] {
			if !yield([2]string{n, to}) {
				return
			}
		}
	}
}

func (g Graph) From(e [2]string) string {
	return e[0]
}

func (g Graph) To(e [2]string) string {
	return e[1]
}

func (g Graph) Weight(_ [2]string) int {
	return 1
}

func (g Graph) Reverse(e [2]string) [2]string {
	e[0], e[1] = e[1], e[0]
	return e
}

// Chooses a random node
func (g Graph) ChooseNode() string {
	for n := range g {
		return n
	}
	panic("no nodes in graph")
}
