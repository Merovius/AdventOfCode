package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"math"
	"os"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
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
	Dump(in.Prog)
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(buf []byte) (Machine, error) {
	type regs struct {
		A int
		B int
		C int
	}
	type input struct {
		Regs regs
		Prog []int
	}
	in, err := parse.Struct[input](
		split.Blocks,
		parse.Struct[regs](
			split.Lines,
			parse.Prefix("Register A: ", parse.Signed[int]),
			parse.Prefix("Register B: ", parse.Signed[int]),
			parse.Prefix("Register C: ", parse.Signed[int]),
		),
		parse.Prefix(
			"Program: ",
			parse.Slice(
				split.On(","),
				parse.Signed[int],
			),
		),
	)(string(buf))
	return Machine{A: in.Regs.A, B: in.Regs.B, C: in.Regs.C, Prog: in.Prog}, err
}

type Machine struct {
	A    int
	B    int
	C    int
	Prog []int
}

const (
	Adv = 0
	Bxl = 1
	Bst = 2
	Jnz = 3
	Bxc = 4
	Out = 5
	Bdv = 6
	Cdv = 7
)

func Part1(m Machine) string {
	return m.Run().String()
}

func (m Machine) Run() Output {
	var (
		pc  int
		out Output
	)
	for pc >= 0 && pc < len(m.Prog)-1 {
		arg := m.Prog[pc+1]
		switch m.Prog[pc] {
		case Adv:
			m.A >>= m.combo(arg)
		case Bxl:
			m.B ^= arg
		case Bst:
			m.B = m.combo(arg) % 8
		case Jnz:
			if m.A != 0 {
				pc = arg
				continue
			}
		case Bxc:
			m.B ^= m.C
		case Out:
			out.Write(m.combo(arg))
		case Bdv:
			m.B = m.A >> m.combo(arg)
		case Cdv:
			m.C = m.A >> m.combo(arg)
		}
		pc += 2
	}
	return out
}

func (m Machine) combo(v int) int {
	switch v {
	case 0, 1, 2, 3:
		return v
	case 4:
		return m.A
	case 5:
		return m.B
	case 6:
		return m.C
	default:
		panic("illegal instruction")
	}
}

func Part2(m Machine) int {
	best := math.MaxInt
	for v := range m.Reverse() {
		best = min(best, v)
	}
	if best < math.MaxInt {
		return best
	}
	return -1
}

func (m Machine) Reverse() iter.Seq[int] {
	// From squinting at the program, we see that it is a loop, which
	// consumes the lowest 3 bits of A, mixes them with up to 7 additional
	// bits from A and outputs them, until A is completely consumed.
	//
	// So the lowest 10 bits determine the next output.
	//
	// We brute-force all possible combinations of 10 bits, selecting all
	// which output the right first number. We then recursively try to
	// extend it by adding three more bits.
	return func(yield func(int) bool) {
		var want Output
		for _, i := range m.Prog {
			want.Write(i)
		}
		for i := 0; i < 128; i++ {
			m.A = i
			m.reverse(7, 0b111, want, yield)
		}
	}
}

func (m Machine) reverse(shift, mask int, want Output, yield func(int) bool) {
	if shift >= 60 {
		return
	}
	a := m.A
	for i := 0; i < 8; i++ {
		m.A = a | (i << shift)
		if got := m.Run(); got == want {
			if !yield(m.A) {
				return
			}
			continue
		} else if got.Mask(mask) != want.Mask(mask) {
			continue
		}
		m.reverse(shift+3, (mask<<3)|0b111, want, yield)
	}
}

// Output is the output of the machine, compressed into a single int64, with 3
// bits per number.
type Output struct {
	v int64
	n int
}

func (o *Output) Write(v int) {
	if o.n > 60 {
		panic("output overflow")
	}
	o.v |= int64(v%8) << o.n
	o.n += 3
}

func (o Output) String() string {
	w := new(strings.Builder)
	for range o.n / 3 {
		if w.Len() > 0 {
			w.WriteByte(',')
		}
		w.WriteByte('0' + byte(o.v&0b111))
		o.v >>= 3
	}
	return w.String()
}

func (o Output) Mask(m int) int {
	return int(o.v) & m
}

func Dump(prog []int) {
	combo := "0123ABC"
	for i := 0; i < len(prog); i += 2 {
		switch prog[i] {
		case 0:
			fmt.Printf("%3d  adv %v\n", i, combo[prog[i+1]])
		case 1:
			fmt.Printf("%3d  bxl %v\n", i, prog[i+1])
		case 2:
			fmt.Printf("%3d  bst %v\n", i, combo[prog[i+1]])
		case 3:
			fmt.Printf("%3d  jnz %v\n", i, prog[i+1])
		case 4:
			fmt.Printf("%3d  bxc\n", i)
		case 5:
			fmt.Printf("%3d  out %v\n", i, combo[prog[i+1]])
		case 6:
			fmt.Printf("%3d  bdv %v\n", i, combo[prog[i+1]])
		case 7:
			fmt.Printf("%3d  cdv %v\n", i, combo[prog[i+1]])
		}
	}
}
