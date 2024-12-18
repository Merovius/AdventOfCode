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
		w     int
		h     int
		n     int
		want1 int
		want2 [2]int
	}{
		{"example", example, 7, 7, 12, 22, [2]int{6, 1}},
		{"input", input, 71, 71, 1024, 314, [2]int{15, 20}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in, tc.w, tc.h, tc.n); got != tc.want1 {
				t.Errorf("Part1(%q) = %v, want %v", tc.name, got, tc.want1)
			}
			if got := Part2(in, tc.w, tc.h); got != tc.want2 {
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
		Part1(in, 71, 71, 1024)
	}
}

func BenchmarkPart2(b *testing.B) {
	in, _ := Parse(input)
	b.ResetTimer()
	for range b.N {
		Part2(in, 71, 71)
	}
}
