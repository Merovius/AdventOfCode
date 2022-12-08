package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	rules, err := ParseRules(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d colors may contain shiny gold bags\n", CountValidContainers(rules, "shiny gold"))
	fmt.Printf("%d bags are in a shiny gold bag\n", CountBags(rules, "shiny gold"))
}

type Rule struct {
	Color    string
	Contains map[string]int
}

func ParseRules(r io.Reader) ([]Rule, error) {
	var rules []Rule
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		i := strings.Index(t, " bags contain ")
		if i <= 0 {
			return nil, fmt.Errorf("can not parse rule %q", t)
		}
		var r Rule
		r.Color = t[:i]
		r.Contains = make(map[string]int)
		t = t[i+len(" bags contain "):]
		t = strings.TrimSuffix(t, ".")
		if t == "no other bags" {
			rules = append(rules, r)
			continue
		}
		sp := strings.Split(t, ", ")
		for _, s := range sp {
			s = strings.TrimSuffix(s, " bags")
			s = strings.TrimSuffix(s, " bag")
			var (
				count int
				color string
			)

			i := strings.IndexByte(s, ' ')
			if i < 0 {
				return nil, fmt.Errorf("can not parse rule %q", t)
			}
			var err error
			if count, err = strconv.Atoi(s[:i]); err != nil {
				return nil, fmt.Errorf("can not parse rule %q", t)
			}
			color = s[i+1:]

			if _, ok := r.Contains[color]; ok {
				return nil, fmt.Errorf("duplicate color %q in rule %q", color, t)
			}
			r.Contains[color] = count
		}
		rules = append(rules, r)
	}
	return rules, s.Err()
}

func CountValidContainers(rs []Rule, color string) int {
	// Build a graph, where a node is a color and there is an edge from A to B,
	// if A can be contained in B.
	type node struct {
		color   string
		mayBeIn []string
		visited int
	}
	nodes := make(map[string]*node)
	find := func(c string) *node {
		if n := nodes[c]; n != nil {
			return n
		}
		n := &node{color: c}
		nodes[c] = n
		return n
	}

	for _, r := range rs {
		for color, count := range r.Contains {
			if count == 0 {
				continue
			}
			n := find(color)
			n.mayBeIn = append(n.mayBeIn, r.Color)
		}
	}

	// Do a DFS of the graph.
	count := 0
	q := []*node{find(color)}
	push := func(n *node) {
		q = append(q, n)
	}
	pop := func() (n *node) {
		q, n = q[:len(q)-1], q[len(q)-1]
		return n
	}
	for len(q) > 0 {
		n := pop()
		n.visited++
		if n.visited > 1 {
			continue
		}
		count++
		for _, c := range n.mayBeIn {
			push(find(c))
		}
	}
	// If we visited color more than once, it can transitively contain
	// itself. Otherwise, we included it erronously in the count.
	if find(color).visited == 1 {
		count--
	}
	return count
}

func CountBags(rs []Rule, color string) int {
	contains := make(map[string]map[string]int)
	for _, r := range rs {
		if contains[r.Color] != nil {
			panic("duplicate color in rules")
		}
		contains[r.Color] = r.Contains
	}
	var count func(c string) int
	count = func(c string) int {
		N := 1
		for C, n := range contains[c] {
			N += n * count(C)
		}
		return N
	}
	// Don't count the outermost bag itself
	return count(color) - 1
}
