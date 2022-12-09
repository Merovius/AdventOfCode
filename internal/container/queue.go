package container

// FIFO is a First-In-First-Out queue.
type FIFO[E any] []E

// Len returns the length of the queue.
func (q *FIFO[E]) Len() int {
	return len(*q)
}

// Push an element unto the queue.
func (q *FIFO[E]) Push(e E) {
	*q = append(*q, e)
}

// Pop an element from the queue.
func (q *FIFO[E]) Pop() (e E) {
	e, *q = (*q)[0], (*q)[1:]
	return e
}

// Peek at the next element in the queue.
func (q *FIFO[E]) Peek() (e E) {
	return (*q)[0]
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
