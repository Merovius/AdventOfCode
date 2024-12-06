package set

import "iter"

type Set[E comparable] map[E]struct{}

func Make[E comparable](e ...E) Set[E] {
	s := make(Set[E])
	for _, e := range e {
		s.Add(e)
	}
	return s
}

func (s Set[E]) Add(e E) {
	s[e] = struct{}{}
}

func (s Set[E]) Delete(e E) {
	delete(s, e)
}

func (s Set[E]) Contains(e E) bool {
	_, ok := s[e]
	return ok
}

func (s Set[E]) Slice() []E {
	var out []E
	for e := range s {
		out = append(out, e)
	}
	return out
}

func (s Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if !yield(e) {
				return
			}
		}
	}
}

func Intersect[E comparable](s ...Set[E]) Set[E] {
	if len(s) == 0 {
		panic("empty intersection")
	}
	out := make(Set[E])
loop:
	for e := range s[0] {
		for _, s := range s[1:] {
			if !s.Contains(e) {
				continue loop
			}
		}
		out.Add(e)
	}
	return out
}

func Collect[E comparable](s iter.Seq[E]) Set[E] {
	out := make(Set[E])
	for e := range s {
		out.Add(e)
	}
	return out
}
