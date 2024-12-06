package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"log"
	"math"
	"os"
	"slices"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/interval"
	"github.com/Merovius/AdventOfCode/internal/xiter"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

type Almanach struct {
	Seeds []int
	Maps  []Map
}

type Map struct {
	SrcKind string
	DstKind string
	Ranges  []Range
}

type Range struct {
	DstStart int
	SrcStart int
	Length   int
}

func (r Range) String() string {
	return fmt.Sprintf("{d=%d, s=%d, l=%d}", r.DstStart, r.SrcStart, r.Length)
}

func Parse(b []byte) (Almanach, error) {
	return parse.Struct[Almanach](
		split.SplitN("\n\n", 2),
		parse.Prefix("seeds: ", parse.Slice(split.Fields, parse.Signed[int])),
		parse.Blocks(
			parse.Struct[Map](
				split.Regexp(`(?s:(\w+)-to-(\w+) map:\n(.*))`),
				parse.String[string],
				parse.String[string],
				parse.Slice(
					split.Lines,
					parse.Struct[Range](
						split.Fields,
						parse.Signed[int],
						parse.Signed[int],
						parse.Signed[int],
					),
				),
			),
		),
	)(string(bytes.TrimSpace(b)))
}

func Part1(in Almanach) int {
	nums := slices.Clone(in.Seeds)
	for _, m := range in.Maps {
		for i, v := range nums {
			for _, r := range m.Ranges {
				if v >= r.SrcStart && v < r.SrcStart+r.Length {
					nums[i] = r.DstStart + (v - r.SrcStart)
					break
				}
			}
		}
	}
	m := nums[0]
	for _, v := range nums[1:] {
		m = min(m, v)
	}
	return m
}

func Part2(a Almanach) int {
	var intervals []Interval
	for i := 0; i < len(a.Seeds); i += 2 {
		intervals = append(intervals, Interval{a.Seeds[i], a.Seeds[i] + a.Seeds[i+1]})
	}
	in := slices.Values(intervals)
	for _, m := range a.Maps {
		in = m.Apply(in)
	}
	return xiter.FoldR(func(i Interval, m int) int {
		return min(i.Min, m)
	}, in, math.MaxInt)
}

type Interval = interval.CO[int]

func (r Range) Src() Interval {
	return Interval{r.SrcStart, r.SrcStart + r.Length}
}

func (r Range) Apply(i Interval) Interval {
	o := r.DstStart - r.SrcStart
	i.Min += o
	i.Max += o
	return i
}

// Split i into a part that is inside r and a part that is left, inside and
// right of r.
func (r Range) Split(i Interval) (left, in, right Interval) {
	j := r.Src()
	if i.Min < j.Min {
		left.Min = i.Min
		if i.Max <= j.Min {
			left.Max = i.Max
			return left, in, right
		}
		left.Max, i.Min = j.Min, j.Min
	}
	if i.Max > j.Max {
		right.Max = i.Max
		if i.Min >= j.Max {
			right.Min = i.Min
			return left, in, right
		}
		right.Min, i.Max = j.Max, j.Max
	}
	in = i
	return left, in, right
}

func (m Map) Apply(in iter.Seq[Interval]) iter.Seq[Interval] {
	return func(yield func(Interval) bool) {
		for i := range in {
			ApplyRanges(m.Ranges, i, yield)
		}
	}
}

func ApplyRanges(rs []Range, i Interval, yield func(Interval) bool) {
	if i.Len() == 0 {
		return
	}
	if len(rs) == 0 {
		yield(i)
		return
	}
	r, rs := rs[0], rs[1:]
	left, in, right := r.Split(i)
	ApplyRanges(rs, left, yield)
	if !in.Empty() && !yield(r.Apply(in)) {
		return
	}
	ApplyRanges(rs, right, yield)
}
