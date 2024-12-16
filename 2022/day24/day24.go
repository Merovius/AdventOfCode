package main

import (
	"flag"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/graph"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	if *debug {
		log.SetFlags(log.Lshortfile)
		logf = log.Printf
	}

	grd, err := grid.Read(os.Stdin, func(r rune) (rune, error) {
		if strings.IndexRune(".#^>v<", r) >= 0 {
			return r, nil
		}
		return 0, fmt.Errorf("invalid codepoint %q in input", r)
	})
	if err != nil {
		log.Fatal(err)
	}
	g := NewGraph(grd)
	start := State{
		time:   0,
		player: grid.Pos{-1, 0},
	}
	path := graph.AStar(g, start, g.Goal, g.Heuristic)
	fmt.Println("Part 1:", len(path))
	path = graph.AStar(g, start, g.Goal2, g.Heuristic2)
	fmt.Println("Part 2:", len(path))
}

var logf = func(format string, args ...any) {}

type Direction int

const (
	_ Direction = iota
	Up
	Right
	Down
	Left
	Wait
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "^"
	case Right:
		return ">"
	case Down:
		return "v"
	case Left:
		return "<"
	default:
		panic("invalid direction")
	}
}

type Blizzard struct {
	Direction Direction
	Start     grid.Pos
}

type Graph struct {
	h         int
	w         int
	blizzards []Blizzard
}

func NewGraph(g *grid.Grid[rune]) *Graph {
	var blizzards []Blizzard
	for r := 1; r < g.H-1; r++ {
		for c := 1; c < g.W-1; c++ {
			if i := strings.IndexRune("^>v<", g.At(grid.Pos{r, c})); i >= 0 {
				logf("blizzard at %v, going %v", grid.Pos{r - 1, c - 1}, Direction(i+1))
				blizzards = append(blizzards, Blizzard{
					Direction: Direction(1 + i),
					Start:     grid.Pos{r - 1, c - 1},
				})
			}
		}
	}
	return &Graph{
		h:         g.H - 2,
		w:         g.W - 2,
		blizzards: blizzards,
	}
}

type State struct {
	player grid.Pos
	time   int
	t1     int // reached goal first time
	t2     int // reached start
	t3     int // reached goal again
}

func (s State) String() string {
	return fmt.Sprintf("%d,%v", s.time, s.player)
}

type Edge struct {
	From State
	To   State
}

func (g *Graph) blizzardPos(b Blizzard, t int) grid.Pos {
	p := b.Start
	switch b.Direction {
	case Up:
		p.Row -= t
	case Right:
		p.Col += t
	case Down:
		p.Row += t
	case Left:
		p.Col -= t
	default:
		panic("invalid direction")
	}
	p.Row = ((p.Row % g.h) + g.h) % g.h
	p.Col = ((p.Col % g.w) + g.w) % g.w
	return p
}

func (g *Graph) startPos() grid.Pos {
	return grid.Pos{-1, 0}
}

func (g *Graph) goalPos() grid.Pos {
	return grid.Pos{g.h, g.w - 1}
}

func (g *Graph) Edges(s State) iter.Seq[Edge] {
	return func(yield func(Edge) bool) {
		logf("Edges(%+v)", s)
		next := set.Make(s.player)
		for _, δ := range []grid.Pos{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			p := s.player.Add(δ)
			if p.Row >= 0 && p.Row < g.h && p.Col >= 0 && p.Col < g.w {
				next.Add(p)
			}
			if p == g.goalPos() || p == g.startPos() {
				next.Add(p)
			}
		}
		for _, b := range g.blizzards {
			p := g.blizzardPos(b, s.time+1)
			if next.Contains(p) {
				logf("%v blocked by blizzard %v", p, b.Start)
				next.Delete(p)
			}
			if len(next) == 0 {
				break
			}
		}
		for p := range next {
			n := s
			n.time++
			n.player = p
			if n.t1 == 0 {
				if p == g.goalPos() {
					n.t1 = n.time
				}
			} else if n.t2 == 0 {
				if p == g.startPos() {
					n.t2 = n.time
				}
			} else if n.t3 == 0 {
				if p == g.goalPos() {
					n.t3 = n.time
				}
			}
			logf("%v → %v", s, n)
			if !yield(Edge{From: s, To: n}) {
				return
			}
		}
	}
}

func (g *Graph) From(e Edge) State {
	return e.From
}

func (g *Graph) To(e Edge) State {
	return e.To
}

func (g *Graph) Weight(e Edge) int {
	return 1
}

func (g *Graph) Goal(s State) bool {
	return s.t1 != 0
}

func (g *Graph) Goal2(s State) bool {
	return s.t3 != 0
}

func (g *Graph) Heuristic(s State) (h int) {
	return g.goalPos().Dist(s.player)
}

func (g *Graph) Heuristic2(s State) (h int) {
	diam := g.w * (g.h + 1)
	switch {
	case s.t1 == 0:
		return 2*diam + g.goalPos().Dist(s.player)
	case s.t2 == 0:
		return diam + g.startPos().Dist(s.player)
	case s.t3 == 0:
		return g.goalPos().Dist(s.player)
	default:
		return 0
	}
}
