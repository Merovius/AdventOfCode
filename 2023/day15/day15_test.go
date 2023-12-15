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
		{"example", example, 1320, 145},
		{"input", input, 513643, 265345},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			t.Run("Part1", func(t *testing.T) {
				if got1 := Part1(in); got1 != tc.wantPart1 {
					t.Errorf("Part1(…) = %v, want %v", got1, tc.wantPart1)
				}
			})
			t.Run("Part2", func(t *testing.T) {
				if got2 := Part2(in); got2 != tc.wantPart2 {
					t.Errorf("Part2(…) = %v, want %v", got2, tc.wantPart2)
				}
			})
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part1(in) != 513643 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part1(in) != 513643 {
				b.Fail()
			}
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part2(in) != 265345 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part2(in) != 265345 {
				b.Fail()
			}
		}
	})
}
