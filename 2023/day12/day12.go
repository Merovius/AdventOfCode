package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
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
	var (
		maxNG, maxG, maxS = math.MinInt, math.MinInt, math.MinInt
	)
	for _, r := range in {
		maxNG = max(maxNG, len(r.Groups))
		for _, g := range r.Groups {
			maxG = max(maxG, g)
		}
		maxS = max(maxS, len(r.Springs))
	}
	fmt.Println(maxNG, maxG, maxS)
}

func Parse(s string) ([]Record, error) {
	return parse.TrimSpace(parse.Lines(
		parse.Struct[Record](
			split.Fields,
			parse.String[string],
			parse.Slice(
				split.On(","),
				parse.Signed[int],
			),
		),
	))(s)
}

func Part1(in []Record) int {
	var total int
	for _, r := range in {
		total += Arrangements(r.Springs, r.Groups)
	}
	return total
}

func Part2(in []Record) int {
	var total int
	for _, r := range in {
		total += Arrangements2(r.Springs, r.Groups)
	}
	return total
}

type Record struct {
	Springs string
	Groups  []int
}

func Arrangements2(springs string, groups []int) int {
	springs = strings.Join(repeat([]string{springs}, 5), "?")
	groups = repeat(groups, 5)
	return Arrangements(springs, groups)
}

func Arrangements(springs string, groups []int) (total int) {
	return divide(strings.Trim(springs, "."), groups)
}

var memo = make(map[any]int)

func memoKey(springs string, groups []int) any {
	var key struct {
		springs string
		groups  [32]uint8
	}
	key.springs = springs
	if len(groups) > 32 {
		panic("can not use more than 32 groups")
	}
	for i, g := range groups {
		key.groups[i] = uint8(g)
	}
	return key
}

// divide divides springs into rough halves on a run of dots. Dots must
// separate groups, so we can subdivide the problem by which groups to put into
// which half.
//
// springs must have no leading or trailing dots.
func divide(springs string, groups []int) (total int) {
	key := memoKey(springs, groups)
	if v, ok := memo[key]; ok {
		return v
	}
	defer func() {
		memo[key] = total
	}()

	if springs == "" {
		if len(groups) == 0 {
			return 1
		}
		return 0
	}

	l, r := splitOnDots(springs)
	if len(l) > 0 && len(r) > 0 {
		for i := 0; i <= len(groups); i++ {
			if v := divide(l, groups[:i]); v > 0 {
				total += v * divide(r, groups[i:])
			}
		}
		return total
	}
	// note: l and r can not both be empty, as springs had no trailing/leading
	// dots but is not empty. So there must be at least one non-dot byte.
	if len(l) == 0 {
		return conquer(r, groups)
	}
	return conquer(l, groups)
}

// conquer calculates all possible ways to put groups into springs.
// springs must not contain any dots.
func conquer(springs string, groups []int) (total int) {
	// memoize result
	key := memoKey(springs, groups)
	if v, ok := memo[key]; ok {
		return v
	}
	defer func() {
		memo[key] = total
	}()

	if len(groups) == 0 {
		if strings.IndexByte(springs, '#') >= 0 {
			return 0
		}
		return 1
	}
	g := groups[0]
	if len(springs) < g {
		return 0
	}
	if springs[0] == '?' {
		// try skipping this
		total = conquer(springs[1:], groups)
	}
	// now springs must start with an (actual or attempted) group
	springs = springs[g:]
	if len(groups) == 1 {
		return total + conquer(springs, groups[1:])
	}
	if len(springs) == 0 || springs[0] == '#' {
		return total
	}
	return total + conquer(springs[1:], groups[1:])
}

func repeat[T any](s []T, n int) []T {
	out := make([]T, 0, len(s)*5)
	for range n {
		out = append(out, s...)
	}
	return out
}

// splitOnDots divides s on a run of . roughly in the middle. If there are no
// dots, splitOnDots returns s, "".
//
// The input must have no leading or trailing dots. The results have leading
// and trailing dots removed.
func splitOnDots(s string) (l, r string) {
	m := len(s) / 2
	// note: No danger of an off-by-one error, as neither the first nor the
	// last byte of s is a dot.
	for i, j := m, m; i >= 0 && j < len(s); i, j = i-1, j+1 {
		if s[i] == '.' {
			return strings.TrimRight(s[:i], "."), strings.TrimLeft(s[i:], ".")
		}
		if s[j] == '.' {
			return strings.TrimRight(s[:j], "."), strings.TrimLeft(s[j:], ".")
		}
	}
	return s, ""
}
