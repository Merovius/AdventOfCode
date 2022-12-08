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
	vm, err := parseProgram(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Accumulator value on first loop: %v\n", findLoopValue(vm))

	fixed := fixProgram(vm)
	fmt.Printf("Need to flip instruction %d to fix the program\n", fixed)
	fmt.Println("Final accumulator value:", vm.Run())
}

func fixProgram(vm *vm) int {
	defer vm.Reset()

	// build a graph, where the nodes are instructions and there is an edge
	// from A to B, if A is the next instruction to execute after B. We mark
	// each node with whether or not it halts the program in the second step.
	type node struct {
		halts     bool
		neighbors []int
	}
	graph := make([]node, len(vm.prog)+1)
	for i, in := range vm.prog {
		switch in.op {
		case "acc", "nop":
			graph[i+1].neighbors = append(graph[i+1].neighbors, i)
		case "jmp":
			graph[i+in.arg].neighbors = append(graph[i+in.arg].neighbors, i)
		}
	}
	// Execute a DFS, marking all nodes that are reachable from the last one as
	// halting.
	var visit func(int)
	visit = func(i int) {
		if graph[i].halts {
			return
		}
		graph[i].halts = true
		for _, n := range graph[i].neighbors {
			visit(n)
		}
	}
	visit(len(graph) - 1)
	halts := func(i int) bool {
		return i < len(graph) && graph[i].halts
	}

	// Step through the program. At each step, check if a flip would cause us
	// to reach a halting instruction next. If so, that is the instruction we
	// need to fix.
	for !vm.Halted() {
		switch in := &vm.prog[vm.pc]; in.op {
		case "acc":
		case "nop":
			if halts(vm.pc + in.arg) {
				in.op = "jmp"
				return vm.pc
			}
		case "jmp":
			if halts(vm.pc + 1) {
				in.op = "nop"
				return vm.pc
			}
		}
		vm.Step()
	}
	return len(vm.prog)
}

func eq(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func findLoopValue(vm *vm) int {
	visited := make([]bool, len(vm.prog))
	for {
		if visited[vm.pc] {
			return vm.acc
		}
		visited[vm.pc] = true
		vm.Step()
	}
}

type vm struct {
	prog []inst
	pc   int
	acc  int
}

func (vm *vm) Run() int {
	for !vm.Halted() {
		vm.Step()
	}
	return vm.acc
}

func (vm *vm) Step() {
	if vm.Halted() {
		return
	}
	switch i := vm.prog[vm.pc]; i.op {
	case "nop":
		vm.pc++
	case "acc":
		vm.acc += i.arg
		vm.pc++
	case "jmp":
		vm.pc += i.arg
	}
}

func (vm *vm) Halted() bool {
	return vm.pc >= len(vm.prog)
}

func (vm *vm) Reset() {
	vm.pc, vm.acc = 0, 0
}

type inst struct {
	op  string
	arg int
}

func parseProgram(r io.Reader) (*vm, error) {
	var prog []inst
	s := bufio.NewScanner(r)
	for s.Scan() {
		sp := strings.Split(s.Text(), " ")
		if len(sp) != 2 {
			return nil, fmt.Errorf("invalid line %q", s.Text())
		}
		if !validOpCodes[sp[0]] {
			return nil, fmt.Errorf("invalid line %q: unknown opcode", s.Text())
		}
		v, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, fmt.Errorf("invalid line %q: %w", s.Text(), err)
		}
		prog = append(prog, inst{sp[0], v})
	}
	return &vm{prog: prog}, nil
}

var validOpCodes = map[string]bool{
	"jmp": true,
	"nop": true,
	"acc": true,
}
