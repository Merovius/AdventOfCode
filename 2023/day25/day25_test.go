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
		name string
		in   string
		want int
	}{
		{"example", example, 54},
		{"input", input, 495607},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			in, err := Parse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			t.Run("Part1", func(t *testing.T) {
				if got1 := Part1(in); got1 != tc.want {
					t.Errorf("Part1(…) = %v, want %v", got1, tc.want)
				}
			})
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			if Part1(in) != 495607 {
				b.Fail()
			}
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if Part1(in) != 495607 {
				b.Fail()
			}
		}
	})
}
