package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"sort"
)

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Require argument")
	}
	cups, err := parseArg(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	g := Game{
		cur:  0,
		cups: cups,
	}
	fmt.Println(g)
	for i := 0; i < 100; i++ {
		g = doMove(g)
	}
	fmt.Println(g)
}

func parseArg(s string) ([]int, error) {
	cups := make([]int, len(s))
	for i, r := range s {
		if '0' > r || r > '9' {
			return nil, errors.New("argument must be only digits")
		}
		cups[i] = int(r - '0')
	}
	if len(cups) > 10 {
		return nil, errors.New("argument contains duplicate digits")
	}
	cups2 := append([]int(nil), cups...)
	sort.Ints(cups2)
	for i, v := range cups2 {
		if i+1 != v {
			return nil, errors.New("argument does not contain all digits 1-9")
		}
	}
	return ring(cups), nil
}

type ring []int

func (r ring) at(i int) int {
	return r[i%len(r)]
}

func (r ring) slice(i, j int) ring {
	i %= len(r)
	j %= len(r)
	if j >= i {
		return r[i:j]
	}
	return append(append(ring(nil), r[i:]...), r[:j]...)
}

func (r ring) has(needle int) bool {
	for _, v := range r {
		if v == needle {
			return true
		}
	}
	return false
}

func (r ring) find(needle int) int {
	for i, v := range r {
		if v == needle {
			return i
		}
	}
	panic("needle not in haystack")
}

// shift rotates the ring around to make n the first index
func (r ring) shift(n int) ring {
	n = ((n % len(r)) + len(r)) % len(r)
	return append(append(ring(nil), r[n:]...), r[:n]...)
}

type Game struct {
	cur  int  // index of the current cup
	cups ring // cups
}

func doMove(g Game) Game {
	c := g.cups[g.cur]
	pu, cups := g.cups.slice(g.cur+1, g.cur+4), g.cups.slice(g.cur+4, g.cur+1)
	dest := c - 1
	if dest < 1 {
		dest = 9
	}
	for pu.has(dest) {
		dest--
		if dest < 1 {
			dest = 9
		}
	}
	d := search(cups, dest)
	var out ring
	out = append(out, dest)
	out = append(out, pu...)
	out = append(out, cups.slice(d+1, d)...)
	return Game{
		cur:  (out.find(c) + 1) % len(out),
		cups: out,
	}
}

func in(a []int, needle int) bool {
	for _, v := range a {
		if v == needle {
			return true
		}
	}
	return false
}

func search(a []int, needle int) int {
	for i, v := range a {
		if v == needle {
			return i
		}
	}
	panic("needle not in slice")
}
