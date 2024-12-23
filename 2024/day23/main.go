package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"unique"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/set"
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

func Parse(in []byte) (Graph, error) {
	edges, err := parse.Slice(
		split.Lines,
		parse.Array[[2]Node](
			split.On("-"),
			ParseNode,
		),
	)(string(in))
	if err != nil {
		return nil, err
	}
	g := make(Graph)
	for _, e := range edges {
		g.Add(e[0], e[1])
	}
	return g, nil
}

func Part1(g Graph) int {
	cliques := make(set.Set[[3]Node])
	for f1, t1 := range g {
		if f1.String()[0] != 't' {
			continue
		}
		for f2, t2 := range g {
			for f3, t3 := range g {
				if !t2.Contains(f1) || !t3.Contains(f1) {
					continue
				}
				if !t1.Contains(f2) || !t3.Contains(f2) {
					continue
				}
				if !t1.Contains(f3) || !t2.Contains(f3) {
					continue
				}
				s := [3]Node{f1, f2, f3}
				slices.SortFunc(s[:], Node.Compare)
				cliques.Add(s)
			}
		}
	}
	return len(cliques)
}

func Part2(g Graph) string {
	var best set.Set[Node]
	for n := range g {
		if c := MaximalClique(g, n); len(c) > len(best) {
			best = c
		}
	}
	s := slices.Collect(xiter.Map(best.All(), Node.String))
	slices.Sort(s)
	return strings.Join(s, ",")
}

func MaximalClique(g Graph, n Node) set.Set[Node] {
	c := make(set.Set[Node])
	c.Add(n)
	for {
		before := len(c)
		for n, edges := range g {
			if c.Contains(n) {
				continue
			}
			if c.SubsetOf(edges) {
				c.Add(n)
			}
		}
		if len(c) == before {
			return c
		}
	}
}

type Node unique.Handle[string]

func ParseNode(s string) (Node, error) {
	return Node(unique.Make(s)), nil
}

func (n Node) Compare(m Node) int {
	return cmp.Compare(n.String(), m.String())
}

func (n Node) String() string {
	return (unique.Handle[string])(n).Value()
}

type Graph map[Node]set.Set[Node]

func (g Graph) Add(from, to Node) {
	if g[from] == nil {
		g[from] = make(set.Set[Node])
	}
	g[from].Add(to)
	if g[to] == nil {
		g[to] = make(set.Set[Node])
	}
	g[to].Add(from)
}
