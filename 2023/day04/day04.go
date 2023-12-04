//go:build goexperiment.rangefunc

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/iter"
	"github.com/Merovius/AdventOfCode/internal/op"
	"github.com/Merovius/AdventOfCode/internal/set"
	"github.com/Merovius/AdventOfCode/internal/slices"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	cards, err := parse.Lines(
		parse.Struct[Card](
			split.Regexp(`(.*): (.*) \| (.*)`),
			func(s string) (int, error) {
				idx, ok := strings.CutPrefix(s, "Card")
				if !ok {
					return 0, errors.New(`expected "Card"`)
				}
				idx = strings.TrimLeftFunc(idx, unicode.IsSpace)
				return strconv.Atoi(idx)
			},
			parse.Slice(split.Fields, parse.Signed[int]),
			parse.Slice(split.Fields, parse.Signed[int]),
		),
	)(string(bytes.TrimSpace(buf)))
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

func Part1(cards []Card) int {
	var total int
	for _, c := range cards {
		if n := iter.Len(WinningNumbers(c)); n > 0 {
			total += 1 << (n - 1)
		}
	}
	return total
}

func Part2(cards []Card) int {
	wins := make([]int, len(cards))
	for i, c := range slices.Backwards(cards) {
		n := iter.Len(WinningNumbers(c))
		wins[i] = iter.FoldR(
			op.Add,
			iter.Right(slices.Elements(wins[i+1:i+1+n])),
			1,
		)
	}
	return iter.FoldR(op.Add, iter.Right(slices.Elements(wins)), 0)
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
