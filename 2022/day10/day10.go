package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	prog, err := parse.Lines(
		parse.Any[Inst](
			func(s string) (Inst, error) {
				if s == "noop" {
					return Inst{Op: "noop"}, nil
				}
				return Inst{}, errors.New(`expected "noop"`)
			},
			parse.Struct[Inst](split.Fields, parse.Enum("addx"), parse.Signed[int]),
		),
	).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	sig, disp := Run(prog)
	fmt.Println("Signal strength total:", sig)
	fmt.Println("CRT display:")
	for _, l := range disp {
		fmt.Println(l)
	}
}

type Inst struct {
	Op  string
	Arg int
}

func Run(prog []Inst) (int, []string) {
	var (
		x     = 1
		t     = 1
		pc    = 0
		stall bool
		sig   int
		crt   = make([]rune, 40*6)
		r, c  = 0, 0
	)
	for i := range crt {
		crt[i] = ' '
	}
	for pc < len(prog) || stall {
		if t%40 == 20 {
			sig += x * t
		}
		r, c = ((t-1)/40)%6, (t-1)%40
		if math.Abs(c-x) <= 1 {
			crt[r*40+c] = '▇'
		} else {
			crt[r*40+c] = ' '
		}
		switch {
		case prog[pc].Op == "noop":
			pc++
		case stall:
			x += prog[pc].Arg
			stall = false
			pc++
		default:
			stall = true
		}
		t++
	}
	var disp []string
	for i := 0; i < len(crt); i += 40 {
		disp = append(disp, string(crt[i:i+40]))
	}
	return sig, disp
}
