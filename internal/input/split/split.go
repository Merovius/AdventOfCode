package split

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Func func(string) ([]string, error)

// Split into blocks, separated by empty lines.
func Blocks(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errors.New("empty input")
	}
	return strings.Split(s, "\n\n"), nil
}

// Split into whitespace-separated fields.
func Fields(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errors.New("empty input")
	}
	return strings.Fields(s), nil
}

// Split into lines.
func Lines(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errors.New("empty input")
	}
	return strings.Split(s, "\n"), nil
}

// Split along sep.
func On(sep string) Func {
	return func(s string) ([]string, error) {
		return strings.Split(s, sep), nil
	}
}

// Split along any of a set of separators. If multiple separators overlap, the
// first match is used.
func Any(sep ...string) Func {
	return func(in string) ([]string, error) {
		var out []string
		for len(in) > 0 {
			i, n := len(in), 0
			for _, s := range sep {
				if j := strings.Index(in, s); j >= 0 && j < i {
					i, n = j, len(s)
				}
			}
			out, in = append(out, in[:i]), in[i+n:]
		}
		return out, nil
	}
}

// Split into at most n pieces, along sep.
func SplitN(sep string, n int) Func {
	return func(s string) ([]string, error) {
		return strings.SplitN(s, sep, n), nil
	}
}

// Split into capture groups of a regular expression. The regular expression
// must match the full parsed string.
func Regexp(re string) Func {
	if !strings.HasPrefix(re, "^") {
		re = "^" + re
	}
	if !strings.HasSuffix(re, "$") {
		re += "$"
	}
	r, err := regexp.Compile(re)
	if err != nil {
		panic(fmt.Errorf("regexp.Compile(%q) = %v", re, err))
	}
	return func(s string) ([]string, error) {
		sp := r.FindStringSubmatch(s)
		if sp == nil || len(sp[0]) != len(s) {
			return nil, fmt.Errorf("%q does not match %q", s, re)
		}
		return sp[1:], nil
	}
}
