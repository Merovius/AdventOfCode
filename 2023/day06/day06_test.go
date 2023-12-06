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
		{"example", example, 288, 71503},
		{"input", input, 633080, 20048741},
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
