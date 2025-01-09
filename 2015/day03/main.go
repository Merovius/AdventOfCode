package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in := string(buf)
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Part1(in string) int {
	seen := make(set.Set[[2]int])
	var x [2]int
	seen.Add(x)
	for _, b := range in {
		switch b {
		case '<':
			x[0]--
		case '^':
			x[1]--
		case '>':
			x[0]++
		case 'v':
			x[1]++
		}
		seen.Add(x)
	}
	return len(seen)
}

func Part2(in string) int {
	seen := make(set.Set[[2]int])
	var (
		x [2]int
		y [2]int
	)
	seen.Add(x)
	for _, b := range in {
		switch b {
		case '<':
			x[0]--
		case '^':
			x[1]--
		case '>':
			x[0]++
		case 'v':
			x[1]++
		}
		seen.Add(x)
		x, y = y, x
	}
	return len(seen)
}
