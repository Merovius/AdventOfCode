package main

import (
	_ "embed"
	"testing"
)

//go:embed input.txt
var input string

func Test(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		want1 int
		want2 int
	}{
		{"example1", "(())", 0, 0},
		{"example2", "()()", 0, 0},
		{"example3", "(((", 3, 0},
		{"example4", "(()(()(", 3, 0},
		{"example5", "))(((((", 3, 1},
		{"example6", "())", -1, 3},
		{"example7", "))(", -1, 1},
		{"example8", ")))", -3, 1},
		{"example9", ")())())", -3, 1},
		{"input", input, 280, 1797},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if got := Part1(tc.input); got != tc.want1 {
				t.Errorf("Part1(%q) = %v, want %v", tc.name, got, tc.want1)
			}
			if got := Part2(tc.input); got != tc.want2 {
				t.Errorf("Part2(%q) = %v, want %v", tc.name, got, tc.want2)
			}
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		Part1(input)
	}
}

func BenchmarkPart2(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		Part2(input)
	}
}
