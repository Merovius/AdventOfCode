package main

import (
	_ "embed"
	"strings"
	"testing"

	"gonih.org/AdventOfCode/internal/grid"
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
		{"example", example, 4361, 467835},
		{"input", input, 531932, 73646890},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g, err := grid.Read[rune](strings.NewReader(tc.in), func(r rune) (rune, error) { return r, nil })
			if err != nil {
				t.Fatal(err)
			}
			if got := Part1(g); got != tc.wantPart1 {
				t.Errorf("Part1() = %v, want %v", got, tc.wantPart1)
			}
			if got := Part2(g); got != tc.wantPart2 {
				t.Errorf("Part2() = %v, want %v", got, tc.wantPart2)
			}
		})
	}
}
