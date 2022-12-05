package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
)

//go:embed example.txt
var exampleFile string

//go:embed input.txt
var inputFile string

var exampleStacks = [][]byte{
	{'Z', 'N'},
	{'M', 'C', 'D'},
	{'P'},
}

var inputStacks = [][]byte{
	{'B', 'P', 'N', 'Q', 'H', 'D', 'R', 'T'},
	{'W', 'G', 'B', 'J', 'T', 'V'},
	{'N', 'R', 'H', 'D', 'S', 'V', 'M', 'Q'},
	{'P', 'Z', 'N', 'M', 'C'},
	{'D', 'Z', 'B'},
	{'V', 'C', 'W', 'Z'},
	{'G', 'Z', 'N', 'C', 'V', 'Q', 'L', 'S'},
	{'L', 'G', 'J', 'M', 'D', 'N', 'V'},
	{'T', 'P', 'M', 'F', 'Z', 'C', 'G'},
}

func main() {
	input := flag.Bool("input", false, "run on input code")
	flag.Parse()
	var (
		stacks [][]byte
		insts  []Inst
		err    error
	)
	if *input {
		stacks = inputStacks
		insts, err = ReadInstructions(strings.NewReader(inputFile))
	} else {
		stacks = exampleStacks
		insts, err = ReadInstructions(strings.NewReader(exampleFile))
	}
	if err != nil {
		log.Fatal(err)
	}
	Exec(stacks, insts)
	fmt.Printf("Tops of the stacks: %q\n", Tops(stacks))
}

func ReadInstructions(r io.Reader) ([]Inst, error) {
	var out []Inst
	for n := 0; ; n++ {
		var i Inst
		j, err := fmt.Fscanf(r, "move %d from %d to %d\n", &i.N, &i.From, &i.To)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if j != 3 || err != nil {
			return nil, fmt.Errorf("line %d: %w", n+1, err)
		}
		i.From -= 1
		i.To -= 1
		out = append(out, i)
	}
	return out, nil
}

type Stacks [][]byte

type Inst struct {
	N    int
	From int
	To   int
}

func Exec(stacks [][]byte, prog []Inst) {
	pop := func(i int) (v byte) {
		n := len(stacks[i]) - 1
		stacks[i], v = stacks[i][:n], stacks[i][n]
		return v
	}
	push := func(i int, b byte) {
		stacks[i] = append(stacks[i], b)
	}
	for _, inst := range prog {
		var buf []byte
		for j := 0; j < inst.N; j++ {
			buf = append(buf, pop(inst.From))
		}
		for j := len(buf) - 1; j >= 0; j-- {
			push(inst.To, buf[j])
		}
	}
}

func Tops(stacks [][]byte) string {
	var out []byte
	for _, s := range stacks {
		if len(s) == 0 {
			continue
		}
		out = append(out, s[len(s)-1])
	}
	return string(out)
}
