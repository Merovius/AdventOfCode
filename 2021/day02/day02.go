package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	prog, err := Parse(f)
	if err != nil {
		return err
	}

	fmt.Printf("Final position, part 1: %+v\n", Run1(prog))
	fmt.Printf("Final position, part 2: %+v\n", Run2(prog))
	return nil
}

func Run1(prog []Command) Pos {
	var p Pos
	for _, c := range prog {
		switch c.Op {
		case OpForward:
			p.Horizontal += c.Count
		case OpDown:
			p.Depth += c.Count
		case OpUp:
			p.Depth -= c.Count
		}
	}
	return p
}

func Run2(prog []Command) Pos {
	var p Pos
	for _, c := range prog {
		switch c.Op {
		case OpForward:
			p.Horizontal += c.Count
			p.Depth += c.Count * p.Aim
		case OpDown:
			p.Aim += c.Count
		case OpUp:
			p.Aim -= c.Count
		}
	}
	return p
}

type Pos struct {
	Horizontal int
	Depth      int
	Aim        int
}

type Command struct {
	Op    Op
	Count int
}

type Op int

const (
	OpForward Op = iota
	OpUp
	OpDown
)

func Parse(r io.Reader) ([]Command, error) {
	var out []Command

	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		i := strings.Index(l, " ")
		if i < 0 {
			return nil, fmt.Errorf("invalid input line %q", l)
		}
		c, err := strconv.Atoi(l[i+1:])
		if err != nil {
			return nil, fmt.Errorf("invalid input line %q: %v", l, err)
		}

		switch op := l[:i]; op {
		case "forward":
			out = append(out, Command{OpForward, c})
		case "up":
			out = append(out, Command{OpUp, c})
		case "down":
			out = append(out, Command{OpDown, c})
		default:
			return nil, fmt.Errorf("invalid operation %q", op)
		}
	}

	return out, nil
}
