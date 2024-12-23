package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/set"
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
	edges, err := parse.Slice(split.Lines, parse.Array[[2]string](split.On("-"), parse.String[string]))(string(in))
	if err != nil {
		return Graph{}, err
	}
	g := Graph{edges: make(map[string]set.Set[string])}
	for _, e := range edges {
		if g.edges[e[0]] == nil {
			g.edges[e[0]] = make(set.Set[string])
		}
		if g.edges[e[1]] == nil {
			g.edges[e[1]] = make(set.Set[string])
		}
		g.edges[e[0]].Add(e[1])
		g.edges[e[1]].Add(e[0])
	}
	return g, nil
}

func Part1(in Graph) int {

	sets := make(set.Set[[3]string])
	for f1, t1 := range in.edges {
		if f1[0] != 't' {
			continue
		}
		for f2, t2 := range in.edges {
			for f3, t3 := range in.edges {
				if !t2.Contains(f1) || !t3.Contains(f1) {
					continue
				}
				if !t1.Contains(f2) || !t3.Contains(f2) {
					continue
				}
				if !t1.Contains(f3) || !t2.Contains(f3) {
					continue
				}
				s := [3]string{f1, f2, f3}
				slices.Sort(s[:])
				sets.Add(s)
			}
		}
	}
	return len(sets)
}

func Part2(in Graph) string {
	var best set.Set[string]
	for n := range in.edges {
		c := maximalClique(in, n)
		if len(c) > len(best) {
			best = c
		}
	}
	s := slices.Collect(maps.Keys(best))
	slices.Sort(s)
	return strings.Join(s, ",")
}

func maximalClique(g Graph, n string) set.Set[string] {
	c := make(set.Set[string])
	c.Add(n)
	for {
		before := len(c)
	loop:
		for n := range g.edges {
			if c.Contains(n) {
				continue
			}
			for m := range c {
				if !g.edges[n].Contains(m) {
					continue loop
				}
			}
			c.Add(n)
		}
		if len(c) == before {
			return c
		}
	}
}

type Graph struct {
	edges map[string]set.Set[string]
}
