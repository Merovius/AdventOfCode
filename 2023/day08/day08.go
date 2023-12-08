package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(s string) (Network, error) {
	return parse.TrimSpace(parse.Struct[Network](
		split.Blocks,
		parse.Slice(split.Bytes, parse.Enum(L, R)),
		parse.Lines(
			parse.Struct[Node](
				split.Regexp(`(\w+)\s+=\s+\((\w+),\s+(\w+)\)`),
				parse.String[string],
				parse.String[string],
				parse.String[string],
			),
		),
	))(s)
}

type Network struct {
	Instructions []Inst
	Nodes        []Node
}

type Inst rune

const (
	L Inst = 'L'
	R Inst = 'R'
)

type Node struct {
	Name  string
	Left  string
	Right string
}

func Part1(net Network) int {
	m := make(map[string]Node)
	for _, n := range net.Nodes {
		m[n.Name] = n
	}
	var (
		N   int
		cur = "AAA"
	)
	for {
		for _, i := range net.Instructions {
			if cur == "ZZZ" {
				return N
			}
			switch i {
			case L:
				cur = m[cur].Left
			case R:
				cur = m[cur].Right
			default:
				panic(fmt.Errorf("invalid Inst %q", i))
			}
			N++
		}
	}
}

func Part2(net Network) int {
	m := make(map[string]Node)
	for _, n := range net.Nodes {
		m[n.Name] = n
	}
	var (
		cycles []cycle
		first  = true
	)
	for _, n := range net.Nodes {
		if !strings.HasSuffix(n.Name, "A") {
			continue
		}
		if first {
			cycles = run(m, net.Instructions, n.Name)
			first = false
			continue
		}
		cycles = merge(cycles, run(m, net.Instructions, n.Name))
	}
	if len(cycles) == 0 {
		panic("no solution")
	}
	v := math.MaxInt
	for _, c := range cycles {
		v = min(v, c.offset)
	}
	return v
}

func run(m map[string]Node, prog []Inst, node string) []cycle {
	type state struct {
		n string
		k int // mod len(prog)
	}
	var (
		N     int                   // current step
		C     int                   // start of cycle
		T     int                   // length of cycle
		seen  = make(map[state]int) // maps state to step
		goals []int                 // all goals we have found
	)
loop:
	for {
		for _, i := range prog {
			s := state{node, N % len(prog)}
			if k, ok := seen[s]; ok {
				C, T = k, N-k
				break loop
			}
			seen[s] = N
			if strings.HasSuffix(node, "Z") {
				goals = append(goals, N)
			}
			switch i {
			case L:
				node = m[node].Left
			case R:
				node = m[node].Right
			}
			N++
		}
	}
	var cycles []cycle
	for _, g := range goals {
		c := cycle{offset: g}
		if g >= C {
			c.length = T
		}
		cycles = append(cycles, c)
	}
	return cycles
}

type cycle struct {
	offset int
	length int
}

func (c cycle) String() string {
	return fmt.Sprintf("(%d+%d)", c.offset, c.length)
}

func merge(a, b []cycle) []cycle {
	out := []cycle{}
	for _, ca := range a {
		for _, cb := range b {
			if c, ok := mergeCycles(ca, cb); ok {
				out = append(out, c)
			}
		}
	}
	return out
}

func mergeCycles(a, b cycle) (cycle, bool) {
	if a.length == 0 && b.length == 0 {
		return cycle{}, false
	}
	if a.length == 0 {
		if a.offset < b.offset {
			return cycle{}, false
		}
		if (a.offset-b.offset)%b.length != 0 {
			return cycle{}, false
		}
		return a, true
	}
	if b.length == 0 {
		return mergeCycles(b, a)
	}
	x, ok := math.CRT(a.offset%a.length, b.offset%b.length, a.length, b.length)
	if !ok {
		return cycle{}, false
	}
	M := math.LCM(a.length, b.length)
	if m := max(a.offset, b.offset); x < m {
		x += ((m - x) / M) * M
		if x < m {
			x += M
		}
	}
	return cycle{offset: x, length: M}, true
}
