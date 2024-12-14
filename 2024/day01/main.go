package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	lines, err := parse.Lines(parse.Array[[2]int](split.Fields, parse.Signed[int]))(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	_ = lines
	a1, a2 := make([]int, len(lines)), make([]int, len(lines))
	for i, l := range lines {
		a1[i], a2[i] = l[0], l[1]
	}
	slices.Sort(a1)
	slices.Sort(a2)
	var total int
	for i := range lines {
		if a1[i] > a2[i] {
			total += a1[i] - a2[i]
		} else {
			total += a2[i] - a1[i]
		}
	}
	fmt.Println(total)
	m := make(map[int]int)
	for _, v := range a2 {
		m[v]++
	}
	total = 0
	for _, v := range a1 {
		total += v * m[v]
	}
	fmt.Println(total)
}
