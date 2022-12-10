package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	prog, err := ReadInput(os.Stdin)
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

func ReadInput(r io.Reader) ([]Inst, error) {
	var out []Inst
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		op, argS, ok := strings.Cut(l, " ")
		var arg int
		if ok {
			var err error
			arg, err = strconv.Atoi(argS)
			if err != nil {
				return nil, err
			}
		}
		if op != "noop" && op != "addx" {
			return nil, fmt.Errorf("invalid instruction %q", op)
		}
		out = append(out, Inst{op, arg})
	}
	return out, nil
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
