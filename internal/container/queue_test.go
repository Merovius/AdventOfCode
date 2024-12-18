package container

import (
	"math/rand/v2"
	"testing"
)

func TestFIFO(t *testing.T) {
	for n := range uint64(1000) {
		var q FIFO[int]
		rnd := rand.New(rand.NewPCG(0, n))
		for i := range rnd.IntN(10000) {
			if q.Len() == 0 {
				n := rnd.IntN(10) + 1
				for j := range n {
					q.Push(i + j)
				}
			}
			if got := q.Pop(); got != i {
				t.Errorf("q.Pop() = %d, want %d", got, i)
			}
		}
	}
}

// TestMakeFIFO checks that MakeFIFO reserves sufficient space to avoid allocations.
func TestMakeFIFO(t *testing.T) {
	const N = 128
	rnd := rand.New(rand.NewPCG(0, 0))
	got := testing.AllocsPerRun(1000, func() {
		q := MakeFIFO[int](N)
		for range rnd.IntN(N) {
			q.Push(rnd.IntN(N))
		}
		for q.Len() > 0 {
			q.Pop()
		}
		for range N {
			q.Push(rnd.IntN(N))
		}
	})
	if got != 1 {
		t.Errorf("MakeFIFO[int](128) does not reserve enough space (AllocsPerRun = %v)", got)
	}
}

func BenchmarkFIFO(b *testing.B) {
	b.Run("PushPop", func(b *testing.B) {
		q := MakeFIFO[int](10)
		for b.Loop() {
			q.Push(0)
			q.Pop()
		}
	})
}
