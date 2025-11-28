package set

import (
	"fmt"
	"iter"
	"math/bits"
)

type negativeElementError int

func (e negativeElementError) Error() string {
	return fmt.Sprintf("negative element %d", int(e))
}

// Bits represents a dense set of non-negative integers. Its zero value
// represents the empty set.
type Bits struct {
	e []uint
}

func MakeBits(n int) Bits {
	N := (n + (bits.UintSize - 1)) / bits.UintSize
	return Bits{
		e: make([]uint, N),
	}
}

func (b *Bits) grow(e int) {
	if e < 0 {
		panic(negativeElementError(e))
	}
	i := e / bits.UintSize
	for i >= len(b.e) {
		b.e = append(b.e, 0)
		b.e = b.e[:cap(b.e)]
	}
}

func (b *Bits) Len() int {
	var n int
	for _, e := range b.e {
		n += bits.OnesCount(e)
	}
	return n
}

func (b *Bits) Add(e int) {
	b.grow(e)
	i, o := e/bits.UintSize, e%bits.UintSize
	b.e[i] |= 1 << o
}

func (b *Bits) Contains(e int) bool {
	if e < 0 {
		panic(negativeElementError(e))
	}
	i, o := e/bits.UintSize, e%bits.UintSize
	if i >= len(b.e) {
		return false
	}
	return b.e[i]&(1<<o) != 0
}

func (b *Bits) Delete(e int) {
	if e < 0 {
		panic(negativeElementError(e))
	}
	i, o := e/bits.UintSize, e%bits.UintSize
	if i >= len(b.e) {
		return
	}
	b.e[i] &^= 1 << o
}

func (b *Bits) Clear() {
	for i := range b.e {
		b.e[i] = 0
	}
}

// All is equivalent to Sorted.
func (b *Bits) All() iter.Seq[int] {
	return b.Sorted()
}

// Sorted yields all elements in b in ascending order. It is O(m), where m is the
// maximum element ever inserted into b.
func (b *Bits) Sorted() iter.Seq[int] {
	return func(yield func(int) bool) {
		var v int
		for _, e := range b.e {
			i := 0
			for e != 0 {
				j := bits.TrailingZeros(e)
				i += j
				if !yield(v | i) {
					return
				}
				e >>= j + 1
				i += 1
			}
			v += bits.UintSize
		}
	}
}

// Descending yields all elements in b in descending order. It is O(m), where m is the
// maximum element ever inserted into b.
func (b *Bits) Descending() iter.Seq[int] {
	return func(yield func(int) bool) {
		v := (len(b.e) - 1) * bits.UintSize
		for n := len(b.e) - 1; n >= 0; n-- {
			e := b.e[n]
			i := bits.UintSize - 1
			for e != 0 {
				j := bits.LeadingZeros(e)
				i -= j
				if !yield(v | i) {
					return
				}
				e <<= j + 1
				i -= 1
			}
			v -= bits.UintSize
		}
	}
}
