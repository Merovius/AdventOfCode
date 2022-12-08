package set

import (
	"fmt"
	"strings"
)

type Set[E comparable] map[E]struct{}

func Make[E comparable](es ...E) Set[E] {
	s := make(Set[E])
	for _, e := range es {
		s.Add(e)
	}
	return s
}

func (s Set[E]) Add(e E) {
	s[e] = struct{}{}
}

func (s Set[E]) Contains(e E) bool {
	_, ok := s[e]
	return ok
}

func (s Set[E]) Delete(e E) bool {
	if _, ok := s[e]; ok {
		delete(s, e)
		return true
	}
	return false
}

func (s Set[E]) String() string {
	w := new(strings.Builder)
	w.WriteByte('{')
	first := true
	for e := range s {
		if !first {
			w.WriteString(", ")
		}
		first = false
		fmt.Fprint(w, e)
	}
	w.WriteByte('}')
	return w.String()
}

func Union[E comparable](s1, s2 Set[E]) Set[E] {
	s := make(Set[E])
	for e := range s1 {
		s.Add(e)
	}
	for e := range s2 {
		s.Add(e)
	}
	return s
}

func Intersection[E comparable](s1, s2 Set[E]) Set[E] {
	s := make(Set[E])
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	for e := range s1 {
		if s2.Contains(e) {
			s.Add(e)
		}
	}
	return s
}

func Map[E1, E2 comparable](s Set[E1], f func(E1) E2) Set[E2] {
	s2 := make(Set[E2])
	for e := range s {
		s2.Add(f(e))
	}
	return s2
}
