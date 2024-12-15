package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
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
	fmt.Println(Part1(in.Clone()))
	fmt.Println(Part2(in))
}

func Parse(buf []byte) (Input, error) {
	type input struct {
		G *grid.Grid[Cell]
		M [][]grid.Direction
	}
	in, err := parse.Struct[input](
		split.Blocks,
		func(in string) (*grid.Grid[Cell], error) {
			return grid.Read(strings.NewReader(in), func(r rune) (Cell, error) {
				if i := strings.Index(".@O#", string(r)); i >= 0 {
					return Cell(i), nil
				}
				return 0, fmt.Errorf("unknown cell %q", r)
			})
		},
		parse.Lines(
			parse.Slice(
				split.Bytes,
				func(s string) (grid.Direction, error) {
					if i := strings.Index("^>v<", s); i >= 0 {
						return grid.Direction(i), nil
					}
					return 0, fmt.Errorf("unknown direction %q", s)
				},
			),
		),
	)(string(buf))
	if err != nil {
		return Input{}, err
	}
	i := slices.Index(in.G.G, Robot)
	in.G.G[i] = Empty
	return Input{
		Pos:  in.G.Pos(i),
		G:    in.G,
		Move: slices.Concat(in.M...),
	}, nil
}

type Input struct {
	Pos  grid.Pos
	G    *grid.Grid[Cell]
	Move []grid.Direction
}

func (in Input) Clone() Input {
	return Input{
		Pos:  in.Pos,
		G:    in.G.Clone(),
		Move: slices.Clone(in.Move),
	}
}

type Cell byte

const (
	Empty Cell = iota
	Robot
	Box
	Wall
	BoxL
	BoxR
)

func Part1(in Input) int {
	return run(in.G, in.Pos, in.Move)
}

func Part2(in Input) int {
	// double grid
	G := make([]Cell, 0, 2*len(in.G.G))
	for _, c := range in.G.G {
		if c == Box {
			G = append(G, BoxL, BoxR)
		} else {
			G = append(G, c, c)
		}
	}
	in.G.G = G
	in.G.W = 2 * in.G.W
	in.Pos.Col *= 2

	return run(in.G, in.Pos, in.Move)
}

func run(g *grid.Grid[Cell], p grid.Pos, moves []grid.Direction) int {
	for _, m := range moves {
		if !canMove(g, p, m) {
			continue
		}
		q := m.Move(p)
		move(g, q, m)
		p = q
	}
	var total int
	for p, c := range g.All() {
		if c == BoxL || c == Box {
			total += 100*p.Row + p.Col
		}
	}
	return total
}

func canMove(g *grid.Grid[Cell], p grid.Pos, d grid.Direction) bool {
	q := d.Move(p)
	switch g.At(q) {
	case Wall:
		return false
	case Box:
		return canMove(g, q, d)
	case BoxL:
		if d == grid.Left || d == grid.Right {
			return canMove(g, d.Move(q), d)
		}
		return canMove(g, q, d) && canMove(g, grid.Right.Move(q), d)
	case BoxR:
		if d == grid.Left || d == grid.Right {
			return canMove(g, d.Move(q), d)
		}
		return canMove(g, q, d) && canMove(g, grid.Left.Move(q), d)
	}
	return true
}

func move(g *grid.Grid[Cell], p grid.Pos, d grid.Direction) {
	if g.At(p) == Empty {
		return
	}
	q := d.Move(p)
	move(g, q, d)
	if d == grid.Up || d == grid.Down {
		switch g.At(p) {
		case BoxL:
			move(g, grid.Right.Move(q), d)
			g.Set(grid.Right.Move(q), g.At(grid.Right.Move(p)))
			g.Set(grid.Right.Move(p), Empty)
		case BoxR:
			move(g, grid.Left.Move(q), d)
			g.Set(grid.Left.Move(q), g.At(grid.Left.Move(p)))
			g.Set(grid.Left.Move(p), Empty)
		}
	}
	g.Set(q, g.At(p))
	g.Set(p, Empty)
}
