package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
)

func main() {
	log.SetFlags(log.Lshortfile)
	data, err := parse.Lines(parse.Signed[int]).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(data))
	fmt.Println(Part2(data))
	_ = data
}

func Part1(l []int) int {
	list, zero := ToList(l, 1)
	Mix(list)
	return CoordinateSum(&list[zero])
}

func Part2(l []int) int {
	list, zero := ToList(l, 811589153)
	for i := 0; i < 10; i++ {
		Mix(list)
	}
	return CoordinateSum(&list[zero])
}

func Mix(l []Node) {
	for i := range l {
		n := &l[i]
		k, swap := n.V, n.SwapRight
		if k < 0 {
			k, swap = -k, n.SwapLeft
		}
		k %= len(l) - 1
		for i := 0; i < k; i++ {
			swap()
		}
	}
}

func CoordinateSum(n *Node) int {
	var c int
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			n = n.R
		}
		c += n.V
	}
	return c
}

type Node struct {
	V int
	R *Node
	L *Node
}

func ToList(s []int, mul int) (n []Node, zero int) {
	n = make([]Node, len(s))
	for i, v := range s {
		n[i].V = v * mul
		if v == 0 {
			zero = i
		}
		switch {
		case i == 0:
			n[i].L = &n[len(n)-1]
			n[i].R = &n[i+1]
		case i == len(n)-1:
			n[i].L = &n[i-1]
			n[i].R = &n[0]
		default:
			n[i].L = &n[i-1]
			n[i].R = &n[i+1]
		}
	}
	return n, zero
}

func (n *Node) SwapLeft() {
	nll, nl, nr := n.L.L, n.L, n.R

	nll.R = n
	nl.L, nl.R = n, nr
	n.L, n.R = nll, nl
	nr.L = nl
}

func (n *Node) SwapRight() {
	nrr, nr, nl := n.R.R, n.R, n.L

	nrr.L = n
	nr.R, nr.L = n, nl
	n.R, n.L = nrr, nr
	nl.R = nr
}

func (n Node) String() string {
	parts := []string{fmt.Sprint(n.V)}
	for m := n.R; m.R != n.R; m = m.R {
		parts = append(parts, fmt.Sprint(m.V))
	}
	return "[" + strings.Join(parts, " ") + "]"
}
