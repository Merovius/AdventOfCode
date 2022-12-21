package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input"
	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	data, err := input.Slice(input.Lines(), func(s string) (Valve, error) {
		var v Valve
		before, after, ok := strings.Cut(s, "; ")
		if !ok {
			return Valve{}, fmt.Errorf("invalid line %q: no ;", s)
		}
		n, err := fmt.Sscanf(before, "Valve %s has flow rate=%d", &v.Name, &v.FlowRate)
		if err != nil || n != 2 {
			return Valve{}, fmt.Errorf("invalid line %q: can'd read name/flow rate", s)
		}
		after, ok = strings.CutPrefix(after, "tunnels lead to valves ")
		if !ok {
			after, ok = strings.CutPrefix(after, "tunnel leads to valve ")
			if !ok {
				return Valve{}, fmt.Errorf("invalid line %q: no tunnels", s)
			}
		}
		v.Tunnels = strings.Split(after, ", ")
		return v, nil
	}).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	idx := make(map[string]int)
	for i, v := range data {
		idx[v.Name] = i
	}

	flow := make(map[int]int)
	l := make([][]int, len(data))
	for i, v := range data {
		if v.FlowRate > 0 {
			flow[i] = v.FlowRate
		}
		l[i] = make([]int, len(data))
		for j := range l {
			l[i][j] = math.MaxInt / 2
		}
		l[i][i] = 0
		for _, t := range v.Tunnels {
			l[i][idx[t]] = 1
		}
	}
	for i := range l {
		for j := range l {
			for k := range l {
				l[j][k] = math.Min(l[j][k], l[j][i]+l[i][k])
			}
		}
	}

	var visit func(n, budget, released int, open Bitset, answer map[Bitset]int) map[Bitset]int
	visit = func(n, budget, released int, open Bitset, answer map[Bitset]int) map[Bitset]int {
		answer[open] = math.Max(answer[open], released)
		for k := range flow {
			nb := budget - l[n][k] - 1
			if open.Contains(k) || nb < 0 {
				continue
			}
			visit(k, nb, released+nb*flow[k], open.Add(k), answer)
		}
		return answer
	}

	best := math.MinInt
	for _, v := range visit(idx["AA"], 30, 0, 0, make(map[Bitset]int)) {
		best = math.Max(best, v)
	}
	fmt.Printf("You alone can release %d pressure in 30m\n", best)

	answer := visit(idx["AA"], 26, 0, 0, make(map[Bitset]int))
	best = math.MinInt
	for s1, v1 := range answer {
		for s2, v2 := range answer {
			if s1.Intersects(s2) {
				continue
			}
			best = math.Max(best, v1+v2)
		}
	}
	fmt.Printf("You and the Elephant can release %d pressure in 26m\n", best)
}

type Valve struct {
	Name     string
	FlowRate int
	Tunnels  []string
}

type Bitset uint64

func (s Bitset) Contains(i int) bool {
	return s&(1<<i) != 0
}

func (s Bitset) Add(i int) Bitset {
	return s | (1 << i)
}

func (s Bitset) Intersects(s2 Bitset) bool {
	return s&s2 != 0
}

func (s Bitset) String() string {
	return fmt.Sprintf("%b", s)
}
