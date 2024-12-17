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
		input []byte
		want1 string
		want2 int
	}{
		{"example", example, "4,6,3,5,6,3,5,2,1,0", -1},
		{"example2", example2, "5,7,3,0", 117440},
		{"input", input, "1,7,2,1,4,1,5,4,0", 37221261688308},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in); got != tc.want1 {
				t.Errorf("Part1(%q) = %q, want %q", tc.name, got, tc.want1)
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
