package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example []byte

//go:embed input.txt
var input []byte

func Test(t *testing.T) {
	tcs := []struct {
		name  string
		input []byte
		want1 int
		want2 int
	}{
		{"example", example, 126384, 154115708116294},
		{"input", input, 217662, 263617786809000},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in); got != tc.want1 {
				t.Errorf("Part1(%q) = %v, want %v", tc.name, got, tc.want1)
			}
			if got := Part2(in); got != tc.want2 {
				t.Errorf("Part2(%q) = %v, want %v", tc.name, got, tc.want2)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	for range b.N {
		Parse(input)
	}
}

func BenchmarkPart1(b *testing.B) {
	in, _ := Parse(input)
	b.ResetTimer()
	for range b.N {
		Part1(in)
	}
}

func BenchmarkPart2(b *testing.B) {
	in, _ := Parse(input)
	b.ResetTimer()
	for range b.N {
		Part2(in)
	}
}
