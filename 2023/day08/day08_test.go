package main

import (
	_ "embed"
	"testing"
)

//go:embed example1.txt
var example1 []byte

//go:embed example2.txt
var example2 []byte

//go:embed example3.txt
var example3 []byte

//go:embed input.txt
var input []byte

func TestPart1(t *testing.T) {
	tcs := []struct {
		name string
		in   []byte
		want int
	}{
		{"example1", example1, 2},
		{"example2", example2, 6},
		{"input", input, 20569},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in); got != tc.want {
				t.Errorf("Part1(…) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tcs := []struct {
		name string
		in   []byte
		want int
	}{
		{"example3", example3, 6},
		{"input", input, 21366921060721},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part2(in); got != tc.want {
				t.Errorf("Part1(…) = %v, want %v", got, tc.want)
			}
		})
	}
}
