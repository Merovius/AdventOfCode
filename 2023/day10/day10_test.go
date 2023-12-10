package main

import (
	_ "embed"
	"testing"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed example3.txt
var example3 string

//go:embed example4.txt
var example4 string

//go:embed input.txt
var input string

func Test(t *testing.T) {
	tcs := []struct {
		name      string
		in        string
		wantPart1 int
		wantPart2 int
	}{
		{"example1", example1, 4, 1},
		{"example2", example2, 8, 1},
		{"example3", example3, 70, 8},
		{"example4", example4, 80, 10},
		{"input", input, 6956, 455},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got1 := Part1(in); got1 != tc.wantPart1 {
				t.Errorf("Part1(…) = %v, want %v", got1, tc.wantPart1)
			}
			if got2 := Part2(in); got2 != tc.wantPart2 {
				t.Errorf("Part2(…) = %v, want %v", got2, tc.wantPart2)
			}
		})
	}
}
