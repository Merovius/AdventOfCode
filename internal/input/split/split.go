package split

import (
	"errors"
	"fmt"
	"iter"
	"regexp"
	"strings"
)

type Func func(string) iter.Seq2[string, error]

var (
	_ Func = Blocks
	_ Func = Fields
	_ Func = Lines
	_ Func = Bytes
)

// Split into blocks, separated by empty lines.
func Blocks(s string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		s = strings.Trim(s, "\n")
		if len(s) == 0 {
			yield("", errors.New("split.Blocks: empty input"))
			return
		}
		for s := range strings.SplitSeq(s, "\n\n") {
			if !yield(s, nil) {
				return
			}
		}
	}
}

// Split into whitespace-separated fields.
func Fields(s string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		if len(s) == 0 {
			yield("", errors.New("split.Fields: empty input"))
			return
		}
		for s := range strings.FieldsSeq(s) {
			if !yield(s, nil) {
				return
			}
		}
	}
}

// Split into lines.
func Lines(s string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		s = strings.Trim(s, "\n")
		if len(s) == 0 {
			yield("", errors.New("split.Lines: empty input"))
			return
		}
		for s := range strings.SplitSeq(s, "\n") {
			if !yield(s, nil) {
				return
			}
		}
	}
}

// Split along sep.
func On(sep string) Func {
	return func(s string) iter.Seq2[string, error] {
		return func(yield func(string, error) bool) {
			for s := range strings.SplitSeq(s, sep) {
				if !yield(s, nil) {
					return
				}
			}
		}
	}
}

// Split along any of a set of separators. If multiple separators overlap, the
// first match is used.
func Any(sep ...string) Func {
	return func(in string) iter.Seq2[string, error] {
		return func(yield func(string, error) bool) {
			for len(in) > 0 {
				for _, s := range sep {
					if before, after, ok := strings.Cut(in, s); ok {
						if !yield(before, nil) {
							return
						}
						in = after
						break
					}
				}
				yield(in, nil)
				return
			}
		}
	}
}

// Split into at most n pieces, along sep.
func SplitN(sep string, n int) Func {
	if n < 0 {
		panic("n must not be negative")
	}
	return func(s string) iter.Seq2[string, error] {
		return func(yield func(string, error) bool) {
			i := 1
			for len(s) > 0 {
				if i == n {
					yield(s, nil)
					return
				}
				before, after, ok := strings.Cut(s, sep)
				if !ok {
					break
				}
				if !yield(before, nil) {
					return
				}
				s = after
				i++
			}
		}
	}
}

// Split into capture groups of a regular expression. The regular expression
// must match the full parsed string.
func Regexp(re string) Func {
	if !strings.HasPrefix(re, `\A`) {
		re = `\A` + re
	}
	if !strings.HasSuffix(re, `\z`) {
		re += `\z`
	}
	r, err := regexp.Compile(re)
	if err != nil {
		panic(fmt.Errorf("regexp.Compile(%q) = %v", re, err))
	}
	return func(s string) iter.Seq2[string, error] {
		return func(yield func(string, error) bool) {
			sp := r.FindStringSubmatch(s)
			if sp == nil || len(sp[0]) != len(s) {
				yield("", fmt.Errorf("%q does not match %q", s, re))
				return
			}
			for _, s := range sp[1:] {
				if !yield(s, nil) {
					return
				}
			}
		}
	}
}

// Split into bytes.
func Bytes(s string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for i := range len(s) {
			if !yield(s[i:i+1], nil) {
				return
			}
		}
	}
}

// Split into pieces of fixed byte sizes. The last piece will be the rest of the string.
func After(ns ...int) Func {
	return func(s string) iter.Seq2[string, error] {
		return func(yield func(string, error) bool) {
			for _, n := range ns {
				if len(s) < n {
					yield("", fmt.Errorf("%q less than %d bytes", s, n))
					return
				}
				if !yield(s[:n], nil) {
					return
				}
				s = s[n:]
			}
			yield(s, nil)
		}
	}
}
