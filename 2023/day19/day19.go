package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/interval"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(s string) (Input, error) {
	return parse.TrimSpace(parse.Struct[Input](
		split.Blocks,
		parse.Map(
			split.Lines,
			split.Regexp(`([^{]+)\{(.*)\}`),
			parse.String[string],
			parse.Slice(
				split.On(","),
				func(s string) (Rule, error) {
					cond, dst, ok := strings.Cut(s, ":")
					if !ok {
						return Rule{Dst: s}, nil
					}
					var op byte
					prop, val, ok := strings.Cut(cond, ">")
					if ok {
						op = '>'
					} else {
						prop, val, ok = strings.Cut(cond, "<")
						if !ok {
							return Rule{}, fmt.Errorf("invalid condition %q", cond)
						}
						op = '<'
					}
					var (
						r   Rule
						err error
					)
					r.Prop, err = ParseProp(prop)
					if err != nil {
						return Rule{}, err
					}
					r.Op = op
					r.Val, err = strconv.Atoi(val)
					if err != nil {
						return Rule{}, err
					}
					r.Dst = dst
					return r, nil
				},
			),
		),
		parse.Lines(
			parse.Array[Part](
				split.Regexp(`\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)\}`),
				parse.Signed[int],
			),
		),
	))(s)
}

func Part1(in Input) int {
	get := func(s string) []Rule {
		if s == "A" {
			return []Rule{{Dst: "A"}}
		}
		if s == "R" {
			return []Rule{{Dst: "R"}}
		}
		return in.Workflows[s]
	}

	var total int
parts:
	for _, p := range in.Parts {
		w := get("in")
		for {
		workflow:
			for _, r := range w {
				switch r.Op {
				case 0:
					switch r.Dst {
					case "A":
						total += p[0] + p[1] + p[2] + p[3]
						fallthrough
					case "R":
						continue parts
					}
					w = get(r.Dst)
					break workflow
				case '>':
					if p[r.Prop] > r.Val {
						w = get(r.Dst)
						break workflow
					}
				case '<':
					if p[r.Prop] < r.Val {
						w = get(r.Dst)
						break workflow
					}
				}
			}
		}
	}
	return total
}

func Part2(in Input) int {
	type Intervals = interval.Set[interval.OO[int], int]
	const (
		min = 0
		max = 4001
	)
	var all [4]Intervals
	for i := 0; i < 4; i++ {
		all[i].Add(interval.MakeOO(min, max))
	}

	var rec func(w string, is [4]Intervals) int
	rec = func(w string, is [4]Intervals) int {
		if w == "R" {
			return 0
		}
		if w == "A" {
			v := is[0].Len() * is[1].Len() * is[2].Len() * is[3].Len()
			return v
		}

		var total int
		for _, r := range in.Workflows[w] {
			switch r.Op {
			case 0:
				return total + rec(r.Dst, is)
			case '>':
				js := is
				js[r.Prop] = *is[r.Prop].Clone()
				js[r.Prop].Intersect(interval.MakeOO(r.Val, max))
				total += rec(r.Dst, js)
				is[r.Prop] = *is[r.Prop].Clone()
				is[r.Prop].Intersect(interval.MakeOO(min, r.Val+1))
			case '<':
				js := is
				js[r.Prop] = *is[r.Prop].Clone()
				js[r.Prop].Intersect(interval.MakeOO(min, r.Val))
				total += rec(r.Dst, js)
				is[r.Prop] = *is[r.Prop].Clone()
				is[r.Prop].Intersect(interval.MakeOO(r.Val-1, max))
			}
		}
		panic("no terminal rule")
	}
	return rec("in", all)
}

type Input struct {
	Workflows map[string][]Rule
	Parts     []Part
}

type Workflow struct {
	Name  string
	Rules []Rule
}

type Rule struct {
	Prop Prop
	Op   byte
	Val  int
	Dst  string
}

func (r Rule) String() string {
	switch r.Op {
	case 0:
		return r.Dst
	default:
		return fmt.Sprintf("%v%c%d:%s", r.Prop, r.Op, r.Val, r.Dst)
	}
}

type Part [4]int

type Prop uint8

const (
	X Prop = iota
	M
	A
	S
)

func ParseProp(s string) (Prop, error) {
	switch s {
	case "x":
		return X, nil
	case "m":
		return M, nil
	case "a":
		return A, nil
	case "s":
		return S, nil
	default:
		return 0, fmt.Errorf("invalid prop %q", s)
	}
}

func (p Prop) String() string {
	switch p {
	case X:
		return "x"
	case M:
		return "m"
	case A:
		return "a"
	case S:
		return "s"
	default:
		panic("invalid prop")
	}
}
