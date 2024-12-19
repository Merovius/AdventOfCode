package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(buf []byte) (Input, error) {
	return parse.Struct[Input](
		split.Blocks,
		func(s string) (*container.RadixSet, error) {
			r := new(container.RadixSet)
			for p := range strings.SplitSeq(s, ", ") {
				r.Add(p)
			}
			return r, nil
		},
		parse.Slice(
			split.Lines,
			parse.String[string],
		),
	)(string(buf))
}

type Input struct {
	Patterns *container.RadixSet
	Designs  []string
}

func Part1(in Input) int {
	memo := make(map[string]int)

	var n int
	for _, d := range in.Designs {
		if count(memo, true, in.Patterns, d) > 0 {
			n++
		}
	}
	return n
}

func Part2(in Input) int {
	memo := make(map[string]int)

	var n int
	for _, d := range in.Designs {
		n += count(memo, false, in.Patterns, d)
	}
	return n
}

var memo map[string]int

func count(memo map[string]int, checkOnly bool, patterns *container.RadixSet, design string) (v int) {
	if v, ok := memo[design]; ok {
		return v
	}
	defer func() { memo[design] = v }()

	if len(design) == 0 {
		return 1
	}
	var n int
	for p := range patterns.PrefixesOf(design) {
		n += count(memo, checkOnly, patterns, design[len(p):])
		if checkOnly && n > 0 {
			return 1
		}
	}
	return n
}
