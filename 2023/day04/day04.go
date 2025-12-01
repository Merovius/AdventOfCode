package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"slices"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/op"
	"gonih.org/AdventOfCode/internal/set"
	"gonih.org/AdventOfCode/internal/xiter"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	cards, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Point total:", Part1(cards))
	fmt.Println("Total number of cards:", Part2(cards))
}

type Card struct {
	Index   int
	Winning []int
	Have    []int
}

func Parse(b []byte) ([]Card, error) {
	return parse.Lines(
		parse.Struct[Card](
			split.Regexp(`Card \s*(\d+): (.*) \| (.*)`),
			parse.Signed[int],
			parse.Slice(split.Fields, parse.Signed[int]),
			parse.Slice(split.Fields, parse.Signed[int]),
		),
	)(string(bytes.TrimSpace(b)))
}

func Part1(cards []Card) int {
	var total int
	for _, c := range cards {
		if n := xiter.Len(WinningNumbers(c)); n > 0 {
			total += 1 << (n - 1)
		}
	}
	return total
}

func Part2(cards []Card) int {
	wins := make([]int, len(cards))
	for i, c := range slices.Backward(cards) {
		n := xiter.Len(WinningNumbers(c))
		wins[i] = xiter.FoldR(
			op.Add,
			xiter.Right(slices.All(wins[i+1:i+1+n])),
			1,
		)
	}
	return xiter.FoldR(op.Add, xiter.Right(slices.All(wins)), 0)
}

func WinningNumbers(c Card) iter.Seq[int] {
	return func(yield func(int) bool) {
		h := set.Make(c.Have...)
		for _, n := range c.Winning {
			if h.Contains(n) {
				if !yield(n) {
					return
				}
			}
		}
	}
}
