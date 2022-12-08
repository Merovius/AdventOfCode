package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	prog, err := ParseProgram(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	x, y := ExecuteNaively(prog)
	fmt.Printf("Naive position is (%d,%d), distance is %d\n", x, y, abs(x)+abs(y))
	x, y = ExecuteCorrectly(prog)
	fmt.Printf("Correct position is (%d,%d), distance is %d\n", x, y, abs(x)+abs(y))
}

func ParseProgram(r io.Reader) ([]Inst, error) {
	var prog []Inst
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if len(l) < 2 {
			return nil, fmt.Errorf("can't parse instruction %q", l)
		}
		arg, err := strconv.Atoi(l[1:])
		if err != nil {
			return nil, fmt.Errorf("can't parse instruction %q: %w", l, err)
		}
		switch l[0] {
		case 'N', 'S', 'E', 'W', 'F':
		case 'L', 'R':
			if arg%90 != 0 {
				return nil, fmt.Errorf("can't parse instruction %q: argument must be multiple of 90", l)
			}
			arg /= 90
		default:
			return nil, fmt.Errorf("can't parse instruction %q: unknown op code", l)
		}
		prog = append(prog, Inst{l[0], arg})
	}
	return prog, nil
}

type Inst struct {
	Op  byte
	Arg int
}

func ExecuteNaively(p []Inst) (x, y int) {
	var (
		dx, dy = 1, 0
	)
	for _, i := range p {
		switch i.Op {
		case 'N':
			y += i.Arg
		case 'E':
			x += i.Arg
		case 'S':
			y -= i.Arg
		case 'W':
			y += i.Arg
		case 'L':
			for j := 0; j < i.Arg; j++ {
				dx, dy = -dy, dx
			}
		case 'R':
			for j := 0; j < i.Arg; j++ {
				dx, dy = dy, -dx
			}
		case 'F':
			x += dx * i.Arg
			y += dy * i.Arg
		}
	}
	return x, y
}

func ExecuteCorrectly(p []Inst) (x, y int) {
	var (
		dx, dy = 10, 1
	)
	for _, i := range p {
		switch i.Op {
		case 'N':
			dy += i.Arg
		case 'E':
			dx += i.Arg
		case 'S':
			dy -= i.Arg
		case 'W':
			dx -= i.Arg
		case 'L':
			for j := 0; j < i.Arg; j++ {
				dx, dy = -dy, dx
			}
		case 'R':
			for j := 0; j < i.Arg; j++ {
				dx, dy = dy, -dx
			}
		case 'F':
			x += dx * i.Arg
			y += dy * i.Arg
		}
	}
	return x, y
}

func abs(n int) int {
	if n > 0 {
		return n
	}
	return -n
}
