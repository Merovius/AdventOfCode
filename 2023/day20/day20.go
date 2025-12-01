package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/container"
	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/math"
)

func main() {
	dump := flag.Bool("dump", false, "dump network in dot format and exit")
	flag.Parse()

	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	if *dump {
		DumpGraph(in)
		return
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(s string) ([]Module, error) {
	return parse.TrimSpace(parse.Lines(
		parse.Struct[Module](
			split.Regexp(`([%&]?)(\w+) -> (.*)`),
			func(s string) (Kind, error) {
				switch s {
				case "%":
					return FlipFlop, nil
				case "&":
					return Conjunction, nil
				default:
					return Broadcaster, nil
				}
			},
			parse.String[string],
			parse.Slice(
				split.On(", "),
				parse.String[string],
			),
		),
	))(s)
}

func Part1(in []Module) int {
	mods := make(map[string]Module)
	conj := make(map[string]map[string]bool)
	for _, m := range in {
		mods[m.Name] = m
		if m.Kind == Conjunction {
			conj[m.Name] = make(map[string]bool)
		}
	}
	for _, m := range in {
		for _, d := range m.Dsts {
			if c, ok := conj[d]; ok {
				c[m.Name] = false
			}
		}
	}
	ff := make(map[string]bool)

	var highN, lowN int
	push := func(q *container.FIFO[Pulse], from, to string, high bool) {
		if high {
			highN++
		} else {
			lowN++
		}
		q.Push(Pulse{from, to, high})
	}

	for i := 0; i < 1000; i++ {
		q := new(container.FIFO[Pulse])
		push(q, "button", "broadcaster", false)
		for q.Len() > 0 {
			p := q.Pop()
			m := mods[p.To]
			switch m.Kind {
			case Broadcaster:
				for _, d := range m.Dsts {
					push(q, m.Name, d, p.High)
				}
			case FlipFlop:
				if !p.High {
					high := !ff[m.Name]
					ff[m.Name] = high
					for _, d := range m.Dsts {
						push(q, m.Name, d, high)
					}
				}
			case Conjunction:
				conj[m.Name][p.From] = p.High
				high := true
				for _, v := range conj[m.Name] {
					high = high && v
				}
				high = !high
				for _, d := range m.Dsts {
					push(q, m.Name, d, high)
				}
			default:
				panic("invalid kind")
			}
		}
	}
	return highN * lowN
}

func Part2(in []Module) int {
	mods := make(map[string]Module)
	for _, m := range in {
		mods[m.Name] = m
	}

	// The general structure of the network is (see input.dot):
	// - The broadcaster is connected to four chains of Flip-flops. Each
	//   doubles the frequency with which the next one sends a low pulse,
	//   building a shift register.
	// - Some Flip-flops are connected to an output Conjunction, which thus
	//   sends a low pulse if those exact Flip-flops send a high pulse. This
	//   can be read as a binary number, determining the count at which the
	//   output sends a low pulse.
	// - The output is then connected to the rest of the Flip-flops, which acts
	//   to "reset" the shift register, once the count is read.
	// - Lastly, each shift register output is connected via two not-gates (a
	//   Conjunction with one input and one output) to a Conjunction which then
	//   feeds into rx.
	// Thus rx turns off at the least common multiple of the frequencies of the
	// shift registers, which we can read off.

	bc := mods["broadcaster"]
	ch := make(chan int, len(bc.Dsts))
	for _, d := range bc.Dsts {
		go func(m Module) {
			var (
				f   int
				bit = 1
			)
			for {
				done := true
				for _, d := range m.Dsts {
					switch mods[d].Kind {
					case Conjunction:
						f |= bit
					case FlipFlop:
						m = mods[d]
						done = false
					}
				}
				bit <<= 1
				if done {
					ch <- f
					return
				}
			}
		}(mods[d])
	}
	var m int = 1
	for range len(bc.Dsts) {
		m = math.LCM(m, <-ch)
	}
	return m
}

type Module struct {
	Kind Kind
	Name string
	Dsts []string
}

type Kind uint8

const (
	Broadcaster Kind = iota
	FlipFlop
	Conjunction
)

type Pulse struct {
	From string
	To   string
	High bool
}

func DumpGraph(in []Module) {
	fmt.Println("digraph G {")
	for _, m := range in {
		shape := "rect"
		switch m.Kind {
		case Broadcaster:
			shape = "diamond"
		case FlipFlop:
			shape = "invtriangle"
		case Conjunction:
			shape = "rect"
		}
		fmt.Printf("\t%s [shape=%q]\n", m.Name, shape)
		for _, d := range m.Dsts {
			fmt.Printf("\t%s -> %s\n", m.Name, d)
		}
	}
	fmt.Println("}")
}
