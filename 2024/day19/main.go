package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
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
	m := new(memo)

	var n int
	for _, d := range in.Designs {
		m.Clear()
		if count(m, true, in.Patterns, d) > 0 {
			n++
		}
	}
	return n
}

func Part2(in Input) int {
	m := new(memo)

	var n int
	for _, d := range in.Designs {
		m.Clear()
		n += count(m, false, in.Patterns, d)
	}
	return n
}

func count(m *memo, checkOnly bool, patterns *container.RadixSet, design string) (v int) {
	if v, ok := m.Get(design); ok {
		return v
	}
	defer func() { m.Set(design, v) }()

	if len(design) == 0 {
		return 1
	}
	var n int
	for p := range patterns.PrefixesOf(design) {
		n += count(m, checkOnly, patterns, design[len(p):])
		if checkOnly && n > 0 {
			return 1
		}
	}
	return n
}

// memo memoizes the result of count. If we memoize per-pattern, we can just
// use the length of the considered suffix as a key. This avoids having to
// hash and allocate. It duplicates work, but is still a big overall win.
//
// To simplify clearing, we store the return value offset by one, so 0
// signifies "no value".
type memo []int

func (m *memo) Get(s string) (int, bool) {
	if len(*m) > len(s) && (*m)[len(s)] > 0 {
		return (*m)[len(s)] - 1, true
	}
	return 0, false
}

func (m *memo) Set(s string, v int) {
	if len(s) >= len(*m) {
		*m = slices.Grow(*m, len(s)+1)[:len(s)+1]
	}
	(*m)[len(s)] = v + 1
}

func (m *memo) Clear() {
	clear(*m)
}
