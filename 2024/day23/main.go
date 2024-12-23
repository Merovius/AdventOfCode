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
	triangles := make(set.Set[[3]Node])
	// a triangle is a path
	// n1→n2→n3
	// where n3 is a neighbor of n1 (assuming no node is connected to
	// itself).
	for n1, neighbors := range g {
		if n1.String()[0] != 't' {
			continue
		}
		for n2 := range neighbors {
			for n3 := range g[n2] {
				if !neighbors.Contains(n3) {
					continue
				}
				t := [3]Node{n1, n2, n3}
				slices.SortFunc(t[:], Node.Compare)
				triangles.Add(t)
			}
		}
	}
	return len(triangles)
}

func Part2(g Graph) string {
	// This is not technically correct. For example:
	// A─B─C
	// │╱│╲│
	// D─┼─E
	// │╲│╱│
	// F─G─H
	// Every node in this graph is part of a maximal 3-clique (containing
	// A,C,F or H). But there is also a 4-clique {B,D,E,G}, which the
	// algorithm might not find.
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
