package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

func Parse(in []byte) (Input, error) {
	return parse.Struct[Input](
		split.Blocks,
		parse.Slice(
			split.On(", "),
			parse.String[string],
		),
		parse.Slice(
			split.Lines,
			parse.String[string],
		),
	)(string(in))
}

type Input struct {
	Patterns []string
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

func count(memo map[string]int, checkOnly bool, patterns []string, design string) (v int) {
	if v, ok := memo[design]; ok {
		return v
	}
	defer func() { memo[design] = v }()

	if len(design) == 0 {
		return 1
	}
	var n int
	for _, p := range patterns {
		rest, ok := strings.CutPrefix(design, p)
		if !ok {
			continue
		}
		n += count(memo, checkOnly, patterns, rest)
		if checkOnly && n > 0 {
			return 1
		}
	}
	return n
}
