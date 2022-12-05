package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"gonih.org/stack"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	stacks, err := ReadStacks(s)
	if err != nil {
		log.Fatal(err)
	}
	insts, err := ReadInstructions(s)
	if err != nil {
		log.Fatal(err)
	}
	Exec(stacks, insts)
	fmt.Printf("Tops of the stacks: %q\n", Tops(stacks))
}

func ReadStacks(s *bufio.Scanner) ([]stack.Stack[byte], error) {
	var lines [][]byte
scanloop:
	for s.Scan() {
		s := s.Text()
		if s == "" {
			break
		}
		var line []byte
		for i := 0; i < len(s); i += 4 {
			cell := s[i : i+3]
			if cell == "   " {
				line = append(line, 0)
			} else if cell == " 1 " {
				continue scanloop
			} else if cell[0] != '[' || cell[2] != ']' {
				return nil, fmt.Errorf("invalid stack line %q", s)
			} else {
				line = append(line, cell[1])
			}
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return nil, errors.New("no stacks given")
	}
	n := len(lines[0])
	for i, l := range lines[1:] {
		if len(l) != n {
			return nil, fmt.Errorf("number of stacks is not consistent along lines (line 1 -> %d, line %d -> %d)", n, i+1, len(l))
		}
	}
	stacks := make([]stack.Stack[byte], n)
	for i := len(lines) - 1; i >= 0; i-- {
		for j, b := range lines[i] {
			if b == 0 {
				continue
			}
			stacks[j].Push(b)
		}
	}
	return stacks, nil
}

func ReadInstructions(s *bufio.Scanner) ([]Inst, error) {
	var out []Inst
	for s.Scan() {
		var i Inst
		j, err := fmt.Sscanf(s.Text(), "move %d from %d to %d\n", &i.N, &i.From, &i.To)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if j != 3 || err != nil {
			return nil, err
		}
		i.From -= 1
		i.To -= 1
		out = append(out, i)
	}
	return out, s.Err()
}

type Stacks []stack.Stack[byte]

type Inst struct {
	N    int
	From int
	To   int
}

func Exec(stacks []stack.Stack[byte], prog []Inst) {
	var buf stack.Stack[byte]
	for _, inst := range prog {
		for j := 0; j < inst.N; j++ {
			buf.Push(stacks[inst.From].Pop())
		}
		for len(buf) > 0 {
			stacks[inst.To].Push(buf.Pop())
		}
	}
}

func Tops(stacks []stack.Stack[byte]) string {
	var out []byte
	for _, s := range stacks {
		if len(s) > 0 {
			out = append(out, s.Top())
		}
	}
	return string(out)
}
