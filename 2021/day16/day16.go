package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

func main() {
	read(os.Stdin)
}

func read(r io.Reader) {
	buf, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	buf = bytes.TrimSpace(buf)
	p, err := ParsePacket(string(buf))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sum of versions: %d\n", p.VersionSum())
	fmt.Printf("Evaluated: %d\n", p.Eval())
}

func ParsePacket(s string) (Packet, error) {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	bits := make([]bool, 0, 8*len(buf))
	for _, b := range buf {
		for i := 0; i < 8; i++ {
			bits = append(bits, b&0x80 != 0)
			b <<= 1
		}
	}
	return parsePacket(bits)
}

type parser struct {
	bits  []bool
	check func(error)
}

func parsePacket(bits []bool) (pk Packet, err error) {
	sentinel := new(uint8)
	defer func() {
		if v := recover(); v != sentinel && v != nil {
			panic(v)
		}
	}()
	check := func(e error) {
		if e != nil {
			err = e
			panic(sentinel)
		}
	}
	p := &parser{bits: bits, check: check}
	return p.parse(), nil
}

func (p *parser) abort(format string, args ...interface{}) {
	p.check(fmt.Errorf(format, args...))
}

func (p *parser) parse() Packet {
	version := p.int(3)
	id := p.int(3)
	if id == 4 {
		return p.literal(version)
	}
	return p.operator(version, id)
}

func (p *parser) literal(ver int) Packet {
	var (
		v    int
		cont = true
	)
	for cont {
		cont = p.bool()
		v <<= 4
		v |= p.int(4)
	}
	return Literal{Ver: ver, Val: v}
}

func (p *parser) operator(ver, id int) Packet {
	var sub []Packet
	if p.bool() {
		m := p.int(11)
		for i := 0; i < m; i++ {
			sub = append(sub, p.parse())
		}
	} else {
		m := p.int(15)
		m = len(p.bits) - m
		for len(p.bits) > m {
			sub = append(sub, p.parse())
		}
	}
	return Operator{Ver: ver, ID: id, Sub: sub}
}

func (p *parser) consume(n int) (v []bool) {
	if len(p.bits) < n {
		p.abort("not enough bits, expected %d, have %d", n, len(p.bits))
	}
	v, p.bits = p.bits[:n], p.bits[n:]
	return v
}

func (p *parser) int(n int) int {
	b := p.consume(n)

	var v int
	for i := 0; i < n; i++ {
		v <<= 1
		if b[i] {
			v |= 1
		}
	}
	return v
}

func (p *parser) bool() bool {
	return p.consume(1)[0]
}

type Packet interface {
	String() string
	VersionSum() int
	Expr() string
	Eval() int
}

type Literal struct {
	Ver int
	Val int
}

func (l Literal) String() string {
	return fmt.Sprintf("Literal(ver=%d,val=%d)", l.Ver, l.Val)
}

func (l Literal) Expr() string {
	return fmt.Sprintf("%d", l.Val)
}

func (l Literal) VersionSum() int {
	return l.Ver
}

func (l Literal) Eval() int {
	return l.Val
}

type Operator struct {
	Ver int
	ID  int
	Sub []Packet
}

func (o Operator) String() string {
	return fmt.Sprintf("Operator(Ver=%d, ID=%d, Sub=%v)", o.Ver, o.ID, o.Sub)
}

func (o Operator) Expr() string {
	var pieces []string
	for _, p := range o.Sub {
		pieces = append(pieces, p.Expr())
	}
	switch o.ID {
	case 0:
		return fmt.Sprintf("sum(%s)", strings.Join(pieces, ", "))
	case 1:
		return fmt.Sprintf("prod(%s)", strings.Join(pieces, ", "))
	case 2:
		return fmt.Sprintf("min(%s)", strings.Join(pieces, ", "))
	case 5:
		return fmt.Sprintf("gt(%s)", strings.Join(pieces, ", "))
	case 6:
		return fmt.Sprintf("lt(%s)", strings.Join(pieces, ", "))
	case 7:
		return fmt.Sprintf("eq(%s)", strings.Join(pieces, ", "))
	default:
		panic(fmt.Sprintf("invalid ID %d", o.ID))
	}
}

func (o Operator) VersionSum() int {
	v := o.Ver
	for _, s := range o.Sub {
		v += s.VersionSum()
	}
	return v
}

func (o Operator) Eval() (total int) {
	switch o.ID {
	case 0:
		for _, s := range o.Sub {
			total += s.Eval()
		}
	case 1:
		total = 1
		for _, s := range o.Sub {
			total *= s.Eval()
		}
	case 2:
		total = math.MaxInt
		for _, s := range o.Sub {
			if v := s.Eval(); v < total {
				total = v
			}
		}
	case 3:
		for _, s := range o.Sub {
			if v := s.Eval(); v > total {
				total = v
			}
		}
	case 5:
		a, b := o.Sub[0].Eval(), o.Sub[1].Eval()
		if a > b {
			return 1
		}
		return 0
	case 6:
		a, b := o.Sub[0].Eval(), o.Sub[1].Eval()
		if a < b {
			return 1
		}
		return 0
	case 7:
		a, b := o.Sub[0].Eval(), o.Sub[1].Eval()
		if a == b {
			return 1
		}
		return 0
	default:
		panic(fmt.Sprintf("invalid operator ID %d", o.ID))
	}
	return total
}
