package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	tiles, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g := make(Grid)
	for _, t := range tiles {
		g.FlipTile(t)
	}
	fmt.Println("Number of green tiles:", g.CountGreenTiles())
	for i := 1; i <= 100; i++ {
		g = g.DailyFlips()
	}
	fmt.Println("Number of green tiles on day 100:", g.CountGreenTiles())
}

type Hex [3]int

func (h Hex) Walk(d Direction) Hex {
	switch d {
	case E:
		return Hex{h[0] + 1, h[1], h[2] - 1}
	case W:
		return Hex{h[0] - 1, h[1], h[2] + 1}
	case SE:
		return Hex{h[0], h[1] + 1, h[2] - 1}
	case NW:
		return Hex{h[0], h[1] - 1, h[2] + 1}
	case SW:
		return Hex{h[0] - 1, h[1] + 1, h[2]}
	case NE:
		return Hex{h[0] + 1, h[1] - 1, h[2]}
	default:
		panic(fmt.Errorf("unknown direction %d", d))
	}
}

type Direction int

const (
	E Direction = iota
	SE
	SW
	W
	NW
	NE
	nDirection
)

func (d Direction) String() string {
	switch d {
	case E:
		return "E"
	case SE:
		return "SE"
	case SW:
		return "SW"
	case W:
		return "W"
	case NW:
		return "NW"
	case NE:
		return "NE"
	default:
		panic(fmt.Errorf("invalid direction %d", d))
	}
}

type Color int

const (
	Blue Color = iota
	Green
)

type Grid map[Hex]Color

func ReadInput(r io.Reader) ([][]Direction, error) {
	var out [][]Direction
	s := bufio.NewScanner(r)
	for s.Scan() {
		var t []Direction
		l := s.Text()
		for len(l) > 0 {
			switch l[0] {
			case 'e':
				t = append(t, E)
				l = l[1:]
				continue
			case 'w':
				t = append(t, W)
				l = l[1:]
				continue
			case 'n':
			case 's':
			default:
				return nil, fmt.Errorf("invalid character %q in input", l[0])
			}
			if len(l) == 1 {
				return nil, fmt.Errorf("truncated direction %q", l[0])
			}
			switch {
			case l[0] == 'n' && l[1] == 'e':
				t = append(t, NE)
			case l[0] == 'n' && l[1] == 'w':
				t = append(t, NW)
			case l[0] == 's' && l[1] == 'e':
				t = append(t, SE)
			case l[0] == 's' && l[1] == 'w':
				t = append(t, SW)
			}
			l = l[2:]
		}
		out = append(out, t)
	}
	return out, s.Err()
}

func (g Grid) Flip(h Hex) {
	if g[h] == Green {
		g[h] = Blue
	} else {
		g[h] = Green
	}
}

func (g Grid) FlipTile(ds []Direction) {
	h := Hex{}
	for _, d := range ds {
		h = h.Walk(d)
	}
	g.Flip(h)
}

func (g Grid) CountGreenTiles() int {
	var n int
	for _, c := range g {
		if c == Green {
			n++
		}
	}
	return n
}

func (g Grid) DailyFlips() Grid {
	out := make(Grid)
	for h := range g {
		for d := Direction(0); d < nDirection; d++ {
			n := h.Walk(d)
			var nb int
			for d := Direction(0); d < nDirection; d++ {
				if g[n.Walk(d)] == Green {
					nb++
				}
			}
			if g[n] == Green && nb != 0 && nb <= 2 {
				out[n] = Green
			} else if g[n] == Blue && nb == 2 {
				out[n] = Green
			}
		}
	}
	return out
}
