package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/internal/container"
	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/math"
	"golang.org/x/exp/slices"
)

func main() {
	data, err := parse.Blocks(parse.Array[[2]Tree](split.Lines, ParseTree)).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var n int
	for i, p := range data {
		if Cmp(p[0], p[1]) <= 0 {
			n += i + 1
		}
	}
	fmt.Println("Index sum in right order:", n)

	// Use a pointer-slice, to make finding the dividers easier.
	all := Reduce(data)
	d1, d2 := Divider(2), Divider(6)
	all = append(all, d1, d2)
	slices.SortFunc(all, func(a, b *Tree) bool { return Cmp(*a, *b) < 0 })
	i1 := slices.Index(all, d1) + 1
	i2 := slices.Index(all, d2) + 1
	fmt.Println("Decoder key:", i1*i2)
}

type Tree struct {
	Val      int
	Children []Tree
}

func (t Tree) String() string {
	if t.Children == nil {
		return strconv.Itoa(t.Val)
	}
	var parts []string
	for _, c := range t.Children {
		parts = append(parts, c.String())
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func ParseTree(s string) (Tree, error) {
	var stack container.LIFO[Tree]
	stack.Push(Tree{})
	for len(s) > 0 {
		switch s[0] {
		case '[':
			t := Tree{Children: []Tree{}}
			stack.Push(t)
			s = s[1:]
		case ',':
			s = s[1:]
		case ']':
			t := stack.Pop()
			top := stack.Pop()
			top.Children = append(top.Children, t)
			stack.Push(top)
			s = s[1:]
		default:
			i := strings.IndexAny(s, "[],")
			if i < 0 {
				i = len(s)
			}
			n, err := strconv.Atoi(s[:i])
			if err != nil {
				return Tree{}, err
			}
			t := stack.Pop()
			t.Children = append(t.Children, Tree{Val: n})
			stack.Push(t)
			s = s[i:]
		}
	}
	top := stack.Pop()
	if len(top.Children) != 1 {
		return Tree{}, errors.New("unbalanced brackets")
	}
	return top.Children[0], nil
}

func Cmp(l, r Tree) (v int) {
	if l.Children == nil && r.Children == nil {
		return math.Cmp(l.Val, r.Val)
	}
	if l.Children != nil && r.Children != nil {
		i := 0
		for ; i < len(l.Children) && i < len(r.Children); i++ {
			switch Cmp(l.Children[i], r.Children[i]) {
			case -1:
				return -1
			case 1:
				return 1
			}
		}
		if i < len(l.Children) {
			return 1
		}
		if i < len(r.Children) {
			return -1
		}
		return 0
	}
	if l.Children != nil {
		return Cmp(l, Tree{Children: []Tree{{Val: r.Val}}})
	}
	return Cmp(Tree{Children: []Tree{{Val: l.Val}}}, r)
}

func Reduce(pairs [][2]Tree) []*Tree {
	var out []*Tree
	for _, p := range pairs {
		a, b := p[0], p[1]
		out = append(out, &a, &b)
	}
	return out
}

func Divider(i int) *Tree {
	return &Tree{Children: []Tree{
		{Children: []Tree{
			{Val: i},
		}},
	}}
}
