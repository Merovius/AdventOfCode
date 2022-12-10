package sync

import "sync"

type Once = sync.Once

type OnceValue[T any] struct {
	o sync.Once
	v T
}

func (o *OnceValue[T]) Do(f func() T) T {
	o.o.Do(func() {
		o.v = f()
	})
	return o.v
}

type OnceValues[T1, T2 any] struct {
	o  sync.Once
	v1 T1
	v2 T2
}

func (o *OnceValues[T1, T2]) Do(f func() (T1, T2)) (T1, T2) {
	o.o.Do(func() {
		o.v1, o.v2 = f()
	})
	return o.v1, o.v2
}
