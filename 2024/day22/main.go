package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"slices"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/pseq"
	"gonih.org/AdventOfCode/internal/set"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(in []byte) ([]int, error) {
	return parse.Slice(split.Lines, parse.Signed[int])(string(in))
}

func Part1(in []int) int {
	return pseq.MapMerge(slices.Values(in), hash, func(a, b int) int { return a + b })
}

func hash(n int) int {
	for range 2000 {
		n = evolve(n)
	}
	return n
}

func evolve(n int) int {
	// All constants in the problem text are powers of 2, so we can
	// implement the operations using shifts and bit masks.
	const M = (1 << 24) - 1
	n ^= n << 6
	n &= M
	n ^= n >> 5
	n &= M
	n ^= n << 11
	n &= M
	return n
}

func Part2(nums []int) int {
	m := make(map[Δ]int)
	for _, n := range nums {
		for δ, v := range AggregatePrices(n) {
			m[δ] += v
		}
	}
	var best int
	for _, v := range m {
		best = max(best, v)
	}
	return best
}

func AggregatePrices(n int) iter.Seq2[Δ, int] {
	return func(yield func(Δ, int) bool) {
		seen := make(set.Set[Δ])
		var δ Δ
		for j := range 2000 {
			m := evolve(n)
			δ = δ.Shift(int8(m%10 - n%10))
			if j > 3 && !seen.Contains(δ) {
				if !yield(δ, m%10) {
					return
				}
				seen.Add(δ)
			}
			n = m
		}
	}
}

// Δ is a Packed vector of the last four changes.
type Δ uint32

func MakeΔ(a, b, c, d int8) Δ {
	return Δ(a+9)<<0 |
		Δ(b+9)<<5 |
		Δ(c+9)<<10 |
		Δ(d+9)<<15
}

// Shift v∈[-9,9] into the 0'th component.
func (δ Δ) Shift(v int8) Δ {
	δ = (δ<<5)&0xfffff | Δ(v+9)
	return δ
}

// Packed is a Δ and a price, Packed into a single uint32.
type Packed uint32

func Pack(δ Δ, price uint8) Packed {
	return Packed(δ) | (Packed(price)&0x1f)<<20
}

func (p Packed) unpack() (δ Δ, price uint8) {
	return Δ(p & 0xfffff), uint8(p >> 20)
}
