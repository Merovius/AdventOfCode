package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	log.SetFlags(log.Lshortfile)

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

func Parse(in []byte) ([]Gate, error) {
	return parse.Slice(
		split.Lines,
		parse.Struct[Gate](
			split.On(" -> "),
			parse.Any[Value](
				parse.Prefix(
					"NOT ",
					parse.Struct[Not](
						split.Nop,
						parse.Any[Value](ParseID, ParseConst),
					),
				),
				parse.Struct[BinOp](
					split.Regexp(`(\w+) ([A-Z]+) (\w+)`),
					parse.Any[Value](ParseID, ParseConst),
					parse.String[string],
					parse.Any[Value](ParseID, ParseConst),
				),
				ParseID,
				ParseConst,
			),
			ParseID,
		),
	)(string(in))
}

type State map[ID]Value

type Gate struct {
	Val Value
	Dst ID
}

func (g Gate) String() string {
	return g.Val.String() + " -> " + g.Dst.String()
}

type Value interface {
	Load(State) uint16
	String() string
}

type ID string

func ParseID(s string) (ID, error) {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return "", errors.New("not an ID")
		}
	}
	return ID(s), nil
}

func (i ID) Load(s State) uint16 {
	// memoize answer
	v := s[i].Load(s)
	s[i] = Const(v)
	return v
}

func (i ID) String() string {
	return string(i)
}

type Const uint16

func ParseConst(s string) (Const, error) {
	v, err := strconv.Atoi(s)
	return Const(v), err
}

func (c Const) Load(s State) uint16 {
	return uint16(c)
}

func (c Const) String() string {
	return strconv.Itoa(int(c))
}

type Not struct {
	Val Value
}

func (n Not) Load(s State) uint16 {
	return ^n.Val.Load(s)
}

func (n Not) String() string {
	return "NOT " + n.Val.String()
}

type BinOp struct {
	Left  Value
	Op    string
	Right Value
}

func (o BinOp) Load(s State) uint16 {
	l := o.Left.Load(s)
	r := o.Right.Load(s)
	switch o.Op {
	case "AND":
		return l & r
	case "OR":
		return l | r
	case "LSHIFT":
		return l << r
	case "RSHIFT":
		return l >> r
	default:
		panic("invalid operator")
	}
}

func (o BinOp) String() string {
	return o.Left.String() + " " + o.Op + " " + o.Right.String()
}

func Part1(in []Gate) uint16 {
	r := make(map[ID]Value)
	for _, g := range in {
		r[g.Dst] = g.Val
	}
	return r["a"].Load(r)
}

func Part2(in []Gate) uint16 {
	a := Part1(in)
	r := make(map[ID]Value)
	for _, g := range in {
		r[g.Dst] = g.Val
	}
	r["b"] = Const(a)
	return r["a"].Load(r)
}
