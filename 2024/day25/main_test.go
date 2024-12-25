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
		want  int
	}{
		{"example", example, 3},
		{"input", input, 2978},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in); got != tc.want {
				t.Errorf("Part1(%q) = %v, want %v", tc.name, got, tc.want)
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
