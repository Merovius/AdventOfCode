package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/set"
)

func main() {
	contents, err := parse.Lines(parse.String[string]).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	byCompartment := SplitCompartments(contents)
	common, err := FindCommonItems(byCompartment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total priority of incorrectly packed items: %d\n", Priorities(common))
	badges, err := BadgeItems(contents)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total priority of badge items: %d\n", Priorities(badges))
}

func SplitCompartments(rs []string) [][2]string {
	var out [][2]string
	for _, r := range rs {
		out = append(out, [2]string{r[:len(r)/2], r[len(r)/2:]})
	}
	return out
}

func FindCommonItems(rs [][2]string) ([]byte, error) {
	var out []byte
	for _, r := range rs {
		i := strings.IndexAny(r[0], r[1])
		if i < 0 {
			return nil, fmt.Errorf("no common items between %q and %q", r[0], r[1])
		}
		out = append(out, r[0][i])
	}
	return out, nil
}

func Priorities[T ~byte | ~rune](common []T) int {
	var total int
	for _, c := range common {
		total += Priority(c)
	}
	return total
}

func BadgeItems(rs []string) ([]rune, error) {
	if len(rs)%3 != 0 {
		return nil, errors.New("number of inputs not divisible by 3")
	}
	var out []rune
	for i := 0; i < len(rs); i += 3 {
		var sets []set.Set[rune]
		for j := 0; j < 3; j++ {
			sets = append(sets, ToSet(rs[i+j]))
		}
		common := set.Intersect(sets...)
		if len(common) > 1 {
			return nil, fmt.Errorf("more than one item in common between %q, %q and %q: %v", rs[i], rs[i+1], rs[i+2], stringify(common))
		}
		out = append(out, common.Slice()[0])
	}
	return out, nil
}

func ToSet(s string) set.Set[rune] {
	out := make(set.Set[rune])
	for _, c := range s {
		out.Add(c)
	}
	return out
}

func Priority[T ~byte | ~rune](c T) int {
	switch {
	case c >= 'a' && c <= 'z':
		return int(c) - 'a' + 1
	default:
		return int(c) - 'A' + 27
	}
}

func stringify(s set.Set[rune]) string {
	var parts []string
	for r := range s {
		parts = append(parts, strconv.QuoteRune(r))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
