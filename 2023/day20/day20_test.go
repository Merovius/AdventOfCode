package main

import (
	_ "embed"
	"testing"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func TestPart1(t *testing.T) {
	tcs := []struct {
		name string
		in   string
		want int
	}{
		{"example1", example1, 32000000},
		{"example2", example2, 11687500},
		{"input", input, 680278040},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in); got != tc.want {
				t.Errorf("Part1(…) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	// Part2 does not work for the examples.
	in, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	want := 243548140870057
	got := Part2(in)
	if got != want {
		t.Errorf("Part2(…) = %d, want %d", got, want)
	}
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part1(in) != 680278040 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part1(in) != 680278040 {
				b.Fail()
			}
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part2(in) != 243548140870057 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part2(in) != 243548140870057 {
				b.Fail()
			}
		}
	})
}
