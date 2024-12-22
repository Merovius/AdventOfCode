// Package pseq contains utilities for parallelisation
package pseq

import (
	"iter"
	"runtime"
	"sync"
)

type pair[A, B any] struct {
	a A
	b B
}

// MapMerge runs f in parallel, merging the results using m. m must be
// associative and commutative.
func MapMerge[A, B any](seq iter.Seq[A], f func(A) B, m func(B, B) B) B {
	N := runtime.NumCPU()
	chA := make(chan A, N)
	// use pointers, to prevent false sharing
	bs := make([]*B, N)

	wg := new(sync.WaitGroup)
	wg.Add(N)
	for i := range N {
		bs[i] = new(B)
		go func(in <-chan A, out *B) {
			defer wg.Done()
			var first = true
			for a := range chA {
				if first {
					*out = f(a)
					first = false
				} else {
					*out = m(f(a), *out)
				}
			}
		}(chA, bs[i])
	}
	for a := range seq {
		chA <- a
	}
	close(chA)
	wg.Wait()
	b := *bs[0]
	for i := 1; i < N; i++ {
		b = m(b, *bs[i])
	}
	return b
}

// Each runs f in parallel, passing in all elements from seq.
func Each[A any](seq iter.Seq[A], f func(A)) {
	N := runtime.NumCPU()
	chA := make(chan A, N)
	wg := new(sync.WaitGroup)
	wg.Add(N)
	for range N {
		go func(wg *sync.WaitGroup, in <-chan A, f func(A)) {
			defer wg.Done()
			for a := range in {
				f(a)
			}
		}(wg, chA, f)
	}
	for a := range seq {
		chA <- a
	}
	close(chA)
	wg.Wait()
}

// Each2 runs f in parallel, passing in all pairs from seq.
func Each2[A, B any](seq iter.Seq2[A, B], f func(A, B)) {
	// not simply using Each[pair[A,B]], as it requires allocating a
	// closure.
	N := runtime.NumCPU()
	chP := make(chan pair[A, B], N)
	wg := new(sync.WaitGroup)
	wg.Add(N)
	for range N {
		go func(wg *sync.WaitGroup, in <-chan pair[A, B], f func(A, B)) {
			defer wg.Done()
			for p := range in {
				f(p.a, p.b)
			}
		}(wg, chP, f)
	}
	for a, b := range seq {
		chP <- pair[A, B]{a, b}
	}
	close(chP)
	wg.Wait()
}
