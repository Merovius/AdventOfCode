package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
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

func Parse(in []byte) ([]Instruction, error) {
	return parse.Slice(
		split.Lines,
		parse.Struct[Instruction](
			split.Regexp(`(turn off|turn on|toggle) (\d+,\d+) through (\d+,\d+)`),
			func(s string) (Op, error) {
				switch s {
				case "turn off":
					return Off, nil
				case "turn on":
					return On, nil
				case "toggle":
					return Toggle, nil
				default:
					return 0, fmt.Errorf("invalid instruction %q", s)
				}
			},
			parse.Struct[grid.Pos](split.On(","), parse.Signed[int], parse.Signed[int]),
			parse.Struct[grid.Pos](split.On(","), parse.Signed[int], parse.Signed[int]),
		),
	)(string(in))
}

type Op int

const (
	_ Op = iota
	Off
	On
	Toggle
)

func (o Op) Apply(b bool) bool {
	switch o {
	case Off:
		return false
	case On:
		return true
	case Toggle:
		return !b
	default:
		panic("invalid Op")
	}
}

func (o Op) Apply2(b int) int {
	switch o {
	case Off:
		return max(b-1, 0)
	case On:
		return b + 1
	case Toggle:
		return b + 2
	default:
		panic("invalid Op")
	}
}

type Instruction struct {
	Op  Op
	Min grid.Pos
	Max grid.Pos
}

func Part1(in []Instruction) int {
	g := grid.New[bool](1000, 1000)
	for _, i := range in {
		r := grid.Rectangle{i.Min, i.Max.Add(grid.Pos{1, 1})}
		for p := range r.All() {
			g.Set(p, i.Op.Apply(g.At(p)))
		}
	}
	var n int
	for _, v := range g.All() {
		if v {
			n++
		}
	}
	return n
}

func Part2(in []Instruction) int {
	g := grid.New[int](1000, 1000)
	for _, i := range in {
		r := grid.Rectangle{i.Min, i.Max.Add(grid.Pos{1, 1})}
		for p := range r.All() {
			g.Set(p, i.Op.Apply2(g.At(p)))
		}
	}
	var n int
	for _, v := range g.All() {
		n += v
	}
	return n

}
