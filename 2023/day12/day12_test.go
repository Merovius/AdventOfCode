package main

import (
	_ "embed"
	"strconv"
	"testing"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

func Test(t *testing.T) {
	tcs := []struct {
		name      string
		in        string
		wantPart1 int
		wantPart2 int
	}{
		{"example", example, 21, 525152},
		{"input", input, 6852, 8475948826693},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			t.Run("Part1", func(t *testing.T) {
				if got1 := Part1(in); got1 != tc.wantPart1 {
					t.Errorf("Part1(…) = %v, want %v", got1, tc.wantPart1)
				}
			})
			t.Run("Part2", func(t *testing.T) {
				if got2 := Part2(in); got2 != tc.wantPart2 {
					t.Errorf("Part2(…) = %v, want %v", got2, tc.wantPart2)
				}
			})
		})
	}
}

//go:embed testdata/part1.txt
var part1 string

//go:embed testdata/part2.txt
var part2 string

func TestArrangements(t *testing.T) {
	type TestCase struct {
		Springs string
		Groups  []int
		Want    int
	}
	p := parse.TrimSpace(parse.Lines(
		parse.Struct[TestCase](
			split.Fields,
			parse.String[string],
			parse.Slice(split.On(","), parse.Signed[int]),
			parse.Signed[int],
		),
	))
	t.Run("Part1", func(t *testing.T) {
		tcs, err := p(part1)
		if err != nil {
			t.Fatal(err)
		}
		for i, tc := range tcs {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				if got := Arrangements(tc.Springs, tc.Groups); got != tc.Want {
					t.Errorf("Arrangements(%q, %v) = %d, want %d", tc.Springs, tc.Groups, got, tc.Want)
				}
			})
		}
	})
	t.Run("Part2", func(t *testing.T) {
		tcs, err := p(part2)
		if err != nil {
			t.Fatal(err)
		}
		for i, tc := range tcs {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				if got := Arrangements2(tc.Springs, tc.Groups); got != tc.Want {
					t.Errorf("Arrangements2(%q, %v) = %d, want %d", tc.Springs, tc.Groups, got, tc.Want)
				}
			})
		}
	})
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			clear(memo)
			in, _ := Parse(input)
			if Part1(in) != 6852 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			clear(memo)
			if Part1(in) != 6852 {
				b.Fail()
			}
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			clear(memo)
			in, _ := Parse(input)
			if Part2(in) != 8475948826693 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			clear(memo)
			if Part2(in) != 8475948826693 {
				b.Fail()
			}
		}
	})
}
