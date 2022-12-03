package set

type Set[E comparable] map[E]struct{}

func (s Set[E]) Add(e E) {
	s[e] = struct{}{}
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
