package priority_queue

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

type Q[T any] struct {
	less func(T, T) bool
	els  []T
}

func New[T constraints.Ordered]() *Q[T] {
	return NewFunc(func(a, b T) bool { return a < b })
}

func NewFunc[T any](less func(T, T) bool) *Q[T] {
	return &Q[T]{less: less}
}

func (q *Q[T]) Len() int {
	return len(q.els)
}

func (q *Q[T]) Push(v T) {
	heap.Push((*impl[T])(q), v)
}

func (q *Q[T]) Pop() T {
	return heap.Pop((*impl[T])(q)).(T)
}

type impl[T any] Q[T]

var _ heap.Interface = new(impl[int])

func (h *impl[T]) Len() int {
	return len(h.els)
}

func (h *impl[T]) Swap(i, j int) {
	h.els[i], h.els[j] = h.els[j], h.els[i]
}

func (h *impl[T]) Less(i, j int) bool {
	return h.less(h.els[i], h.els[j])
}

func (h *impl[T]) Push(x interface{}) {
	h.els = append(h.els, x.(T))
}

func (h *impl[T]) Pop() (v interface{}) {
	h.els, v = h.els[:len(h.els)-1], h.els[len(h.els)-1]
	return v
}
