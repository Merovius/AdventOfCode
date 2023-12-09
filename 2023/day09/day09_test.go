package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

func Test(t *testing.T) {
	tcs := []struct {
		name      string
		in        string
		wantPart1 int
		wantPart2 int
	}{
		{"example", example, 114, 2},
		{"input", input, 2075724761, 1072},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got1, got2 := Solve(in); got1 != tc.wantPart1 || got2 != tc.wantPart2 {
				t.Errorf("Solve(_) = %v, %v, want %v, %v", got1, got2, tc.wantPart1, tc.wantPart2)
			}
		})
	}
}
