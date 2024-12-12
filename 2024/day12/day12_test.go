package main

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed example.txt
var example []byte

//go:embed example2.txt
var example2 []byte

//go:embed example3.txt
var example3 []byte

//go:embed example4.txt
var example4 []byte

//go:embed example5.txt
var example5 []byte

//go:embed input.txt
var input []byte

func Test(t *testing.T) {
	tcs := []struct {
		name  string
		input []byte
		want1 int
		want2 int
	}{
		{"example", example, 140, 80},
		{"example2", example2, 772, 436},
		{"example3", example3, 1930, 1206},
		{"example4", example4, 692, 236},
		{"example5", example5, 1184, 368},
		{"input", input, 1451030, 859494},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(bytes.NewReader(tc.input))
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
		Parse(bytes.NewReader(input))
	}
}

func BenchmarkPart1(b *testing.B) {
	in, _ := Parse(bytes.NewReader(input))
	b.ResetTimer()
	for range b.N {
		Part1(in)
	}
}

func BenchmarkPart2(b *testing.B) {
	in, _ := Parse(bytes.NewReader(input))
	b.ResetTimer()
	for range b.N {
		Part2(in)
	}
}
