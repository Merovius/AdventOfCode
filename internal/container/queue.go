package container

// FIFO is a First-In-First-Out queue. A zero FIFO is an empty queue ready to use.
type FIFO[E any] struct {
	s []E
	r int // index to read
	w int // index to write, of len(s) if queue is full.
}

// MakeFIFO returns a FIFO that can hold up to size elements without allocating.
func MakeFIFO[E any](size int) FIFO[E] {
	return FIFO[E]{s: make([]E, size)}
}

// Len returns the length of the queue.
func (q *FIFO[E]) Len() int {
	if q.w == len(q.s) {
		return len(q.s)
	} else if q.w < q.r {
		return q.w - q.r + len(q.s)
	}
	return q.w - q.r
}

// Push an element unto the queue.
func (q *FIFO[E]) Push(e E) {
	if q.w == len(q.s) {
		s := append(q.s, e)
		i := copy(s, q.s[q.r:])
		copy(s[i:], q.s[:q.r])
		q.r, q.s = 0, s[:cap(s)]
	}
	q.s[q.w] = e
	q.w++
	if q.w == len(q.s) {
		q.w = 0
	}
	if q.w == q.r {
		q.w = len(q.s)
	}
}

// Pop an element from the queue.
func (q *FIFO[E]) Pop() (e E) {
	if q.w == q.r {
		panic("Pop from empty FIFO")
	}
	if q.w == len(q.s) {
		q.w = q.r
	}
	e, q.r = q.s[q.r], (q.r + 1)
	if q.r == len(q.s) {
		q.r = 0
	}
	return e
}

// Peek at the next element in the queue.
func (q *FIFO[E]) Peek() (e E) {
	if q.w == q.r {
		panic("Peek into empty FIFO")
	}
	return q.s[q.r]
}

// LIFO is a Last-In-First-Out queue (also known as a stack).
type LIFO[E any] []E

// Len returns the length of the queue.
func (q *LIFO[E]) Len() int {
	return len(*q)
}

// Push an element unto the queue.
func (q *LIFO[E]) Push(e E) {
	*q = append(*q, e)
}

// Pop an element from the queue.
func (q *LIFO[E]) Pop() (e E) {
	*q, e = (*q)[:len(*q)-1], (*q)[len(*q)-1]
	return e
}

func (q *LIFO[E]) Peek() (e E) {
	return (*q)[len(*q)-1]
}
