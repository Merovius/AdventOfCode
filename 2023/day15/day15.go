package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

func Parse(s string) ([]string, error) {
	return parse.TrimSpace(parse.Slice(split.On(","), parse.String[string]))(s)
}

func Part1(in []string) int {
	var out int
	for _, s := range in {
		out += int(HASH(s))
	}
	return out
}

func Part2(in []string) int {
	boxes := make([][]Lens, 256)
	for _, s := range in {
		l, vs, ok := strings.Cut(s, "=")
		if ok {
			v, err := strconv.Atoi(vs)
			if err != nil {
				panic(err)
			}
			i := HASH(l)
			j := slices.IndexFunc(boxes[i], func(lens Lens) bool {
				return lens.Label == l
			})
			if j < 0 {
				boxes[i] = append(boxes[i], Lens{l, v})
			} else {
				boxes[i][j].Length = v
			}
		} else {
			l, ok = strings.CutSuffix(s, "-")
			if !ok {
				panic(fmt.Sprintf("invalid input %q", s))
			}
			i := HASH(l)
			boxes[i] = slices.DeleteFunc(boxes[i], func(lens Lens) bool {
				return lens.Label == l
			})
		}
	}
	var total int
	for i, b := range boxes {
		for j, l := range b {
			total += (i + 1) * (j + 1) * l.Length
		}
	}
	return total
}

func HASH(s string) uint8 {
	var v uint8
	for i := 0; i < len(s); i++ {
		v += s[i]
		v *= 17
	}
	return v
}

type Lens struct {
	Label  string
	Length int
}
