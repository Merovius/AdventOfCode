package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/internal/math"
	"gonih.org/AdventOfCode/internal/xiter"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(in []byte) ([]string, error) {
	return strings.Split(string(bytes.TrimSpace(in)), "\n"), nil
}

func Part1(in []string) int {
	return Complexities(in, 2)
}

func Part2(in []string) int {
	return Complexities(in, 25)
}

func Complexities(in []string, lvl int) int {
	c := Hand()
	for range lvl {
		c = DirPad(c)
	}

	var total int
	for _, l := range in {
		best := xiter.FoldL(math.Min, PathCosts(l, c), math.MaxInt)
		i, err := strconv.Atoi(l[:len(l)-1])
		if err != nil {
			panic(err)
		}
		total += i * best
	}
	return total
}

// Controller of a pad. Calling it with a button to push returns the cost of
// pushing that button.
type Controller func(d Dir) int

// Hand is the hand of the player, which moves instantanously and just pushes
// the button.
func Hand() Controller {
	return func(d Dir) int { return 1 }
}

// DirPad is a directional pad.
func DirPad(c Controller) func(b Dir) int {
	var memo [nDir][nDir]int
	pos := dA
	return func(b Dir) int {
		if m := memo[pos][b]; m > 0 {
			pos = b
			return m
		}
		if pos == b {
			memo[pos][b] = c(dA)
			return memo[pos][b]
		}
		var best int = math.MaxInt
		for _, p := range DirPaths[pos][b] {
			var n int
			for _, d := range p {
				n += c(d)
			}
			n += c(dA)
			best = min(best, n)
		}
		memo[pos][b] = best
		pos = b
		return best
	}
}

// PathCosts yields the cost of all paths that type in s on the number pad
// using c.
func PathCosts(s string, c Controller) iter.Seq[int] {
	return func(yield func(int) bool) {
		var rec func(n int, start Num, s string) bool
		rec = func(n int, start Num, s string) bool {
			if len(s) == 0 {
				return yield(n)
			}
			to := NumFromByte(s[0])
			for _, p := range NumPaths[start][to] {
				m := n
				for _, b := range p {
					m += c(b)
				}
				m += c(dA)
				if !rec(m, to, s[1:]) {
					return false
				}
			}
			return true
		}
		rec(0, nA, s)
	}
}

//    +---+---+---+
//    | 7 | 8 | 9 |
//    +---+---+---+
//    | 4 | 5 | 6 |
//    +---+---+---+
//    | 1 | 2 | 3 |
//    +---+---+---+
//        | 0 | A |
//        +---+---+

type Num byte

const (
	n0 Num = iota
	n1
	n2
	n3
	n4
	n5
	n6
	n7
	n8
	n9
	nA

	nNum
)

func NumFromByte(b byte) Num {
	if b == 'A' {
		return nA
	}
	if b >= '0' && b <= '9' {
		return Num(b - '0')
	}
	panic(fmt.Errorf("invalid number pad button %q", b))
}

var NumPaths = [nNum][nNum][][]Dir{
	n7: {
		n7: {path("")},
		n8: {path(">")},
		n9: {path(">>")},
		n4: {path("v")},
		n5: {path("v>"), path(">v")},
		n6: {path("v>>"), path(">>v")},
		n1: {path("vv")},
		n2: {path("vv>"), path(">vv")},
		n3: {path("vv>>"), path(">>vv")},
		n0: {path(">vvv")},
		nA: {path(">>vvv")},
	},
	n8: {
		n7: {path("<")},
		n8: {path("")},
		n9: {path(">")},
		n4: {path("<v"), path("v<")},
		n5: {path("v")},
		n6: {path("v>"), path(">v")},
		n1: {path("<vv"), path("vv<")},
		n2: {path("vv")},
		n3: {path("vv>"), path(">vv")},
		n0: {path("vvv")},
		nA: {path("vvv>"), path(">vvv")},
	},
	n9: {
		n7: {path("<<")},
		n8: {path("<")},
		n9: {path("")},
		n4: {path("<<v"), path("v<<")},
		n5: {path("<v"), path("v<")},
		n6: {path("v")},
		n1: {path("<<v"), path("v<<")},
		n2: {path("<vv"), path("vv<")},
		n3: {path("vv")},
		n0: {path("<vvv"), path("vvv<")},
		nA: {path("vvv")},
	},
	n4: {
		n7: {path("^")},
		n8: {path("^>"), path(">^")},
		n9: {path("^>>"), path(">>^")},
		n4: {path("")},
		n5: {path(">")},
		n6: {path(">>")},
		n1: {path("v")},
		n2: {path("v>"), path(">v")},
		n3: {path("v>>"), path(">>v")},
		n0: {path(">vv")},
		nA: {path(">>vv")},
	},
	n5: {
		n7: {path("^<"), path("<^")},
		n8: {path("^")},
		n9: {path("^>"), path(">^")},
		n4: {path("<")},
		n5: {path("")},
		n6: {path(">")},
		n1: {path("<v"), path("v<")},
		n2: {path("v")},
		n3: {path(">v"), path("v>")},
		n0: {path("vv")},
		nA: {path("vv>"), path(">vv")},
	},
	n6: {
		n7: {path("^<<"), path("<<^")},
		n8: {path("^<"), path("<^")},
		n9: {path("^")},
		n4: {path("<<")},
		n5: {path("<")},
		n6: {path("")},
		n1: {path("<<v"), path("v<<")},
		n2: {path("<v"), path("v<")},
		n3: {path("v")},
		n0: {path("<vv"), path("vv<")},
		nA: {path("vv")},
	},
	n1: {
		n7: {path("^^")},
		n8: {path("^^>"), path(">^^")},
		n9: {path("^^>>"), path(">>^^")},
		n4: {path("^")},
		n5: {path("^>"), path(">^")},
		n6: {path("^>>"), path(">>^")},
		n1: {path("")},
		n2: {path(">")},
		n3: {path(">>")},
		n0: {path(">v")},
		nA: {path(">>v")},
	},
	n2: {
		n7: {path("^^<"), path("<^^")},
		n8: {path("^^")},
		n9: {path("^^>"), path(">^^")},
		n4: {path("^<"), path("<^")},
		n5: {path("^")},
		n6: {path("^>"), path(">^")},
		n1: {path("<")},
		n2: {path("")},
		n3: {path(">")},
		n0: {path("v")},
		nA: {path("v>"), path(">v")},
	},
	n3: {
		n7: {path("^^<<"), path("<<^^")},
		n8: {path("^^<"), path("<^^")},
		n9: {path("^^")},
		n4: {path("^<<"), path("<<^")},
		n5: {path("^<"), path("<^")},
		n6: {path("^")},
		n1: {path("<<")},
		n2: {path("<")},
		n3: {path("")},
		n0: {path("<v"), path("v<")},
		nA: {path("v")},
	},
	n0: {
		n7: {path("^^^<")},
		n8: {path("^^^")},
		n9: {path("^^^>"), path(">^^^")},
		n4: {path("^^<")},
		n5: {path("^^")},
		n6: {path("^^>"), path(">^^")},
		n1: {path("^<")},
		n2: {path("^")},
		n3: {path("^>"), path(">^")},
		n0: {path("")},
		nA: {path(">")},
	},
	nA: {
		n7: {path("^^^<<")},
		n8: {path("^^^<"), path("<^^^")},
		n9: {path("^^^")},
		n4: {path("^^<<")},
		n5: {path("^^<"), path("<^^")},
		n6: {path("^^")},
		n1: {path("^<<")},
		n2: {path("^<"), path("<^")},
		n3: {path("^")},
		n0: {path("<")},
		nA: {path("")},
	},
}

//        +---+---+
//        | ^ | A |
//    +---+---+---+
//    | < | v | > |
//    +---+---+---+

type Dir byte

const (
	dU Dir = iota
	dL
	dD
	dR
	dA

	nDir
)

func DirFromByte(b byte) Dir {
	if i := strings.IndexByte("^<v>A", b); i >= 0 {
		return Dir(i)
	}
	panic(fmt.Errorf("invalid direction pad button %q", b))
}

var DirPaths = [nDir][nDir][][]Dir{
	dU: {
		dU: {path("")},
		dA: {path(">")},
		dL: {path("v<")},
		dD: {path("v")},
		dR: {path("v>"), path(">v")},
	},
	dA: {
		dU: {path("<")},
		dA: {path("")},
		dL: {path("v<<")},
		dD: {path("v<"), path("<v")},
		dR: {path("v")},
	},
	dL: {
		dU: {path(">^")},
		dA: {path(">>^")},
		dL: {path("")},
		dD: {path(">")},
		dR: {path(">>")},
	},
	dD: {
		dU: {path("^")},
		dA: {path("^>"), path(">^")},
		dL: {path("<")},
		dD: {path("")},
		dR: {path(">")},
	},
	dR: {
		dU: {path("^<"), path("<^")},
		dA: {path("^")},
		dL: {path("<<")},
		dD: {path("<")},
		dR: {path("")},
	},
}

func path(s string) []Dir {
	p := make([]Dir, len(s))
	for i, d := range []byte(s) {
		p[i] = DirFromByte(d)
	}
	return p
}
