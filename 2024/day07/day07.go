package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := parse.Lines(parse.Struct[Equation](split.On(": "), parse.Signed[int], parse.Slice(split.Fields, parse.Signed[int])))(string(buf))
	if err != nil {
		log.Fatal(err)
	}

	var total int
	for _, e := range in {
		if solvable(e.Test, e.Inputs) {
			total += e.Test
		}
	}
	fmt.Println(total)

	total = 0
	for _, e := range in {
		if solvable2(e.Test, e.Inputs) {
			total += e.Test
		}
	}
	fmt.Println(total)
}

type Equation struct {
	Test   int
	Inputs []int
}

func solvable(res int, ins []int) bool {
	n, ins := ins[len(ins)-1], ins[:len(ins)-1]
	if len(ins) == 0 {
		return res == n
	}
	if solvable(res-n, ins) {
		return true
	}
	return res%n == 0 && solvable(res/n, ins)
}

func solvable2(res int, ins []int) bool {
	n, ins := ins[len(ins)-1], ins[:len(ins)-1]
	if len(ins) == 0 {
		return res == n
	}
	if res >= n {
		if solvable2(res-n, ins) {
			return true
		}
	}
	if res%n == 0 {
		if solvable2(res/n, ins) {
			return true
		}
	}
	d := log10(n)
	p := pow10[d]
	return res%p == n && solvable2(res/p, ins)
}

func log10(n int) int {
	var d int = 1
	for n >= 10 {
		d++
		n /= 10
	}
	return d
}

var pow10 = [...]int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000}
