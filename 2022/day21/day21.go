package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if *debug {
		log.SetFlags(log.Lshortfile)
		logf = log.Printf
	}
	data, err := input.Slice(input.Lines(), func(s string) (Decl, error) {
		name, expr, ok := strings.Cut(s, ": ")
		if !ok {
			return Decl{}, fmt.Errorf("invalid line %q", s)
		}
		v, err := strconv.Atoi(expr)
		if err == nil {
			return Decl{name, Const(v)}, nil
		}
		if len(expr) != 11 {
			return Decl{}, fmt.Errorf("invalid line %q", s)
		}
		a1, a2 := Var(expr[:4]), Var(expr[7:])
		switch expr[4:7] {
		case " + ":
			return Decl{name, Binary{OpAdd, a1, a2}}, nil
		case " - ":
			return Decl{name, Binary{OpSub, a1, a2}}, nil
		case " * ":
			return Decl{name, Binary{OpMul, a1, a2}}, nil
		case " / ":
			return Decl{name, Binary{OpDiv, a1, a2}}, nil
		default:
			return Decl{}, fmt.Errorf("invalid line %q", s)
		}
	}).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("root yells %v\n", Eval(data))
	fmt.Printf("You need to yell %v\n", Solve(data))
}

var logf = func(format string, args ...any) {}

type Expr interface {
	Eval(m map[string]Expr) int
	Simpl(m map[string]Expr) Expr
	Solve(v int) (int, error)
	String() string
}

type Decl struct {
	Name string
	Expr Expr
}

type Const int

func (c Const) Eval(m map[string]Expr) int {
	return int(c)
}

func (c Const) Simpl(m map[string]Expr) Expr {
	return c
}

func (c Const) Solve(v int) (int, error) {
	if int(c) != v {
		return 0, fmt.Errorf("constant %d can not be %d", c, v)
	}
	return v, nil
}

func (c Const) String() string {
	return strconv.Itoa(int(c))
}

type Var string

func (v Var) Eval(m map[string]Expr) int {
	return m[string(v)].Eval(m)
}

func (v Var) Simpl(m map[string]Expr) (e Expr) {
	if v == "humn" {
		return v
	}
	return m[string(v)].Simpl(m)
}

func (v Var) Solve(w int) (int, error) {
	logf("Solve(%s == %d)", v, w)
	return w, nil
}

func (v Var) String() string {
	return string(v)
}

type Op int

const (
	_ Op = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpEq
)

func (o Op) String() string {
	switch o {
	case OpAdd:
		return "+"
	case OpSub:
		return "-"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpEq:
		return "="
	default:
		panic("invalid op")
	}
}

type Binary struct {
	Op   Op
	Arg1 Expr
	Arg2 Expr
}

func (b Binary) Eval(m map[string]Expr) (r int) {
	a1 := b.Arg1.Eval(m)
	a2 := b.Arg2.Eval(m)
	switch b.Op {
	case OpAdd:
		return a1 + a2
	case OpSub:
		return a1 - a2
	case OpMul:
		return a1 * a2
	case OpDiv:
		return a1 / a2
	case OpEq:
		panic("OpEq can not be evaluated")
	default:
		panic("invalid op")
	}
}

func (b Binary) Simpl(m map[string]Expr) (e Expr) {
	a1 := b.Arg1.Simpl(m)
	a2 := b.Arg2.Simpl(m)
	c1, ok1 := a1.(Const)
	c2, ok2 := a2.(Const)
	if ok1 && ok2 && b.Op != OpEq {
		return Const(Binary{b.Op, c1, c2}.Eval(m))
	}
	return Binary{b.Op, a1, a2}
}

func (b Binary) Solve(v int) (w int, err error) {
	c1, ok1 := b.Arg1.(Const)
	c2, ok2 := b.Arg2.(Const)
	if ok1 == ok2 {
		return 0, errors.New("binary expression needs exactly one constant to be solved")
	}
	if ok1 {
		logf("Solve(%v %v x == %v)", c1, b.Op, v)
	} else {
		logf("Solve(x %v %v == %v)", b.Op, c2, v)
	}

	switch b.Op {
	case OpAdd:
		if ok1 {
			return b.Arg2.Solve(v - int(c1))
		}
		return b.Arg1.Solve(v - int(c2))
	case OpSub:
		if ok1 {
			return b.Arg2.Solve(int(c1) - v)
		}
		return b.Arg1.Solve(v + int(c2))
	case OpMul:
		if ok1 {
			if v%int(c1) != 0 {
				return 0, fmt.Errorf("%d is not divisible by %d", v, c1)
			}
			return b.Arg2.Solve(v / int(c1))
		}
		if v%int(c2) != 0 {
			return 0, fmt.Errorf("%d is not divisible by %d", v, c2)
		}
		return b.Arg1.Solve(v / int(c2))
	case OpDiv:
		if ok1 {
			if v%int(c1) != 0 {
				panic(fmt.Errorf("%d is not divisible by %d", v, c1))
			}
			return b.Arg2.Solve(int(c1) / v)
		}
		return b.Arg1.Solve(int(c2) * v)
	case OpEq:
		if v != 1 {
			return 0, errors.New("OpEq can only be solved for 1")
		}
		if ok1 {
			return b.Arg2.Solve(int(c1))
		}
		return b.Arg1.Solve(int(c2))
	default:
		panic("invalid op")
	}
}

func (b Binary) String() string {
	return fmt.Sprintf("(%v %v %v)", b.Arg1, b.Op, b.Arg2)
}

func Eval(d []Decl) int {
	m := make(map[string]Expr)
	for _, d := range d {
		m[d.Name] = d.Expr
	}
	return m["root"].Eval(m)
}

func Solve(d []Decl) int {
	m := make(map[string]Expr)
	for _, d := range d {
		m[d.Name] = d.Expr
	}
	e := m["root"].Simpl(m).(Binary)
	e.Op = OpEq
	v, err := e.Solve(1)
	if err != nil {
		panic(err)
	}
	return v
}
