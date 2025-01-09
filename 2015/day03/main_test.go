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
		{"example1", ">", 2, 2},
		{"example2", "^>v<", 4, 3},
		{"example3", "^v^v^v^v^v", 2, 11},
		{"example2", "^v", 2, 3},
		{"input", input, 2081, 2341},
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
	for b.Loop() {
		Part1(input)
	}
}

func BenchmarkPart2(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		Part2(input)
	}
}
