package main

import (
	"flag"
	"fmt"
)

var (
	P = 426
	M = 72058
)

func main() {
	flag.IntVar(&P, "p", 9, "Number of players")
	flag.IntVar(&M, "m", 25, "Highest marble")
	flag.Parse()

	scores := make([]int, P)
	c := newCircle()
	for i := 1; i <= M; i++ {
		if i%23 != 0 {
			c.cw(2)
			c.insert(i)
			continue
		}
		p := i % P
		c.ccw(7)
		scores[p] += c.delete()
		scores[p] += i
	}
	high := 0
	for _, s := range scores {
		if s > high {
			high = s
		}
	}
	fmt.Println("Highscore:", high)
}

type circle struct {
	m *element
}

type element struct {
	v    int
	next *element
	prev *element
}

func newCircle() *circle {
	e := new(element)
	e.next = e
	e.prev = e
	return &circle{m: e}
}

// cw moves i times clockwise around the circle
func (c *circle) cw(i int) {
	for ; i > 0; i-- {
		c.m = c.m.next
	}
}

// ccw moves i times counterclockwise around the circle
func (c *circle) ccw(i int) {
	for ; i > 0; i-- {
		c.m = c.m.prev
	}
}

// insert inserts v ccw of the current marble, making the inserted one the new
// current.
func (c *circle) insert(v int) {
	e := &element{
		v:    v,
		prev: c.m.prev,
		next: c.m,
	}
	c.m.prev.next = e
	c.m.prev = e
	c.m = e
}

// delete removes the current marble, making the cw-next one the new current.
func (c *circle) delete() int {
	e := c.m
	e.prev.next = e.next
	e.next.prev = e.prev
	c.m = e.next
	return e.v
}

func (c *circle) print(p int) {
	min := c.m
	z := c.m.next
	for z != c.m {
		if z.v < min.v {
			min = z
		}
		z = z.next
	}
	fmt.Printf("[%d]", p)
	z = min
	for {
		if z == c.m {
			fmt.Printf("(%2d)", z.v)
		} else {
			fmt.Printf(" %2d ", z.v)
		}
		z = z.next
		if z == min {
			break
		}
	}
	fmt.Println()
}
