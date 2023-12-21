package main

import (
	_ "embed"
	"testing"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

func TestPart1(t *testing.T) {
	type tc struct {
		n    int
		want int
	}

	tcs := []struct {
		name string
		in   string
		tcs  []tc
	}{
		{"example", example, []tc{
			{6, 16},
			{64, 42},
		}},
		{"input", input, []tc{{64, 3600}}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g, start, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			for _, tc := range tc.tcs {
				if got := NReachable(g, start, tc.n); got != tc.want {
					t.Errorf("NReachable(g, %v, %v) = %v, want %v", start, tc.n, got, tc.want)
				}
			}
		})
	}
}

func TestPart2(t *testing.T) {
	// Part2 does not work for the example.
	g, start, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	want := 599763113936220
	got := Part2(g, start)
	if got != want {
		t.Errorf("Part2(…) = %d, want %d", got, want)
	}
}
