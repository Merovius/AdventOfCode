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
		name      string
		in        []byte
		wantPart1 int
		wantPart2 int
	}{
		{"example", example, 35, 46},
		{"input", input, 379811651, 27992443},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got1 := Part1(cards); got1 != tc.wantPart1 {
				t.Errorf("Part1(…) = %v, want %v", got1, tc.wantPart1)
			}
			if got2 := Part2(cards); got2 != tc.wantPart2 {
				t.Errorf("Part2(…) = %v, want %v", got2, tc.wantPart2)
			}
		})
	}
}
