package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type Graph map[string][]string

func main() {
	g, err := Parse("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Paths for part 1: %v\n", len(FindPaths(g, abort1)))
	fmt.Printf("Paths for part 2: %v\n", len(FindPaths(g, abort2)))

}

func Parse(name string) (Graph, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	g := make(Graph)
	for s.Scan() {
		sp := strings.Split(s.Text(), "-")
		if len(sp) != 2 {
			return nil, fmt.Errorf("invalid graph element %q", s.Text())
		}
		g[sp[0]] = append(g[sp[0]], sp[1])
		g[sp[1]] = append(g[sp[1]], sp[0])
	}
	return g, s.Err()
}

type State struct {
	path    []string
	visited map[string]int
}

func NewState() *State {
	return &State{visited: make(map[string]int)}
}

func isLower(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !unicode.IsLower(r)
	}) < 0
}

func (s *State) Push(n string) {
	s.path = append(s.path, n)
	if isLower(n) {
		s.visited[n] += 1
	}
}

func (s *State) Pop() {
	var n string
	s.path, n = s.path[:len(s.path)-1], s.path[len(s.path)-1]
	if isLower(n) {
		s.visited[n] -= 1
	}
}

func FindPaths(g Graph, abort func(*State) bool) [][]string {
	s := NewState()
	var paths [][]string

	var dfs func(n string)
	dfs = func(n string) {
		s.Push(n)
		defer s.Pop()

		if n == "end" {
			paths = append(paths, append([]string(nil), s.path...))
			return
		}
		if abort(s) {
			return
		}
		for _, m := range g[n] {
			dfs(m)
		}
	}
	dfs("start")
	return paths
}

func abort1(s *State) bool {
	for _, v := range s.visited {
		if v > 1 {
			return true
		}
	}
	return false
}

func abort2(s *State) bool {
	if s.visited["start"] > 1 {
		return true
	}
	revisited := false
	for _, v := range s.visited {
		if v > 2 {
			return true
		}
		if v > 1 {
			if revisited {
				return true
			}
			revisited = true
		}
	}
	return false
}
