package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	log.SetFlags(log.Lshortfile)

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	data = bytes.TrimSpace(data)

	fmt.Printf("After 2022 rocks, the tower is %d high\n", Simulate(data, 2022))
	fmt.Printf("After 1000000000000 rocks, the tower is %d high\n", Simulate(data, 1000000000000))
}

const frameDelay = time.Second

type Cell int8

const (
	Empty Cell = iota
	Rock
	Shadow
)

func (c Cell) String() string {
	switch c {
	case Empty:
		return " "
	case Rock:
		return "▇"
	case Shadow:
		return "░"
	default:
		panic("invalid cell")
	}
}

type Shape [][]Cell

func (s Shape) String() string {
	var parts []string
	for _, r := range s {
		var l string
		for _, c := range r {
			l += c.String()
		}
		parts = append(parts, l)
	}
	return strings.Join(parts, "\n")
}

// Shapes are flipped.
var Shapes = []Shape{
	Shape{{Rock, Rock, Rock, Rock}},
	Shape{{Empty, Rock, Empty}, {Rock, Rock, Rock}, {Empty, Rock, Empty}},
	Shape{{Rock, Rock, Rock}, {Empty, Empty, Rock}, {Empty, Empty, Rock}},
	Shape{{Rock}, {Rock}, {Rock}, {Rock}},
	Shape{{Rock, Rock}, {Rock, Rock}},
}

func Simulate(data []byte, n int) int {
	s := State{
		prog: string(data),
	}

	type loop struct {
		i     int
		total int
	}
	memo := make(map[State]loop)
	var total int

	for i := 0; i < n; i++ {
		if l, ok := memo[s]; ok {
			δi := i - l.i
			δt := total - l.total
			for ; i+δi < n; i += δi {
				total += δt
			}
		} else {
			memo[s] = loop{i, total}
		}
		next, dropped := s.DropRock()
		s, total = next, total+dropped
	}
	return s.Height() + total
}

const (
	Cols = 7
	Rows = 1 << 12
)

type State struct {
	prog  string
	pc    int
	shape int
	grid  [Rows * Cols]Cell
}

func (s *State) Height() int {
yloop:
	for y := 0; y < Rows; y++ {
		for x := 0; x < Cols; x++ {
			if s.at(x, y) != Empty {
				continue yloop
			}
		}
		return y
	}
	panic("no empty rows")
}

func (s State) DropRock() (next State, dropped int) {
	x, y := 2, s.Height()+3
	shape := Shapes[s.shape]
	s.shape = (s.shape + 1) % len(Shapes)
	for {
		nx := x
		if s.read() == '<' {
			nx = math.Max(x-1, 0)
		} else {
			nx = math.Min(x+1, Cols-len(shape[0]))
		}
		if s.valid(shape, nx, y) {
			x = nx
		}

		if y == 0 || !s.valid(shape, x, y-1) {
			s.apply(shape, x, y, Rock)
			break
		}
		y--
	}
	n := s.dropRows()
	return s, n
}

func (s *State) read() byte {
	c := s.prog[s.pc]
	s.pc++
	if s.pc == len(s.prog) {
		s.pc = 0
	}
	return c
}

func (s *State) at(x, y int) Cell {
	if y*Cols+x >= len(s.grid) {
		panic(fmt.Errorf("(%d,%d) out of bounds", x, y))
	}
	return s.grid[y*Cols+x]
}

func (s *State) set(x, y int, c Cell) {
	s.grid[y*Cols+x] = c
}

func (s *State) valid(sh Shape, x, y int) bool {
	for δy := 0; δy < len(sh); δy++ {
		for δx := 0; δx < len(sh[0]); δx++ {
			if sh[δy][δx] == Empty {
				continue
			}
			if s.at(x+δx, y+δy) != Empty {
				return false
			}
		}
	}
	return true
}

func (s *State) apply(sh Shape, x, y int, c Cell) {
	for δy := 0; δy < len(sh); δy++ {
		for δx := 0; δx < len(sh[0]); δx++ {
			if sh[δy][δx] == Empty {
				continue
			}
			s.set(x+δx, y+δy, c)
		}
	}
}

func (s *State) dropRows() int {
rowLoop:
	for y := 0; y < Rows; y++ {
		for x := 0; x < Cols; x++ {
			if s.at(x, y) == Empty {
				continue rowLoop
			}
		}
		// y is a full row. Drop it and everything below
		m := (y + 1) * Cols
		copy(s.grid[:], s.grid[m:])
		for i := len(s.grid) - m; i < len(s.grid); i++ {
			s.grid[i] = Empty
		}
		return y + 1
	}
	return 0
}

func (s State) String() string {
	w := new(strings.Builder)
	fmt.Fprintf(w, "pc=%d, shape=%d\n", s.pc, s.shape)
	var maxy int
maxyLoop:
	for maxy = Rows - 1; maxy >= 0; maxy-- {
		for x := 0; x < Cols; x++ {
			if s.at(x, maxy) != Empty {
				maxy++
				break maxyLoop
			}
		}
	}
	w.WriteString("┌" + strings.Repeat("─", Cols) + "┐\n")
	for y := maxy; y >= 0; y-- {
		w.WriteString("│")
		for x := 0; x < Cols; x++ {
			w.WriteString(s.at(x, y).String())
		}
		w.WriteString("│\n")
	}
	w.WriteString("└" + strings.Repeat("─", Cols) + "┘")
	return w.String()
}
