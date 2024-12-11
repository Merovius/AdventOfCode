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
		want1 int
		want2 int
	}{
		{"example", example, 55312, 65601038650482},
		{"input", input, 188902, 223894720281135},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(in); got != tc.want1 {
				t.Errorf("Part1(%v) = %v, want %v", in, got, tc.want1)
			}
			if got := Part2(in); got != tc.want2 {
				t.Errorf("Part2(%v) = %v, want %v", in, got, tc.want2)
			}
		})
	}
}
