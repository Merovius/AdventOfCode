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
		{"example", example, 62, 952408144115},
		{"input", input, 50603, 96556251590677},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			t.Run("Part1", func(t *testing.T) {
				if got1 := Enclosed(part1); got1 != tc.wantPart1 {
					t.Errorf("Enclosed(part1) = %v, want %v", got1, tc.wantPart1)
				}
			})
			t.Run("Part2", func(t *testing.T) {
				if got2 := Enclosed(part2); got2 != tc.wantPart2 {
					t.Errorf("Enclosed(part2) = %v, want %v", got2, tc.wantPart2)
				}
			})
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			part1, _, _ := Parse(input)
			if Enclosed(part1) != 50603 {
				b.Fail()
			}
		}
	})
	part1, _, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Enclosed(part1) != 50603 {
				b.Fail()
			}
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, part2, _ := Parse(input)
			if Enclosed(part2) != 96556251590677 {
				b.Fail()
			}
		}
	})
	_, part2, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Enclosed(part2) != 96556251590677 {
				b.Fail()
			}
		}
	})
}
