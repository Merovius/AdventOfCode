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
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(in string) ([][3]int, error) {
	return parse.Slice(
		split.Lines,
		parse.Array[[3]int](
			split.On("x"),
			parse.Signed[int],
		),
	)(in)
}

func Part1(in [][3]int) int {
	var total int
	for _, x := range in {
		a := x[0] * x[1]
		b := x[0] * x[2]
		c := x[1] * x[2]
		total += 2*(a+b+c) + min(a, b, c)
	}
	return total
}

func Part2(in [][3]int) int {
	var total int
	for _, x := range in {
		a := x[0] + x[1]
		b := x[0] + x[2]
		c := x[1] + x[2]
		total += 2*min(a, b, c) + x[0]*x[1]*x[2]
	}
	return total
}
