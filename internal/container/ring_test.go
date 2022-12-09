package container

import (
	"fmt"
	"testing"
)

func TestRingBufferBasic(t *testing.T) {
	b := NewRingBuffer(make([]int, 0, 4))
	if l, c := b.Len(), b.Cap(); l != 0 || c != 4 {
		t.Fatalf("Len(), Cap() = %d, %d, want 0, 4", l, c)
	}
	const (
		push = iota
		pop
	)
	wantLen := 0

	tcs := []struct {
		op  int
		val int
	}{
		{push, 1},
		{pop, 1},
		{push, 2},
		{push, 3},
		{push, 4},
		{push, 5},
		{pop, 2},
		{pop, 3},
		{push, 6},
		{pop, 4},
		{pop, 5},
		{pop, 6},
	}
	for _, tc := range tcs {
		switch tc.op {
		case push:
			t.Logf("Push(%d)", tc.val)
			b.Push(tc.val)
			wantLen++
		case pop:
			if got := b.Pop(); got != tc.val {
				t.Errorf("Pop() = %d, want %d", got, tc.val)
			} else {
				t.Logf("Pop() = %d", tc.val)
			}
			wantLen--
		default:
			panic(fmt.Sprintf("invalid op %v", tc.op))
		}
		if l := b.Len(); l != wantLen {
			t.Errorf("Len() = %d, want %d", l, wantLen)
		}
	}
}

func TestRingBufferPanic(t *testing.T) {
	b := NewRingBuffer(make([]int, 0, 4))
	if !panics(t, func() { b.Pop() }) {
		t.Fatalf("Pop() of empty buffer does not panic")
	}
	ringBufferPush(t, b, 1, 2, 3, 4)
	if !panics(t, func() { b.Push(5) }) {
		t.Fatalf("Push() to full buffer does not panic")
	}
	ringBufferPop(t, b, 4)
	if !panics(t, func() { b.Pop() }) {
		t.Fatalf("Pop() of emptied buffer does not panic")
	}
	// check when filling-point does not align with buffer-size
	ringBufferPush(t, b, 1, 2)
	ringBufferPop(t, b, 2)
	if !panics(t, func() { b.Pop() }) {
		t.Fatalf("Pop() of emptied buffer does not panic")
	}
	ringBufferPush(t, b, 1, 2, 3, 4)
	if !panics(t, func() { b.Push(5) }) {
		t.Fatalf("Push() of filled buffer does not panic")
	}
}

func ringBufferPush[E any](t *testing.T, b *RingBuffer[E], e ...E) {
	for _, e := range e {
		if panics(t, func() { b.Push(e) }) {
			t.Fatalf("Push(%v) panics", e)
		}
	}
}

func ringBufferPop[E any](t *testing.T, b *RingBuffer[E], n int) []E {
	var out []E
	for i := 0; i < n; i++ {
		if panics(t, func() { out = append(out, b.Pop()) }) {
			t.Fatal("Pop() panics")
		}
	}
	return out
}

func panics(t *testing.T, f func()) (b bool) {
	defer func() {
		b = (recover() != nil)
	}()
	f()
	return false
}
