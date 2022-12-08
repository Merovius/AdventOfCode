package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	rules, words, err := ParseInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if err := ResolveRefs(rules); err != nil {
		log.Fatal(err)
	}
	var n int
	for _, w := range words {
		if representsMatch(rules[0].Matches(w)) {
			fmt.Println(w)
			n++
		}
	}
	fmt.Printf("%d words match the rules\n", n)
}

func representsMatch(rests []string) bool {
	for _, r := range rests {
		if len(r) == 0 {
			return true
		}
	}
	return false
}

func ParseInput(r io.Reader) ([]Rule, []string, error) {
	var rules []Rule
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			break
		}
		r, idx, err := ParseRule(l)
		if err != nil {
			return nil, nil, err
		}
		if len(rules) <= idx {
			tmp := make([]Rule, idx+1)
			copy(tmp, rules)
			rules = tmp
		}
		rules[idx] = r
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	var words []string
	for s.Scan() {
		words = append(words, strings.TrimSpace(s.Text()))
	}
	return rules, words, nil
}

func ResolveRefs(rs []Rule) error {
	for _, r := range rs {
		if r == nil {
			continue
		}
		if err := r.ResolveRefs(rs); err != nil {
			return err
		}
	}
	return nil
}

type Rule interface {
	Matches(s string) (rest []string)
	ResolveRefs(rs []Rule) error
	String() string
}

type literal string

func (r literal) Matches(s string) []string {
	if strings.HasPrefix(s, string(r)) {
		return []string{s[len(r):]}
	}
	return nil
}

func (r literal) ResolveRefs(rs []Rule) error {
	return nil
}

func (r literal) String() string {
	return strconv.Quote(string(r))
}

type concatenation []Rule

func (r concatenation) Matches(s string) []string {
	rests := []string{s}
	for _, r := range r {
		var next []string
		for _, s := range rests {
			m := r.Matches(s)
			next = append(next, m...)
		}
		rests = next
	}
	return rests
}

func (r concatenation) ResolveRefs(rs []Rule) error {
	for _, x := range r {
		if err := x.ResolveRefs(rs); err != nil {
			return err
		}
	}
	return nil
}

func (r concatenation) String() string {
	var parts []string
	for _, s := range r {
		parts = append(parts, s.String())
	}
	return strings.Join(parts, " . ")
}

type alternation []Rule

func (r alternation) Matches(s string) []string {
	var rests []string
	for _, r := range r {
		rests = append(rests, r.Matches(s)...)
	}
	return rests
}

func (r alternation) ResolveRefs(rs []Rule) error {
	for _, r := range r {
		if err := r.ResolveRefs(rs); err != nil {
			return err
		}
	}
	return nil
}

func (r alternation) String() string {
	var parts []string
	for _, s := range r {
		parts = append(parts, s.String())
	}
	return "(" + strings.Join(parts, " | ") + ")"
}

type ref struct {
	idx int
	r   Rule // resolved subrule
}

func (r *ref) Matches(s string) []string {
	return r.r.Matches(s)
}

func (r *ref) ResolveRefs(rs []Rule) error {
	r.r = rs[r.idx]
	if r.r == nil {
		return fmt.Errorf("unresolved reference %d", r.idx)
	}
	return nil
}

func (r *ref) String() string {
	return strconv.Itoa(int(r.idx))
}

func ParseRule(l string) (r Rule, idx int, err error) {
	i := strings.IndexByte(l, ':')
	if i < 0 {
		return nil, 0, errors.New("missing : in rule definition")
	}
	idx, err = strconv.Atoi(l[:i])
	if err != nil {
		return nil, 0, fmt.Errorf("can not parse rule index: %w", err)
	}

	l = strings.TrimSpace(l[i+1:])
	if len(l) == 0 {
		return nil, 0, errors.New("empty rule")
	}
	if l[0] == '"' {
		if l[len(l)-1] != '"' {
			return nil, 0, errors.New("unterminated string literal")
		}
		return literal(l[1 : len(l)-1]), idx, nil
	}
	var a alternation
	sp := strings.Split(l, "|")
	for _, s := range sp {
		s = strings.TrimSpace(s)
		sp := strings.Split(s, " ")
		var r concatenation
		for _, s := range sp {
			s = strings.TrimSpace(s)
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, 0, fmt.Errorf("could not parse rule ref %q: %w", s, err)
			}
			r = append(r, &ref{idx: v})
		}
		if len(r) == 1 {
			a = append(a, r[0])
		} else {
			a = append(a, r)
		}
	}
	if len(a) == 1 {
		return a[0], idx, nil
	}
	return a, idx, nil
}
