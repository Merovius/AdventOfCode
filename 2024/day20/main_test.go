package main

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed example.txt
var example []byte

//go:embed input.txt
var input []byte

func TestPart1(t *testing.T) {
	tcs := []struct {
		name  string
		input []byte
		n     int
		want  int
	}{
		{"example", example, 2, 44},
		{"example", example, 4, 30},
		{"example", example, 6, 16},
		{"example", example, 8, 14},
		{"example", example, 10, 10},
		{"example", example, 12, 8},
		{"example", example, 20, 5},
		{"example", example, 36, 4},
		{"example", example, 38, 3},
		{"example", example, 40, 2},
		{"example", example, 64, 1},
		{"input", input, 100, 1399},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%s%d", tc.name, tc.n), func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in, tc.n); got != tc.want {
				t.Errorf("Part1(%q, %d) = %v, want %v", tc.name, tc.n, got, tc.want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tcs := []struct {
		name  string
		input []byte
		n     int
		want  int
	}{
		{"example", example, 76, 3},
		{"example", example, 74, 7},
		{"example", example, 72, 29},
		{"example", example, 70, 41},
		{"example", example, 68, 55},
		{"example", example, 66, 67},
		{"example", example, 64, 86},
		{"example", example, 62, 106},
		{"example", example, 60, 129},
		{"example", example, 58, 154},
		{"example", example, 56, 193},
		{"example", example, 54, 222},
		{"example", example, 52, 253},
		{"example", example, 50, 285},
		{"input", input, 100, 994807},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%s%d", tc.name, tc.n), func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part2(in, tc.n); got != tc.want {
				t.Errorf("Part2(%q, %d) = %v, want %v", tc.name, tc.n, got, tc.want)
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
		Part1(in, 100)
	}
}

func BenchmarkPart2(b *testing.B) {
	in, _ := Parse(input)
	b.ResetTimer()
	for range b.N {
		Part2(in, 100)
	}
}
