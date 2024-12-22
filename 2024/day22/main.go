package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"slices"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/math"
	"github.com/Merovius/AdventOfCode/internal/pseq"
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
	m := Preprocess(nums)
	return pseq.MapMerge(Δs(), m.Match, math.Max)
}

// Monkey keeps preprocessed data to speed up matching change-sequences.
type Monkey [][1997]Packed

func Preprocess(in []int) Monkey {
	m := make(Monkey, len(in))
	pseq.Each2(slices.All(in), func(i, n int) {
		var δ Δ
		for j := range 2000 {
			k := evolve(n)
			δ = δ.Shift(int8(k%10 - n%10))
			if j > 3 {
				m[i][j-3] = Pack(δ, uint8(k%10))
			}
			n = k
		}
	})
	return m
}

func (m Monkey) Match(δ Δ) int {
	var N int
	for _, ps := range m {
		for _, p := range ps {
			if d, price := p.unpack(); d == δ {
				N += int(price)
				break
			}
		}
	}
	return N
}

// Δ is a Packed vector of the last four changes.
type Δ uint32

func MakeΔ(a, b, c, d int8) Δ {
	return Δ(a+9)<<0 |
		Δ(b+9)<<5 |
		Δ(c+9)<<10 |
		Δ(d+9)<<15
}

// Δs iterates over all possible vectors of four changes.
func Δs() iter.Seq[Δ] {
	return func(yield func(Δ) bool) {
		for a := int8(-9); a <= 9; a++ {
			for b := int8(-9); b <= 9; b++ {
				for c := int8(-9); c <= 9; c++ {
					for d := int8(-9); d <= 9; d++ {
						if !yield(MakeΔ(a, b, c, d)) {
							return
						}
					}
				}
			}
		}
	}
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
