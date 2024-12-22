// Package pseq contains utilities for parallelisation
package pseq

import (
	"iter"
	"runtime"
	"sync"
)

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
