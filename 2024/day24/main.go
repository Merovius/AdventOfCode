package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"unique"

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

func Parse(in []byte) (System, error) {
	return parse.Struct[System](
		split.Blocks,
		parse.Slice(
			split.Lines,
			parse.Struct[WireState](
				split.On(": "),
				ParseWire,
				func(s string) (bool, error) {
					switch s {
					case "0":
						return false, nil
					case "1":
						return true, nil
					default:
						return false, fmt.Errorf("invalid wire state %q", s)
					}
				},
			),
		),
		parse.Slice(
			split.Lines,
			parse.Struct[Gate](
				split.Regexp(`(\w+) (\w+) (\w+) -> (\w+)`),
				ParseWire,
				ParseOp,
				ParseWire,
				ParseWire,
			),
		),
	)(string(in))
}

type Wire unique.Handle[string]

func (w Wire) String() string {
	return unique.Handle[string](w).Value()
}

func ParseWire(s string) (Wire, error) {
	return Wire(unique.Make(s)), nil
}

func MakeWire(s string) Wire {
	return Wire(unique.Make(s))
}

type System struct {
	Inputs []WireState
	Gates  []Gate
}

type WireState struct {
	Wire
	State bool
}

type Gate struct {
	A  Wire
	Op Op
	B  Wire
	C  Wire
}

type Op int

const (
	AND Op = iota
	OR
	XOR
)

func (o Op) String() string {
	switch o {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	default:
		panic("invalid Op")
	}
}

func ParseOp(s string) (Op, error) {
	switch s {
	case "AND":
		return AND, nil
	case "OR":
		return OR, nil
	case "XOR":
		return XOR, nil
	default:
		return 0, fmt.Errorf("invalid gate %q", s)
	}
}

func Part1(in System) int {
	state := make(map[Wire]func() bool)
	for _, s := range in.Inputs {
		state[s.Wire] = func() bool { return s.State }
	}
	for _, g := range in.Gates {
		if _, ok := state[g.C]; ok {
			panic(fmt.Errorf("multiple gates output to %v", g.C))
		}
		state[g.C] = func() bool {
			switch g.Op {
			case AND:
				return state[g.A]() && state[g.B]()
			case OR:
				return state[g.A]() || state[g.B]()
			case XOR:
				return state[g.A]() != state[g.B]()
			default:
				panic("invalid op")
			}
		}
	}
	var out int
	for w, f := range state {
		ns, ok := strings.CutPrefix(w.String(), "z")
		if !ok {
			continue
		}
		d, err := strconv.Atoi(ns)
		if err != nil {
			panic(fmt.Errorf("can not parse output wire %q: %w", ns, err))
		}
		if d > 63 || d < 0 {
			panic(fmt.Errorf("output wire %d out of range", d))
		}
		if f() {
			out |= (1 << d)
		}
	}
	return out
}

func Part2(in System) string {
	// No programmatic solution (yet?).
	//
	// Instead, I used Dump below to output the network in graphviz format.
	// It also finds all the one-bit adders and replaces the XOR/AND gate
	// making it up with a single block that labels its carry and output.
	//
	// I then stared at the resulting graph to find the swapped wires.
	//
	// I could probably walk the graph, making sure that for all the
	// adders, their actually output/carry wire is connected to the right
	// gate. But that would take a lot of time to figure out and I don't
	// want to spend that time right now.
	return "cbd,gmh,jmq,qrh,rqf,z06,z13,z38"
}

type Adder struct {
	A Wire
	B Wire
	O Wire
	C Wire
}

func (s System) Simplify() (System, []Adder) {
	var (
		out    System
		adders []Adder
	)
	out.Inputs = slices.Clone(s.Inputs)
	adder := func(f, g Gate) (Adder, bool) {
		if f.Op > g.Op {
			f, g = g, f
		}
		if !(f.Op == AND && g.Op == XOR) {
			return Adder{}, false
		}
		if f.A.String() > f.B.String() {
			f.A, f.B = f.B, f.A
		}
		if g.A.String() > g.B.String() {
			g.A, g.B = g.B, g.A
		}
		if f.A != g.A || f.B != g.B {
			return Adder{}, false
		}
		return Adder{
			A: f.A,
			B: f.B,
			O: g.C,
			C: f.C,
		}, true
	}
	used := make(set.Set[Gate])
gates:
	for i, g := range s.Gates {
		if used.Contains(g) {
			continue
		}
		if g.Op == OR {
			out.Gates = append(out.Gates, g)
			used.Add(g)
			continue
		}
		for _, g2 := range s.Gates[i+1:] {
			if used.Contains(g2) {
				continue
			}
			if a, ok := adder(g, g2); ok {
				adders = append(adders, a)
				used.Add(g)
				used.Add(g2)
				continue gates
			}
		}
		out.Gates = append(out.Gates, g)
		used.Add(g)
	}
	return out, adders
}

func (s System) Dump() {
	s, adders := s.Simplify()

	seen := make(set.Set[Wire])
	fmt.Println("digraph G {")
	for _, w := range s.Inputs {
		fmt.Printf("\t%v;\n", w.Wire)
		seen.Add(w.Wire)
	}
	for _, g := range s.Gates {
		if !seen.Contains(g.A) {
			fmt.Printf("\t%v;\n", g.A)
			seen.Add(g.A)
		}
		if !seen.Contains(g.B) {
			fmt.Printf("\t%v;\n", g.B)
			seen.Add(g.B)
		}
		if !seen.Contains(g.C) {
			fmt.Printf("\t%v;\n", g.C)
			seen.Add(g.C)
		}
	}

	for _, g := range s.Gates {
		a, b, c := g.A, g.B, g.C

		fmt.Printf("\t%v_%v_%v [label=%q, shape=rect];\n", a, g.Op, b, g.Op.String())
		fmt.Printf("\t%v -> %v_%v_%v;\n", a, a, g.Op, b)
		fmt.Printf("\t%v -> %v_%v_%v;\n", b, a, g.Op, b)
		fmt.Printf("\t%v_%v_%v -> %v;\n", a, g.Op, b, c)
	}

	for _, a := range adders {
		name := strings.Join([]string{a.A.String(), a.B.String(), a.O.String(), a.C.String()}, "_")
		fmt.Printf("\t%s [label=\"ADD\", shape=rect];\n", name)
		fmt.Printf("\t%v -> %s;\n", a.A, name)
		fmt.Printf("\t%v -> %s;\n", a.B, name)
		fmt.Printf("\t%s:se -> %v [color=red,taillabel=C];\n", name, a.C)
		fmt.Printf("\t%s:sw -> %v [color=green,taillabel=O];\n", name, a.O)
	}
	fmt.Println("}")
}
