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

func TestPart1(t *testing.T) {
	tcs := []struct {
		name  string
		input []byte
		want  int
	}{
		{"example", example, 4},
		{"example2", example2, 2024},
		{"input", input, 46362252142374},
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

func TestPart2(t *testing.T) {
	tcs := []struct {
		name  string
		input []byte
		want  string
	}{
		{"input", input, "cbd,gmh,jmq,qrh,rqf,z06,z13,z38"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part2(in); got != tc.want {
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

func BenchmarkPart2(b *testing.B) {
	in, _ := Parse(input)
	b.ResetTimer()
	for range b.N {
		Part2(in)
	}
}
