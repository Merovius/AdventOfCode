package main

import (
	"flag"
	"fmt"

	"golang.org/x/exp/slices"
)

func main() {
	inp := flag.Bool("input", false, "use the real input")
	flag.Parse()
	data := example
	if *inp {
		data = input
	}
	var data2 []*Monkey
	for _, m := range data {
		m2 := *m
		m2.Items = append([]int(nil), m2.Items...)
		data2 = append(data2, &m2)
	}
	fmt.Printf("Monkey business part 1: %d\n", MonkeyBusiness(data, false))
	fmt.Printf("Monkey business part 2: %d\n", MonkeyBusiness(data2, true))
}

var example = []*Monkey{
	&Monkey{
		[]int{79, 98},
		OpMul,
		19,
		23,
		2,
		3,
	},
	&Monkey{
		[]int{54, 65, 75, 74},
		OpAdd,
		6,
		19,
		2,
		0,
	},
	&Monkey{
		[]int{79, 60, 97},
		OpSquare,
		0,
		13,
		1,
		3,
	},
	&Monkey{
		[]int{74},
		OpAdd,
		3,
		17,
		0,
		1,
	},
}

var input = []*Monkey{
	&Monkey{
		[]int{56, 56, 92, 65, 71, 61, 79},
		OpMul,
		7,
		3,
		3,
		7,
	},
	&Monkey{
		[]int{61, 85},
		OpAdd,
		5,
		11,
		6,
		4,
	},
	&Monkey{
		[]int{54, 96, 82, 78, 69},
		OpSquare,
		0,
		7,
		0,
		7,
	},
	&Monkey{
		[]int{57, 59, 65, 95},
		OpAdd,
		4,
		2,
		5,
		1,
	},
	&Monkey{
		[]int{62, 67, 80},
		OpMul,
		17,
		19,
		2,
		6,
	},
	&Monkey{
		[]int{91},
		OpAdd,
		7,
		5,
		1,
		4,
	},
	&Monkey{
		[]int{79, 83, 64, 52, 77, 56, 63, 92},
		OpAdd,
		6,
		17,
		2,
		0,
	},
	&Monkey{
		[]int{50, 97, 76, 96, 80, 56},
		OpAdd,
		3,
		13,
		3,
		5,
	},
}

type Monkey struct {
	Items    []int
	Operator Op
	Operand  int
	CheckArg int
	Then     int
	Else     int
}

type Op int

const (
	OpAdd Op = iota
	OpMul
	OpSquare
)

func MonkeyBusiness(ms []*Monkey, part2 bool) int {
	insp := make([]int, len(ms))
	mod := 1
	for _, m := range ms {
		mod *= m.CheckArg
	}
	rounds := 20
	if part2 {
		rounds = 10000
	}
	for i := 0; i < rounds; i++ {
		for j, m := range ms {
			insp[j] += len(m.Items)
			for _, it := range m.Items {
				switch m.Operator {
				case OpAdd:
					it += m.Operand
				case OpMul:
					it *= m.Operand
				case OpSquare:
					it *= it
				default:
					panic("invalid operation")
				}
				if part2 {
					it %= mod
				} else {
					it /= 3
				}
				if it%m.CheckArg == 0 {
					ms[m.Then].Items = append(ms[m.Then].Items, it)
				} else {
					ms[m.Else].Items = append(ms[m.Else].Items, it)
				}
			}
			m.Items = m.Items[:0]
		}
	}
	fmt.Println(insp)
	slices.Sort(insp)
	return insp[len(insp)-1] * insp[len(insp)-2]
}
