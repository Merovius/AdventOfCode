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

// MapReduce runs a MapReduce pipeline. It runs m (the map phase) in parallel
// to split the sequence of inputs into key/value pairs. It then spawns one
// goroutine running r per key (the reduce phase). Lastly, the outputs of all
// reducers are output as a sequence.
func MapReduce[A, B, C any, K comparable](seq iter.Seq[A], m func(A) iter.Seq2[K, B], r func(K, iter.Seq[B]) C) iter.Seq2[K, C] {
	return func(yield func(K, C) bool) {
		N := runtime.NumCPU()

		chA := make(chan A, N)
		chKB := make(chan pair[K, B], N)
		chKC := make(chan pair[K, C])
		stop := make(chan struct{})

		// map
		wg := new(sync.WaitGroup)
		wg.Add(N)
		for range N {
			go func(wg *sync.WaitGroup, in <-chan A, m func(A) iter.Seq2[K, B], out chan<- pair[K, B]) {
				defer wg.Done()
				for a := range in {
					for k, b := range m(a) {
						out <- pair[K, B]{k, b}
					}
				}
			}(wg, chA, m, chKB)
		}
		go func(wg *sync.WaitGroup, thenClose chan<- pair[K, B]) {
			wg.Wait()
			close(thenClose)
		}(wg, chKB)

		// shuffle
		wg = new(sync.WaitGroup)
		wg.Add(1)
		go func(wg *sync.WaitGroup, in <-chan pair[K, B], r func(K, iter.Seq[B]) C, chKC chan<- pair[K, C]) {
			defer wg.Done()
			reducers := make(map[K]chan<- B)
			for kb := range in {
				k, b := kb.a, kb.b
				if chB, ok := reducers[k]; ok {
					chB <- b
					continue
				}

				// reduce
				chB := make(chan B)
				defer close(chB)
				wg.Add(1)
				go func(wg *sync.WaitGroup, k K, in <-chan B, r func(K, iter.Seq[B]) C, out chan<- pair[K, C]) {
					defer wg.Done()
					c := r(k, func(yield func(B) bool) {
						var stopped bool
						for b := range in {
							if !stopped && !yield(b) {
								// TODO: figure out if we want to do something else here.
								stopped = true
							}
						}
					})
					out <- pair[K, C]{k, c}
				}(wg, k, chB, r, chKC)
				reducers[k] = chB
				chB <- b
			}
		}(wg, chKB, r, chKC)
		go func(wg *sync.WaitGroup, thenClose chan<- pair[K, C]) {
			wg.Wait()
			close(thenClose)
		}(wg, chKC)

		// consume input sequence
		go func(stop <-chan struct{}, seq iter.Seq[A], out chan<- A) {
			defer close(out)
			for a := range seq {
				select {
				case out <- a:
				case <-stop:
					return
				}
			}
		}(stop, seq, chA)

		var stopped bool
		for kc := range chKC {
			if !stopped && !yield(kc.a, kc.b) {
				close(stop)
				stopped = true
			}
		}
	}
}
