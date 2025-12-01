package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/math"
	"gonih.org/AdventOfCode/internal/set"
)

func main() {
	dump := flag.Bool("dump", false, "dump dot of network and exit")
	flag.Parse()
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	if *dump {
		WriteDot(os.Stdout, in)
		return
	}
	if v, err := Part1(in); err != nil {
		log.Println(err)
	} else {
		fmt.Println("Part 1:", v)
	}
	if v, err := Part2(in); err != nil {
		log.Println(err)
	} else {
		fmt.Println("Part 2:", v)
	}
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

func Part1(net Network) (int, error) {
	m, err := validate(net)
	if err != nil {
		return 0, err
	}
	type state struct {
		n string
		k int
	}
	var (
		N    int
		cur  = "AAA"
		seen = make(set.Set[state])
	)
	for {
		for _, i := range net.Instructions {
			s := state{cur, N % len(net.Instructions)}
			if seen.Contains(s) {
				return 0, errors.New("end state is unreachable")
			}
			seen.Add(s)
			if cur == "ZZZ" {
				return N, nil
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

func Part2(net Network) (int, error) {
	m, err := validate(net)
	if err != nil {
		return 0, err
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
			if len(cycles) == 0 {
				return 0, fmt.Errorf("node %q does not reach end state", n.Name)
			}
			first = false
			continue
		}
		cycles = merge(cycles, run(m, net.Instructions, n.Name))
		if len(cycles) == 0 {
			return 0, fmt.Errorf("node %q never reaches end state with other nodes", n.Name)
		}
	}
	v := math.MaxInt
	for _, c := range cycles {
		v = min(v, c.offset)
	}
	return v, nil
}

func validate(net Network) (map[string]Node, error) {
	m := make(map[string]Node)
	starts, ends := make(set.Set[string]), make(set.Set[string])
	for _, n := range net.Nodes {
		if _, ok := m[n.Name]; ok {
			return nil, fmt.Errorf("duplicate node %q", n.Name)
		}
		if strings.HasSuffix(n.Name, "A") {
			starts.Add(n.Name)
		} else if strings.HasSuffix(n.Name, "Z") {
			ends.Add(n.Name)
		}
		m[n.Name] = n
	}
	if len(starts) != len(ends) {
		return nil, fmt.Errorf("have %d start but %d end nodes", len(starts), len(ends))
	}
	for _, n := range net.Nodes {
		if _, ok := m[n.Left]; !ok {
			return nil, fmt.Errorf("node %q has undefined left neighbor %q", n.Name, n.Left)
		}
		if _, ok := m[n.Right]; !ok {
			return nil, fmt.Errorf("node %q has undefined right neighbor %q", n.Name, n.Right)
		}
	}
	if len(net.Instructions) == 0 {
		return nil, errors.New("no instructions")
	}
	return m, nil
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
		if a.offset == b.offset {
			return a, true
		}
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
	x, M, ok := math.ChineseRemainder(a.offset, a.length, b.offset, b.length)
	if !ok {
		return cycle{}, false
	}
	if m := max(a.offset, b.offset); x < m {
		x += ((m - x) / M) * M
		if x < m {
			x += M
		}
	}
	return cycle{offset: x, length: M}, true
}

func WriteDot(w io.Writer, net Network) error {
	wc := keepError(w)
	defer wc.Close()
	fmt.Fprintln(wc, "digraph G {")
	fmt.Fprintln(wc, "\tsubgraph cluster {")
	fmt.Fprintln(wc, "\tlabel = \"Network\"")
	fmt.Fprintln(wc, "\tcolor=black")
	for _, n := range net.Nodes {
		shape := "ellipse"
		if strings.HasSuffix(n.Name, "A") {
			shape = "diamond"
		} else if strings.HasSuffix(n.Name, "Z") {
			shape = "rect"
		}
		fmt.Fprintf(wc, "\t\t_%s [label=%q,shape=%s]\n", n.Name, n.Name, shape)
		fmt.Fprintf(wc, "\t\t_%s -> _%s [color=green,label=L]\n", n.Name, n.Left)
		fmt.Fprintf(wc, "\t\t_%s -> _%s [color=red,label=R]\n", n.Name, n.Right)
	}
	fmt.Fprintln(wc, "\t}")
	m := make(map[string]Node)
	for _, n := range net.Nodes {
		m[n.Name] = n
	}
	prog := net.Instructions
	type state struct {
		n string
		k int
	}
	fmt.Fprintln(wc)
	fmt.Fprintln(wc, "\tsubgraph cluster_states {")
	fmt.Fprintln(wc, "\t\tlabel = \"State machine\"")
	fmt.Fprintln(wc, "\t\tcolor = black")
	printed := make(map[any]bool)
	printState := func(s state, attr string) {
		if printed[s] {
			return
		}
		printed[s] = true
		shape := "ellipse"
		if strings.HasSuffix(s.n, "A") {
			shape = "diamond"
		} else if strings.HasSuffix(s.n, "Z") {
			shape = "rect"
		}
		fmt.Fprintf(wc,
			"\t\t_walk_%s_%d [shape=%s,label=<%s<br/>%s<font color=\"dodgerblue\"><b>%s</b></font>%s>%s]\n",
			s.n,
			s.k,
			shape,
			s.n,
			string(prog[:s.k]),
			string(prog[s.k:s.k+1]),
			string(prog[s.k+1:]),
			attr,
		)
	}
	printEdge := func(from, to state, attr string) {
		type edge struct {
			from state
			to   state
			attr string
		}
		e := edge{from, to, attr}
		if printed[e] {
			return
		}
		printed[e] = true
		fmt.Fprintf(wc,
			"\t\t_walk_%s_%d -> _walk_%s_%d [%s]\n",
			from.n,
			from.k,
			to.n,
			to.k,
			attr,
		)
	}
	for _, n := range net.Nodes {
		if !strings.HasSuffix(n.Name, "A") {
			continue
		}
		var (
			N    int
			seen = make(map[state]int)
		)
	loop:
		for {
			for _, i := range prog {
				j := N % len(prog)
				s := state{n.Name, j}
				if _, ok := seen[s]; ok {
					break loop
				}
				seen[s] = N
				printState(s, "")
				switch i {
				case L:
					printEdge(s, state{n.Left, (j + 1) % len(prog)}, "label=L")
					n = m[n.Left]
				case R:
					printEdge(s, state{n.Right, (j + 1) % len(prog)}, "label=R")
					n = m[n.Right]
				}
				N++
			}
		}
	}
	fmt.Fprintln(wc, "\t}")
	fmt.Fprintln(wc, "}")
	return wc.Close()
}

type errWriter struct {
	w   *bufio.Writer
	err error
}

func keepError(w io.Writer) io.WriteCloser {
	return &errWriter{w: bufio.NewWriter(w)}
}

func (w errWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	n, w.err = w.w.Write(p)
	return n, w.err
}

func (w errWriter) Close() error {
	if w.err != nil {
		return w.err
	}
	return w.w.Flush()
}
