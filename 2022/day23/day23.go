package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Merovius/AdventOfCode/internal/frame"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/math"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	debug = flag.Bool("debug", false, "Output debug logs")
	logf  = func(format string, args ...any) {}
)

func main() {
	flag.Parse()
	if *debug {
		log.SetFlags(log.Lshortfile)
		logf = log.Printf
	}
	g, err := grid.Read(os.Stdin, func(r rune) (Cell, error) {
		if i := strings.IndexRune(".#", r); i >= 0 {
			return Cell(i), nil
		}
		return 0, fmt.Errorf("invalid cell %q", r)
	})
	if err != nil {
		log.Fatal(err)
	}
	s := NewSparse(g)
	s2 := s.Clone()
	Move(s, 10)
	fmt.Println("Number of ground tiles:", CountEmptyTiles(s))
	fmt.Println("Number of rounds until elves stop moving:", Move(s2, -1))
}

type Cell byte

const (
	Ground Cell = iota
	Elf
)

func Move(g *Sparse[Cell], n int) int {
	consider := [][]grid.Pos{
		{{-1, -1}, {-1, 0}, {-1, 1}},
		{{1, -1}, {1, 0}, {1, 1}},
		{{-1, -1}, {0, -1}, {1, -1}},
		{{-1, 1}, {0, 1}, {1, 1}},
	}
	if n < 0 {
		n = math.MaxInt
	}
	for c := 0; c < n; c++ {
		logf("c = %d", c)
		Dump(g)
		props := make(map[grid.Pos]int)
		moves := make(map[grid.Pos]grid.Pos)
	proposals:
		for e := range g.m {
			isolated := !slices.ContainsFunc(g.Neigh8(e), func(p grid.Pos) bool {
				return g.At(p) != Ground
			})
			if isolated {
				continue proposals
			}
		considerations:
			for i := range consider {
				check := consider[(c+i)%len(consider)]
				for _, δ := range check {
					if g.At(e.Add(δ)) != Ground {
						continue considerations
					}
				}
				t := e.Add(check[1])
				props[t]++
				moves[e] = t
				break
			}
		}
		if len(props) == 0 {
			logf("c = %d", c+1)
			Dump(g)
			return c + 1
		}
		for e, t := range moves {
			if props[t] == 1 {
				g.Set(e, Ground)
				g.Set(t, Elf)
			}
		}
	}
	logf("c = %d", n)
	Dump(g)
	return n
}

func CountEmptyTiles(g *Sparse[Cell]) int {
	min, max := g.Bounds()
	return (max.Row-min.Row)*(max.Col-min.Col) - len(g.m)
}

type Sparse[T comparable] struct {
	m map[grid.Pos]T
}

func NewSparse[T comparable](g *grid.Grid[T]) *Sparse[T] {
	m := make(map[grid.Pos]T)
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			p := grid.Pos{Row: r, Col: c}
			if v := g.At(p); v != *new(T) {
				m[p] = v
			}
		}
	}
	return &Sparse[T]{m}
}

func (g *Sparse[T]) Clone() *Sparse[T] {
	return &Sparse[T]{maps.Clone(g.m)}
}

func (g *Sparse[T]) At(p grid.Pos) T {
	return g.m[p]
}

func (g *Sparse[T]) Set(p grid.Pos, v T) {
	if v == *new(T) {
		delete(g.m, p)
	} else {
		g.m[p] = v
	}
}

func (g *Sparse[T]) Neigh8(p grid.Pos) []grid.Pos {
	var out []grid.Pos
	for δr := -1; δr <= 1; δr++ {
		for δc := -1; δc <= 1; δc++ {
			if δr == 0 && δc == 0 {
				continue
			}
			out = append(out, p.Add(grid.Pos{δr, δc}))
		}
	}
	return out
}

func (g *Sparse[T]) Bounds() (min, max grid.Pos) {
	min.Row, min.Col = math.MaxInt, math.MaxInt
	max.Row, max.Col = math.MinInt, math.MinInt
	for e := range g.m {
		min.Row = math.Min(min.Row, e.Row)
		min.Col = math.Min(min.Col, e.Col)
		max.Row = math.Max(max.Row, e.Row+1)
		max.Col = math.Max(max.Col, e.Col+1)
	}
	return min, max
}

func Dump(g *Sparse[Cell]) {
	if !*debug {
		return
	}
	defer time.Sleep(time.Second)
	w := frame.New(os.Stdout, frame.Simple)
	defer w.Close()
	min, max := g.Bounds()
	for r := min.Row - 1; r < max.Row+1; r++ {
		for c := min.Col - 1; c < max.Col+1; c++ {
			switch g.At(grid.Pos{r, c}) {
			case Ground:
				io.WriteString(w, " ")
			case Elf:
				io.WriteString(w, "█")
			default:
				panic("invalid Cell")
			}
		}
		io.WriteString(w, "\n")
	}
}
