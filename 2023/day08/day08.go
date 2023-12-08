package main

import (
	"bytes"
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
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(b []byte) (Network, error) {
	return parse.Struct[Network](
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
	)(string(bytes.TrimSpace(b)))
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
	var steps = make(map[string]int)

	m := make(map[string]Node)
	for _, n := range net.Nodes {
		m[n.Name] = n
		if strings.HasSuffix(n.Name, "A") {
			steps[n.Name] = 0
		}
	}
	for n := range steps {
		steps[n] = run(m, net.Instructions, n)
	}

	// output is the least common multiple of all path-lengths
	N := 1
	for _, s := range steps {
		N = lcm(N, s)
	}
	return N
}

func run(m map[string]Node, prog []Inst, node string) int {
	var N int
	for {
		for _, i := range prog {
			if strings.HasSuffix(node, "Z") {
				// Note: Theoretically, we could reach multiple goal nodes from
				// a single start. It doesn't happen in my input, though.
				return N
			}
			switch i {
			case L:
				node = m[node].Left
			case R:
				node = m[node].Right
			default:
				panic(fmt.Errorf("invalid Inst %q", i))
			}
			N++
		}
	}
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

func lcm(a, b int) int {
	g, _, _ := math.GCD(a, b)
	return (a / g) * b
}
