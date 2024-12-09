package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example []byte

//go:embed example2.txt
var example2 []byte

//go:embed input.txt
var input []byte

func Test(t *testing.T) {
	tcs := []struct {
		name  string
		in    []byte
		want1 int
		want2 int
	}{
		{"example", example, 1928, 2858},
		{"example2", example2, 2353, 3158},
		{"input", input, 6200294120911, 6227018762750},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if got1 := Part1(tc.in); got1 != tc.want1 {
				t.Errorf("Part1(%q) = %d, want %d", tc.name, got1, tc.want1)
			}
			if got2 := Part2(tc.in); got2 != tc.want2 {
				t.Errorf("Part2(%q) = %d, want %d", tc.name, got2, tc.want2)
			}
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	for range b.N {
		Part1(input)
	}
}

func BenchmarkPart2(b *testing.B) {
	for range b.N {
		Part2(input)
	}
}
