package container

type RingBuffer[E any] struct {
	b    []E
	r    int
	w    int
	full bool
}

// NewRingBuffer uses b as a ring buffer. The capacity of the buffer is cap(b)
// and b[:len(b)] are the initial contents. So make([]E, 0, N) is an empty
// buffer of size N.
func NewRingBuffer[E any](b []E) *RingBuffer[E] {
	return &RingBuffer[E]{
		b:    b[:cap(b)],
		r:    0,
		w:    len(b),
		full: len(b) == cap(b),
	}
}

func (b *RingBuffer[E]) inc(v int) int {
	v += 1
	if v == len(b.b) {
		v = 0
	}
	return v
}

// Len returns the number of elements currently in the buffer.
func (b *RingBuffer[E]) Len() int {
	if b.full {
		return len(b.b)
	}
	if b.r > b.w {
		return b.w - b.r + len(b.b)
	}
	return b.w - b.r
}

// Cap returns the capacity of the buffer.
func (b *RingBuffer[E]) Cap() int {
	return len(b.b)
}

// Push adds e to the buffer.
func (b *RingBuffer[E]) Push(e E) {
	if b.full {
		panic("Push to full RingBuffer")
	}
	b.b[b.w] = e
	b.w = b.inc(b.w)
	b.full = (b.w == b.r)
}

// Pop reads an element from the buffer.
func (b *RingBuffer[E]) Pop() E {
	if b.r == b.w && !b.full {
		panic("Pop from empty RingBuffer")
	}
	b.full = false
	e := b.b[b.r]
	b.r = b.inc(b.r)
	return e
}
