package main

import (
	_ "embed"
	"math/rand"
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
		{"example", example, 6440, 5905},
		{"input", input, 246409899, 244848487},
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

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}

func BenchmarkPart1(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			Part1(in)
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			rand.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
			b.StartTimer()
			Part1(in)
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	b.Run("WithParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			in, _ := Parse(input)
			Part2(in)
		}
	})
	in, _ := Parse(input)
	b.Run("WithoutParse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			rand.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
			b.StartTimer()
			Part2(in)
		}
	})
}
