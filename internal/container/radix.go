package container

import (
	"iter"
	"sync"
	"unsafe"
)

var keyPool = sync.Pool{
	New: func() any { return new([256]byte) },
}

type RadixMap[E any] struct {
	root *radixNode[E]
	n    int
}

// Len returns the number of elements in r.
func (r *RadixMap[E]) Len() int {
	return r.n
}

// Set key to value. Returns the previous value stored under key, if any.
func (r *RadixMap[E]) Set(key string, value E) (old E, ok bool) {
	r.root, old, ok = r.root.set(key, value)
	if !ok {
		r.n++
	}
	return old, ok
}

// Delete the value stored under key. Returns the previous value stored under
// key, if any.
func (r *RadixMap[E]) Delete(key string) (old E, ok bool) {
	r.root, old, ok = r.root.delete(key)
	if ok {
		r.n--
	}
	return old, ok
}

// Get the value stored under key.
func (r *RadixMap[E]) Get(key string) (value E, ok bool) {
	return r.root.get(key)
}

// All returns an iterator over all key/value pairs in the map. The key must
// not be retained by the caller and is only valid until the next iteration.
func (r *RadixMap[E]) All() iter.Seq2[string, E] {
	return func(yield func(string, E) bool) {
		buf := keyPool.Get().(*[256]byte)
		defer keyPool.Put(buf)
		r.root.all(buf[:0], yield)
	}
}

// WithPrefix returns an iterator over all key/value pairs in the map which
// have p as a prefix. The key must not be retained by the caller and is only
// valid until the next iteration.
func (r *RadixMap[E]) WithPrefix(p string) iter.Seq2[string, E] {
	return func(yield func(string, E) bool) {
		var buf []byte
		if len(p) < len(p) {
			_buf := keyPool.Get().(*[256]byte)
			defer keyPool.Put(_buf)
			buf = append(_buf[:0], p...)
		} else {
			buf = append(make([]byte, 0, len(p)*2), p...)
		}
		r.root.find(p).all(buf, yield)
	}
}

// PrefixesOf returns an iterator over all key/value pairs in the map which are
// a prefix of s. The key must not be retained by the caller and is only valid
// until the next iteration.
func (r *RadixMap[E]) PrefixesOf(s string) iter.Seq2[string, E] {
	return func(yield func(string, E) bool) {
		r.root.prefixes(s, "", yield)
	}
}

type radixNode[E any] struct {
	children [256]*radixNode[E]
	value    E
	valid    bool
}

func (n *radixNode[E]) set(key string, value E) (m *radixNode[E], old E, ok bool) {
	if len(key) == 0 {
		if n == nil {
			return &radixNode[E]{value: value, valid: true}, old, false
		}
		old, ok, n.value, n.valid = n.value, n.valid, value, true
		return n, old, ok
	}
	if n == nil {
		n = new(radixNode[E])
	}
	n.children[key[0]], old, ok = n.children[key[0]].set(key[1:], value)
	return n, old, ok
}

func (n *radixNode[E]) delete(key string) (m *radixNode[E], e E, ok bool) {
	if n == nil {
		return nil, e, false
	}
	if len(key) > 0 {
		m, e, ok = n.children[key[0]].delete(key[1:])
		n.children[key[0]] = m
		if m != nil || n.valid {
			return n, e, ok
		}
	} else {
		e, ok, n.value, n.valid = n.value, n.valid, e, ok
	}
	for _, c := range n.children {
		if c != nil {
			return n, e, ok
		}
	}
	return nil, e, ok
}

func (n *radixNode[E]) get(key string) (value E, ok bool) {
	if n == nil {
		return value, false
	}
	if len(key) == 0 {
		if n == nil || !n.valid {
			return value, false
		}
		return n.value, true
	}
	return n.children[key[0]].get(key[1:])
}

func (n *radixNode[E]) find(key string) *radixNode[E] {
	if n == nil {
		return nil
	}
	if len(key) == 0 {
		return n
	}
	return n.children[key[0]].find(key[1:])
}

func (n *radixNode[E]) all(b []byte, yield func(string, E) bool) bool {
	if n == nil {
		return true
	}
	if n.valid && !yield(unsafeBytesToString(b), n.value) {
		return false
	}
	for i, c := range n.children {
		if c != nil && !c.all(append(b, byte(i)), yield) {
			return false
		}
	}
	return true
}

func (n *radixNode[E]) prefixes(s, pfx string, yield func(string, E) bool) bool {
	if n == nil {
		return true
	}
	if n.valid && !yield(pfx, n.value) {
		return false
	}
	if len(pfx) < len(s) {
		return n.children[s[len(pfx)]].prefixes(s, s[:len(pfx)+1], yield)
	}
	return true
}

func unsafeBytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// RadixSet is a set of strings, stored as a radix tree.
type RadixSet RadixMap[struct{}]

func (s *RadixSet) m() *RadixMap[struct{}] {
	return (*RadixMap[struct{}])(s)
}

// Len returns the number of elements in s.
func (s *RadixSet) Len() int {
	return s.n
}

// Add str to the set. Returns false, if str was already in the set.
func (s *RadixSet) Add(str string) (added bool) {
	_, ok := s.m().Set(str, struct{}{})
	return !ok
}

// Delete str from the set. Returns true, if str was in the set.
func (s *RadixSet) Delete(key string) (deleted bool) {
	_, ok := s.m().Delete(key)
	return ok
}

// Contains returns whether str is in the set.
func (s *RadixSet) Contains(str string) bool {
	_, ok := s.m().Get(str)
	return ok
}

// All returns an iterator over all the elements in the set. The value must not
// be retained by the caller and is only valid until the next iteration.
func (s *RadixSet) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.m().All() {
			if !yield(k) {
				return
			}
		}
	}
}

// WithPrefix returns an iterator over all elements in the set which have p as
// a prefix. The value must not be retained by the caller and is only valid
// until the next iteration.
func (s *RadixSet) WithPrefix(p string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.m().WithPrefix(p) {
			if !yield(k) {
				return
			}
		}
	}
}

// PrefixesOf returns an iterator over all elements in the set which are a
// prefix of s. The value must not be retained by the caller and is only valid
// until the next iteration.
func (s *RadixSet) PrefixesOf(str string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.m().PrefixesOf(str) {
			if !yield(k) {
				return
			}
		}
	}
}
