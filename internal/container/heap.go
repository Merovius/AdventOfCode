// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is largely copied from the standard library container/heap
// package, with minor edits to add type parameters.

package container

import "golang.org/x/exp/constraints"

func less[E constraints.Ordered](a, b E) bool { return a < b }

// Heap implements a min-heap.
type Heap[E constraints.Ordered] []E

// Init establishes the heap invariants required by the other routines in this
// package. Init is idempotent with respect to the heap invariants and may be
// called whenever the heap invariants may have been invalidated.
//
// The complexity is O(n) where n = h.Len().
func (h *Heap[E]) Init() {
	heapify(*h, less[E])
}

// Len returns the numuüber of elements in h.
func (h *Heap[E]) Len() int {
	return len(*h)
}

// Push pushes the element e onto the heap.
//
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Push(e E) {
	*h = heapPush(*h, less[E], e)
}

// Pop removes and returns the minimum element (according to Less) from the
// heap. Pop is equivalent to Remove(h, 0).
//
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Pop() E {
	var e E
	*h, e = heapPop(*h, less[E])
	return e
}

// Remove removes and returns the element at index i from the heap.
//
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Remove(i int) E {
	var e E
	*h, e = remove(*h, less[E], i)
	return e
}

// Fix re-establishes the heap ordering after the element at index i has
// changed its value. Changing the value of the element at index i and then
// calling Fix is equivalent to, but less expensive than, calling Remove(h, i)
// followed by a Push of the new value.
//
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Fix(i int) {
	heapFix(*h, less[E], i)
}

// HeapFunc implements a min-heap with custom comparison.
type HeapFunc[E any] struct {
	Elements []E
	Less     func(E, E) bool
}

// Init establishes the heap invariants required by the other routines in this
// package. Init is idempotent with respect to the heap invariants and may be
// called whenever the heap invariants may have been invalidated.
//
// The complexity is O(n) where n = h.Len().
func (h *HeapFunc[E]) Init() {
	heapify(h.Elements, h.Less)
}

// Len returns the numuüber of elements in h.
func (h *HeapFunc[E]) Len() int {
	return len(h.Elements)
}

// Push pushes the element e onto the heap.
//
// The complexity is O(log n) where n = h.Len().
func (h *HeapFunc[E]) Push(e E) {
	h.Elements = heapPush(h.Elements, h.Less, e)
}

// Pop removes and returns the minimum element (according to Less) from the
// heap. Pop is equivalent to Remove(h, 0).
//
// The complexity is O(log n) where n = h.Len().
func (h *HeapFunc[E]) Pop() E {
	var e E
	h.Elements, e = heapPop(h.Elements, h.Less)
	return e
}

// Remove removes and returns the element at index i from the heap.
//
// The complexity is O(log n) where n = h.Len().
func (h *HeapFunc[E]) Remove(i int) E {
	var e E
	h.Elements, e = remove(h.Elements, h.Less, i)
	return e
}

// Fix re-establishes the heap ordering after the element at index i has
// changed its value. Changing the value of the element at index i and then
// calling Fix is equivalent to, but less expensive than, calling Remove(h, i)
// followed by a Push of the new value.
//
// The complexity is O(log n) where n = h.Len().
func (h *HeapFunc[E]) Fix(i int) {
	heapFix(h.Elements, h.Less, i)
}

func heapify[E any](h []E, less func(E, E) bool) {
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		heapDown(h, less, i, n)
	}
}

func heapPush[S ~[]E, E any](h S, less func(E, E) bool, e E) S {
	h = append(h, e)
	heapUp(h, less, len(h)-1)
	return h
}

func heapPop[S ~[]E, E any](h S, less func(E, E) bool) (S, E) {
	n := len(h) - 1
	h[0], h[n] = h[n], h[0]
	heapDown(h, less, 0, n)
	return h[:n], h[n]
}

func remove[S ~[]E, E any](h S, less func(E, E) bool, i int) (S, E) {
	n := len(h) - 1
	if n != i {
		h[i], h[n] = h[n], h[i]
		if !heapDown(h, less, i, n) {
			heapUp(h, less, i)
		}
	}
	return h[:n], h[n]
}

func heapFix[E any](h []E, less func(E, E) bool, i int) {
	if !heapDown(h, less, i, len(h)) {
		heapUp(h, less, i)
	}
}

func heapUp[E any](h []E, less func(E, E) bool, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(h[j], h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		j = i
	}
}

func heapDown[E any](h []E, less func(E, E) bool, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(h[j2], h[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(h[j], h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		i = j
	}
	return i > i0
}
