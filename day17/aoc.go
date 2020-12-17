package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func main() {
	grid3, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	grid4 := grid3.toGrid4()
	for i := 0; i < 6; i++ {
		grid3 = grid3.simulate()
	}
	fmt.Println("Number of alive cells:", len(grid3))
	for i := 0; i < 6; i++ {
		grid4 = grid4.simulate()
	}
	fmt.Println("Number of alive cells:", len(grid4))
}

func readInput(r io.Reader) (grid3, error) {
	s := bufio.NewScanner(r)
	y := 0
	grid3 := make(grid3)
	for s.Scan() {
		for x, b := range s.Bytes() {
			if b == '.' {
				continue
			}
			c := cell3{x, y, 0}
			grid3[c] = true
		}
		y++
	}
	return grid3, s.Err()
}

type grid3 map[cell3]bool

type cell3 [3]int

func (g grid3) toGrid4() grid4 {
	g4 := make(grid4)
	for c := range g {
		g4[cell4{c[0], c[1], c[2], 0}] = true
	}
	return g4
}

func (g grid3) simulate() grid3 {
	next := make(grid3)
	for c := range g {
		n := g.countAliveNeighbors(c)
		if n == 2 || n == 3 {
			next[c] = true
		}
		g.walkNeighbors(c, func(c cell3) {
			if g[c] {
				return
			}
			if n := g.countAliveNeighbors(c); n == 3 {
				next[c] = true
			}
		})
	}
	return next
}

func (g grid3) walkNeighbors(c cell3, f func(cell3)) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := -1; k <= 1; k++ {
				if i == 0 && j == 0 && k == 0 {
					continue
				}
				c2 := cell3{
					c[0] + i,
					c[1] + j,
					c[2] + k,
				}
				f(c2)
			}
		}
	}
}

func (g grid3) countAliveNeighbors(c cell3) int {
	var n int
	g.walkNeighbors(c, func(c cell3) {
		if g[c] {
			n++
		}
	})
	return n
}

func (g grid3) bounds() (min, max cell3) {
	min = cell3{math.MaxInt64, math.MaxInt64, math.MaxInt64}
	max = cell3{math.MinInt64, math.MinInt64, math.MinInt64}
	for c := range g {
		if c[0] < min[0] {
			min[0] = c[0]
		}
		if c[0] > max[0] {
			max[0] = c[0]
		}
		if c[1] < min[1] {
			min[1] = c[1]
		}
		if c[1] > max[1] {
			max[1] = c[1]
		}
		if c[2] < min[2] {
			min[2] = c[2]
		}
		if c[2] > max[2] {
			max[2] = c[2]
		}
	}
	return min, max
}

func (g grid3) dump() {
	min, max := g.bounds()
	for z := min[2]; z <= max[2]; z++ {
		fmt.Printf("z=%d\n", z)
		for y := min[1]; y <= max[1]; y++ {
			for x := min[0]; x <= max[0]; x++ {
				if g[cell3{x, y, z}] {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

type grid4 map[cell4]bool

type cell4 [4]int

func (g grid4) simulate() grid4 {
	next := make(grid4)
	for c := range g {
		n := g.countAliveNeighbors(c)
		if n == 2 || n == 3 {
			next[c] = true
		}
		g.walkNeighbors(c, func(c cell4) {
			if g[c] {
				return
			}
			if n := g.countAliveNeighbors(c); n == 3 {
				next[c] = true
			}
		})
	}
	return next
}

func (g grid4) walkNeighbors(c cell4, f func(cell4)) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := -1; k <= 1; k++ {
				for l := -1; l <= 1; l++ {
					if i == 0 && j == 0 && k == 0 && l == 0 {
						continue
					}
					c2 := cell4{
						c[0] + i,
						c[1] + j,
						c[2] + k,
						c[3] + l,
					}
					f(c2)
				}
			}
		}
	}
}

func (g grid4) countAliveNeighbors(c cell4) int {
	var n int
	g.walkNeighbors(c, func(c cell4) {
		if g[c] {
			n++
		}
	})
	return n
}

func (g grid4) bounds() (min, max cell4) {
	min = cell4{math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}
	max = cell4{math.MinInt64, math.MinInt64, math.MinInt64, math.MinInt64}
	for c := range g {
		for i := 0; i < 4; i++ {
			if c[i] < min[i] {
				min[i] = c[i]
			}
			if c[i] > max[i] {
				max[i] = c[i]
			}
		}
	}
	return min, max
}

func (g grid4) dump() {
	min, max := g.bounds()
	fmt.Println(min, max)
	for w := min[3]; w <= max[3]; w++ {
		for z := min[2]; z <= max[2]; z++ {
			fmt.Printf("z=%d, w=%d\n", z, w)
			for y := min[1]; y <= max[1]; y++ {
				for x := min[0]; x <= max[0]; x++ {
					if g[cell4{x, y, z, w}] {
						fmt.Print("#")
					} else {
						fmt.Print(".")
					}
				}
				fmt.Println()
			}
			fmt.Println()
		}
	}
}
