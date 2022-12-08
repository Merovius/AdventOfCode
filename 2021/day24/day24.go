package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	dump := flag.Bool("dump", false, "dump a graph of the computation in graphViz format")
	flag.Parse()

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	prog, err := read(f)
	if err != nil {
		log.Fatal(err)
	}
	g := flowGraph(prog)
	g.optimize()
	if *dump {
		g.dump()
		return
	}

	fmt.Println("Try inputs. Outputs will be given as base26 decoded vectors")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		s := s.Text()
		if len(s) != 14 {
			log.Printf("Wrong input length")
			continue
		}

		v := g.vars[varZ].eval(s)
		var vec []int
		for v > 0 {
			vec = append(vec, v%26)
			v /= 26
		}
		for i := 0; i < len(vec)/2; i++ {
			j := len(vec) - 1 - i
			vec[i], vec[j] = vec[j], vec[i]
		}
		fmt.Println(vec)
	}
}

func read(r io.Reader) ([]instruction, error) {
	var out []instruction
	s := bufio.NewScanner(r)
	for s.Scan() {
		f := strings.Fields(s.Text())
		var i instruction
		switch f[0] {
		case "inp":
			i.op = opInp
		case "add":
			i.op = opAdd
		case "mul":
			i.op = opMul
		case "div":
			i.op = opDiv
		case "mod":
			i.op = opMod
		case "eql":
			i.op = opEql
		default:
			return nil, fmt.Errorf("unknown instruction %q", f[0])
		}
		var err error
		i.arg1, err = parseVar(f[1])
		if err != nil {
			return nil, err
		}
		if i.op != opInp {
			i.arg2, err = parseArg(f[2])
			if err != nil {
				return nil, err
			}
		}
		out = append(out, i)
	}
	return out, s.Err()
}

type instruction struct {
	op   op
	arg1 _var
	arg2 arg
}

func (i instruction) String() string {
	switch i.op {
	case opInp:
		return fmt.Sprintf("inp %v", i.arg1)
	case opAdd:
		return fmt.Sprintf("add %v %v", i.arg1, i.arg2)
	case opMul:
		return fmt.Sprintf("mul %v %v", i.arg1, i.arg2)
	case opDiv:
		return fmt.Sprintf("div %v %v", i.arg1, i.arg2)
	case opMod:
		return fmt.Sprintf("mod %v %v", i.arg1, i.arg2)
	case opEql:
		return fmt.Sprintf("eql %v %v", i.arg1, i.arg2)
	case opNeq:
		return fmt.Sprintf("neq %v %v", i.arg1, i.arg2)
	default:
		return fmt.Sprintf("op(%d) %v %v", i.op, i.arg1, i.arg2)
	}
}

type op int

const (
	opInvalid op = iota
	opInp
	opAdd
	opMul
	opDiv
	opMod
	opEql
	opNeq
)

func (n op) eval(l, r int) int {
	switch n {
	case opAdd:
		return l + r
	case opMul:
		return l * r
	case opDiv:
		if r == 0 {
			panic("invalid div")
		}
		return l / r
	case opMod:
		if l < 0 || r <= 0 {
			panic("invalid mod")
		}
		return l % r
	case opEql:
		if l == r {
			return 1
		}
		return 0
	case opNeq:
		if l != r {
			return 1
		}
		return 0
	default:
		panic(fmt.Sprintf("invalid op code %d", int(n)))
	}
}

func (o op) String() string {
	switch o {
	case opInp:
		return "input"
	case opAdd:
		return "+"
	case opMul:
		return "*"
	case opDiv:
		return "/"
	case opMod:
		return "%"
	case opEql:
		return "=="
	case opNeq:
		return "!="
	default:
		panic(fmt.Sprintf("invalid op code %d", int(o)))
	}
}

type _var int

const (
	varW _var = iota
	varX
	varY
	varZ
	nVars
)

func (v _var) String() string {
	switch v {
	case varW:
		return "w"
	case varX:
		return "x"
	case varY:
		return "y"
	case varZ:
		return "z"
	default:
		return "var(" + strconv.Itoa(int(v)) + ")"
	}
}

func parseVar(s string) (_var, error) {
	switch s {
	case "w":
		return varW, nil
	case "x":
		return varX, nil
	case "y":
		return varY, nil
	case "z":
		return varZ, nil
	default:
		return nVars, fmt.Errorf("unknown var %q", s)
	}
}

type arg interface{}

func parseArg(s string) (arg, error) {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i, nil
	}
	return parseVar(s)
}

const nInputs = 14

type graph struct {
	vars [nVars]node
}

func flowGraph(prog []instruction) *graph {
	g := new(graph)
	for i := _var(0); i < nVars; i++ {
		g.vars[i] = newConstNode(0)
	}
	var inp int
	for _, inst := range prog {
		var n node
		if inst.op == opInp {
			n = newInputNode(inp)
			inp++
		} else {
			var arg node
			if i, ok := inst.arg2.(int); ok {
				arg = newConstNode(i)
			} else {
				arg = g.vars[inst.arg2.(_var)]
			}
			n = &opNode{inst.op, g.vars[inst.arg1], arg}
		}
		g.vars[inst.arg1] = n
	}
	return g
}

func (g *graph) dump() {
	ids := make(map[node]int)
	fmt.Println("digraph G {")
	for i, n := range g.vars {
		g.dumpNode(n, ids)
		fmt.Printf("\t%q -> %d\n", _var(i), ids[n])
	}
	fmt.Println("}")
}

func (g *graph) dumpNode(n node, ids map[node]int) {
	id, ok := ids[n]
	if !ok {
		id = len(ids)
		ids[n] = id

		var label string
		switch n := n.(type) {
		case *constNode:
			label = strconv.Itoa(int(*n))
		case *inputNode:
			label = n.String()
		case *opNode:
			label = n.op().String()
			g.dumpNode(n.left, ids)
			fmt.Printf("\t%d -> %d [label=l]\n", id, ids[n.left])
			g.dumpNode(n.right, ids)
			fmt.Printf("\t%d -> %d [label=r]\n", id, ids[n.right])
		default:
			panic(fmt.Sprintf("unknown node type %T", n))
		}
		var min string
		if m := n.min(); m == math.MinInt {
			min = "-∞"
		} else {
			min = strconv.Itoa(m)
		}
		var max string
		if m := n.max(); m == math.MaxInt {
			max = "∞"
		} else {
			max = strconv.Itoa(m)
		}
		label += fmt.Sprintf(" [%s,%s]", min, max)
		fmt.Printf("\t%d [label=%q]\n", id, label)
	}
}

func (g *graph) optimize() {
	o := newOptimizer()
	for i, n := range g.vars {
		g.vars[i] = o.optimize(n)
	}
}

type kind int

const (
	kindConst kind = iota
	kindInput
	kindOp
)

type optimizer struct {
}

func newOptimizer() *optimizer {
	return &optimizer{}
}

func (o *optimizer) optimize(n node) node {
	on, ok := n.(*opNode)
	if !ok {
		return n
	}
	on.left = o.optimize(on.left)
	on.right = o.optimize(on.right)

	if on.left.kind() == kindConst && on.right.kind() == kindConst {
		return newConstNode(on.op().eval(on.left.val(), on.right.val()))
	}

	switch on.op() {
	case opAdd:
		if on.left.kind() == kindConst && on.left.val() == 0 {
			return on.right
		}
		if on.right.kind() == kindConst && on.right.val() == 0 {
			return on.left
		}
		if on.right.kind() == kindConst && on.left.op() == opAdd && on.left.(*opNode).right.kind() == kindConst {
			c := on.right.val()
			c += on.left.(*opNode).right.val()
			on.left = on.left.(*opNode).left
			*on.right.(*constNode) = constNode(c)
		}
	case opMul:
		if on.left.kind() == kindConst {
			switch on.left.val() {
			case 0:
				return newConstNode(0)
			case 1:
				return on.right
			}
		}
		if on.right.kind() == kindConst {
			switch on.right.val() {
			case 0:
				return newConstNode(0)
			case 1:
				return on.left
			}
		}
	case opDiv:
		if on.right.kind() == kindConst && on.right.val() == 1 {
			return on.left
		}
		if on.right.kind() == kindConst && on.left.kind() == kindOp {
			base := on.right.val()
			add, ok := on.left.(*opNode)
			if !ok || add.op() != opAdd {
				return n
			}
			mul, ok := add.left.(*opNode)
			if !ok || mul.op() != opMul {
				return n
			}
			c, ok := mul.right.(*constNode)
			if !ok || int(*c) != base {
				return n
			}
			return mul.left
		}
	case opMod:
		if on.left.min() >= 0 && on.left.max() < on.right.min() {
			return on.left
		}
		if on.right.kind() == kindConst && on.left.kind() == kindOp {
			base := on.right.val()
			add, ok := on.left.(*opNode)
			if !ok || add.op() != opAdd {
				return n
			}
			mul, ok := add.left.(*opNode)
			if !ok || mul.op() != opMul {
				return n
			}
			c, ok := mul.right.(*constNode)
			if !ok || int(*c) != base {
				return n
			}
			if int(*c) != base {
				return n
			}
			return add.right
		}
	case opEql:
		if on.left.max() < on.right.min() || on.left.min() > on.right.max() {
			return newConstNode(0)
		}
		if on.left.kind() == kindOp && on.right.kind() == kindConst && on.right.val() == 0 {
			ol := on.left.(*opNode)
			on.o = opNeq
			on.left = ol.left
			on.right = ol.right
			break
		}
		if on.right.kind() == kindOp && on.left.kind() == kindConst && on.left.val() == 0 {
			or := on.right.(*opNode)
			on.o = opNeq
			on.left = or.left
			on.right = or.right
			break
		}
	}
	return n
}

type node interface {
	fmt.Stringer
	eval(string) int
	kind() kind
	min() int
	max() int
	val() int
	op() op
}

type constNode int

func newConstNode(v int) node {
	return (*constNode)(&v)
}

func (n *constNode) String() string {
	return fmt.Sprintf("const(%d)", int(*n))
}

func (n *constNode) eval(_ string) int {
	return int(*n)
}

func (n *constNode) kind() kind {
	return kindConst
}

func (n *constNode) min() int {
	return int(*n)
}

func (n *constNode) max() int {
	return int(*n)
}

func (n *constNode) val() int {
	return int(*n)
}

func (n *constNode) op() op {
	return opInvalid
}

type inputNode int

func newInputNode(i int) node {
	return (*inputNode)(&i)
}

func (n *inputNode) String() string {
	return fmt.Sprintf("input[%d]", int(*n))
}

func (n *inputNode) eval(s string) int {
	if s[*n] < '0' || s[*n] > '9' {
		panic(fmt.Sprintf("invalid byte %q in input", s[*n]))
	}
	return int(s[*n] - '0')
}

func (n *inputNode) kind() kind {
	return kindInput
}

func (n *inputNode) min() int {
	return 1
}

func (n *inputNode) max() int {
	return 9
}

func (n *inputNode) val() int {
	panic("val on inputNode")
}

func (n *inputNode) op() op {
	return opInvalid
}

type opNode struct {
	o     op
	left  node
	right node
}

func (n opNode) String() string {
	return fmt.Sprintf("(%T %v %T)", n.left, n.o, n.right)
}

func (n *opNode) eval(s string) int {
	return n.o.eval(n.left.eval(s), n.right.eval(s))
}

func (n *opNode) kind() kind {
	return kindOp
}

func (n *opNode) min() int {
	switch n.o {
	case opAdd:
		return n.left.min() + n.right.min()
	case opMul:
		min := math.MaxInt
		for _, a := range []int{n.left.min(), n.left.max()} {
			for _, b := range []int{n.right.min(), n.right.max()} {
				if a*b < min {
					min = a * b
				}
			}
		}
		return min
	case opDiv:
		if n.left.min() < 0 || n.right.min() < 0 {
			panic("can't determine ranges for division including negative numbers")
		}
		return n.left.min() / n.right.max()
	case opMod:
		return 0
	case opEql:
		return 0
	case opNeq:
		return 1
	default:
		panic(fmt.Sprintf("unknown op %v", n.o))
	}
}

func (n *opNode) max() int {
	switch n.o {
	case opAdd:
		return n.left.max() + n.right.max()
	case opMul:
		max := math.MinInt
		for _, a := range []int{n.left.min(), n.left.max()} {
			for _, b := range []int{n.right.min(), n.right.max()} {
				if a*b > max {
					max = a * b
				}
			}
		}
		return max
	case opDiv:
		if n.left.min() < 0 || n.right.min() < 0 {
			panic("can't determine ranges for division including negative numbers")
		}
		return n.left.max() / n.right.min()
	case opMod:
		max := n.right.max() - 1
		if m := n.left.max(); m < max {
			max = m
		}
		return max
	case opEql:
		return 1
	case opNeq:
		return 0
	default:
		panic(fmt.Sprintf("unknown op %v", n.o))
	}
}

func (n *opNode) val() int {
	panic("val on opNode")
}

func (n *opNode) op() op {
	return n.o
}

func toStackMachine(g *graph) []stackInst {
	var prog []stackInst

	var rec func(n node)
	rec = func(n node) {

		switch n := n.(type) {
		case *inputNode:
			prog = append(prog, stackInst{stackRead, 0})
		case *constNode:
			prog = append(prog, stackInst{stackConst, n.val()})
		case *opNode:
			if n.op() == opMul && n.right.kind() == kindConst && n.right.val() == 26 {
				rec(n.left)
				prog = append(prog, stackInst{stackConst, 0})
				return
			}
			rec(n.left)
			rec(n.right)
			if n.op() == opAdd {
				prog = append(prog, stackInst{stackAdd, 0})
			}
		}
	}
	rec(g.vars[varZ])

	return prog
}

type stackInst struct {
	op  stackOp
	arg int
}

func (i stackInst) String() string {
	switch i.op {
	case stackRead:
		return "read"
	case stackConst:
		return fmt.Sprintf("push %d", i.arg)
	case stackAdd:
		return fmt.Sprintf("add")
	default:
		return fmt.Sprintf("invalid(%d, %d)", i.op, i.arg)
	}
}

type stackOp int

const (
	stackInvalid stackOp = iota
	stackRead
	stackConst
	stackAdd
)
