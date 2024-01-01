package main

import (
	_ "embed"
	"testing"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func TestPart1(t *testing.T) {
	tcs := []struct {
		name string
		in   string
		want int
	}{
		{"example", example1, 142},
		{"input", input, 55477},
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
		in   string
		want int
	}{
		{"example", example2, 281},
		{"input", input, 54431},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if got := Part2(in); got != tc.want {
				t.Errorf("Part2(…) = %v, want %v", got, tc.want)
			}
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part1(in) != 55477 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part1(in) != 55477 {
				b.Fail()
			}
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part2(in) != 54431 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part2(in) != 54431 {
				b.Fail()
			}
		}
	})
}
