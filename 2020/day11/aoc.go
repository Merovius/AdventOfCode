package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	g, err := readGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g2 := g.copy()
	for g.stepNeighbor() {
	}
	fmt.Printf("%d occupied seats when considering neighbors\n", g.occupiedSeats())
	for g2.stepRay() {
	}
	fmt.Printf("%d occupied seats when considering rays\n", g2.occupiedSeats())
}

type grid []string

func readGrid(r io.Reader) (grid, error) {
	var out grid
	s := bufio.NewScanner(r)
	for s.Scan() {
		out = append(out, s.Text())
		if len(out) > 1 && (len(out[len(out)-1]) != len(out[0])) {
			return nil, errors.New("lines have different lengths")
		}
	}
	return out, s.Err()
}

func (g grid) copy() grid {
	return append(grid(nil), g...)
}

func (g grid) stepNeighbor() (changed bool) {
	next := make(grid, len(g))
	for i, l := range g {
		var w strings.Builder
		for j, b := range []byte(l) {
			if b == '.' {
				w.WriteByte(b)
				continue
			}
			var n int
			n += g.occupied(i-1, j-1)
			n += g.occupied(i-1, j)
			n += g.occupied(i-1, j+1)
			n += g.occupied(i, j-1)
			n += g.occupied(i, j+1)
			n += g.occupied(i+1, j-1)
			n += g.occupied(i+1, j)
			n += g.occupied(i+1, j+1)
			if b == 'L' && n == 0 {
				w.WriteByte('#')
			} else if b == '#' && n >= 4 {
				w.WriteByte('L')
			} else {
				w.WriteByte(b)
			}
		}
		nl := w.String()
		changed = changed || (nl != l)
		next[i] = nl
	}
	copy(g, next)
	return changed
}

func (g grid) stepRay() (changed bool) {
	next := make(grid, len(g))
	for i, l := range g {
		var w strings.Builder
		for j, b := range []byte(l) {
			if b == '.' {
				w.WriteByte(b)
				continue
			}
			var n int
			n += g.seesOccupied(i, j, -1, -1)
			n += g.seesOccupied(i, j, -1, 0)
			n += g.seesOccupied(i, j, -1, 1)
			n += g.seesOccupied(i, j, 0, -1)
			n += g.seesOccupied(i, j, 0, 1)
			n += g.seesOccupied(i, j, 1, -1)
			n += g.seesOccupied(i, j, 1, 0)
			n += g.seesOccupied(i, j, 1, 1)
			if b == 'L' && n == 0 {
				w.WriteByte('#')
			} else if b == '#' && n >= 5 {
				w.WriteByte('L')
			} else {
				w.WriteByte(b)
			}
		}
		nl := w.String()
		changed = changed || (nl != l)
		next[i] = nl
	}
	copy(g, next)
	return changed
}

// occupied returns 1 if the given seat is occupied and 0 else. Positions
// outside the grid are treated as unoccupied.
func (g grid) occupied(y, x int) int {
	if x < 0 || y < 0 || x >= len(g[0]) || y >= len(g) {
		return 0
	}
	if g[y][x] == '#' {
		return 1
	}
	return 0
}

// seesOccupied returns 1 if the next seat seen in direction di,dj from seat
// i,j is occupied and 0 else.
func (g grid) seesOccupied(i, j, di, dj int) int {
	for {
		i += di
		j += dj
		if i < 0 || j < 0 || i >= len(g) || j >= len(g[0]) {
			return 0
		}
		switch g[i][j] {
		case '.':
		case 'L':
			return 0
		case '#':
			return 1
		}
	}
}

func (g grid) dump() {
	for _, l := range g {
		fmt.Println(l)
	}
	fmt.Println()
}

func (g grid) occupiedSeats() int {
	var n int
	for _, l := range g {
		for _, c := range []byte(l) {
			if c == '#' {
				n++
			}
		}
	}
	return n
}
