package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

func main() {
	print := flag.Bool("print", false, "print grid state")
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	st := readState()
	if *print {
		fmt.Print("\x1b[2J\x1b[H")
		st.print()
		time.Sleep(time.Second)
	}

	for i := 0; ; i++ {
		x, y, done := st.tick()
		if *print {
			fmt.Print("\x1b[H")
			st.print()
			time.Sleep(time.Second)
		}
		if done {
			fmt.Printf("%d,%d\n", x, y)
			return
		}
	}
}

type state struct {
	w     int
	h     int
	grid  []cell
	carts []int
}

func (st *state) tick() (x, y int, done bool) {
	if len(st.carts) == 0 {
		log.Fatal("no carts")
	}
	for i := 0; i < len(st.carts); i++ {
		idx := st.carts[i]
		d, cy, ok := st.grid[idx].cart()
		if !ok {
			log.Fatalf("No cart at %d", i)
		}
		st.grid[idx] &= rMask
		switch d {
		case dUp:
			idx -= st.w
		case dDown:
			idx += st.w
		case dLeft:
			if idx%st.w == 0 {
				log.Fatal("Cart crashed into wall")
			}
			idx -= 1
		case dRight:
			idx += 1
			if idx%st.w == 0 {
				log.Fatal("Cart crashed into wall")
			}
		}
		if idx < 0 || idx >= len(st.grid) {
			log.Fatal("Cart crashed into wall")
		}
		st.carts[i] = idx
		if _, _, ok := st.grid[idx].cart(); ok {
			// remove both carts
			st.carts = append(st.carts[:i], st.carts[i+1:]...)
			for j, jdx := range st.carts {
				if jdx == idx {
					if j < i {
						i--
					}
					st.carts = append(st.carts[:j], st.carts[j+1:]...)
					break
				}
			}
			st.grid[idx] &= rMask
			i--
			continue
		}
		switch st.grid[idx].rail() {
		case rCurveR:
			// switches left<->down and right<->up
			d ^= 0x10
		case rCurveL:
			// switches left<->up and right<->down
			d ^= 0x30
		case rCross:
			switch cy {
			case cRight:
				d = (d + 0x10) & dMask
			case cLeft:
				d = (d - 0x10) & dMask
			}
			cy += 0x40
			if cy == 0xC0 {
				cy = 0x00
			}
		}
		st.grid[idx] |= 0x08 | cell(d) | cell(cy)
	}
	sort.Ints(st.carts)
	if len(st.carts) == 1 {
		return st.carts[0] % st.w, st.carts[0] / st.w, true
	}
	return 0, 0, false
}

func red(s string) string {
	return "\x1b[1;31m" + s + "\x1b[m"
}

func green(s string) string {
	return "\x1b[1;32m" + s + "\x1b[m"
}

func (st *state) print() {
	for i, c := range st.grid {
		if i != 0 && (i%st.w) == 0 {
			fmt.Println()
		}
		if c == empty {
			fmt.Print(" ")
			continue
		}
		if c == crash {
			fmt.Print(red("X"))
			continue
		}
		_, _, ok := c.cart()
		if ok {
			c &= dMask
		}
		switch c &^ cMask {
		case cell(rVert):
			fmt.Print("|")
		case cell(rHoriz):
			fmt.Print("-")
		case cell(rCurveR):
			fmt.Print("/")
		case cell(rCurveL):
			fmt.Print("\\")
		case cell(rCross):
			fmt.Print("+")
		case cell(dUp):
			fmt.Print(green("^"))
		case cell(dRight):
			fmt.Print(green(">"))
		case cell(dDown):
			fmt.Print(green("v"))
		case cell(dLeft):
			fmt.Print(green("<"))
		default:
			log.Fatalf("invalid cell: %#x", c&^cMask)
		}
	}
	fmt.Println()
}

type cell uint8

const (
	empty cell = cell(cMask) + iota
	crash
)

type rail uint8

const (
	rVert rail = iota
	rHoriz
	rCurveR
	rCurveL
	rCross

	rMask = 0x07
)

type direction uint8

const (
	dUp direction = 0x08 | (0x10 * iota)
	dRight
	dDown
	dLeft

	dMask = 0x38
)

type cycle uint8

const (
	cLeft     cycle = 0x00
	cStraight cycle = 0x40
	cRight    cycle = 0x80

	cMask = 0xC0
)

func (c cell) rail() rail {
	return rail(c & 0x7)
}

func (c cell) cart() (direction, cycle, bool) {
	if c == empty {
		return 0, 0, false
	}
	if c == crash {
		return 0, 0, true
	}
	return direction(c) & dMask, cycle(c) & cMask, (c & 0x08) != 0
}

// railtype  = |-/\+ -> 3 bit
// hasCart 		     -> 1 bit
// cartdir   = ^v<>  -> 2 bit
// cartcycle = <>|   -> 2 bit

var rune2cell = map[rune]cell{
	' ':  empty,
	'|':  cell(rVert),
	'-':  cell(rHoriz),
	'/':  cell(rCurveR),
	'\\': cell(rCurveL),
	'+':  cell(rCross),
	'v':  cell(rVert) | cell(dDown),
	'^':  cell(rVert) | cell(dUp),
	'>':  cell(rHoriz) | cell(dRight),
	'<':  cell(rHoriz) | cell(dLeft),
}

func readState() state {
	var st state
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if st.w != 0 && len(l) != st.w {
			log.Fatal("lines have differing lengths:", st.w, len(l))
		}
		st.w = len(l)
		for i, r := range l {
			c, ok := rune2cell[r]
			if !ok {
				log.Fatalf("invalid rune %q", r)
			}
			st.grid = append(st.grid, c)
			if _, _, ok := c.cart(); ok {
				st.carts = append(st.carts, st.w*st.h+i)
			}
		}
		st.h++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return st
}
