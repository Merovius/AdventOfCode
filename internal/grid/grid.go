package grid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/slices"
)

type Pos struct {
	Row int
	Col int
}

func (p Pos) Add(δ Pos) Pos {
	return Pos{
		Row: p.Row + δ.Row,
		Col: p.Col + δ.Col,
	}
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row+1, p.Col+1)
}

type Grid[T any] struct {
	W int
	H int
	G []T
}

func Read[T any](r io.Reader, fromRune func(rune) (T, error)) (*Grid[T], error) {
	var lines [][]rune
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		if !utf8.ValidString(l) {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		lines = append(lines, []rune(l))
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("empty grid")
	}
	g := &Grid[T]{
		W: len(lines[0]),
		H: len(lines),
		G: make([]T, len(lines)*len(lines[0])),
	}
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			v, err := fromRune(lines[r][c])
			if err != nil {
				return nil, err
			}
			g.Set(Pos{r, c}, v)
		}
	}
	return g, nil
}

func New[T any](w, h int) *Grid[T] {
	return &Grid[T]{
		W: w,
		H: h,
		G: make([]T, w*h),
	}
}

func (g *Grid[T]) At(p Pos) T {
	return g.G[g.W*p.Row+p.Col]
}

func (g *Grid[T]) Set(p Pos, v T) {
	g.G[g.W*p.Row+p.Col] = v
}

func (g *Grid[T]) Clone() *Grid[T] {
	return &Grid[T]{
		W: g.W,
		H: g.H,
		G: slices.Clone(g.G),
	}
}

func (g *Grid[T]) Pos(i int) Pos {
	return Pos{
		Row: i / g.W,
		Col: i % g.W,
	}
}

func (g *Grid[T]) Valid(p Pos) bool {
	return p.Row >= 0 && p.Row < g.H && p.Col >= 0 && p.Col < g.W
}

func (g *Grid[T]) Neigh4(p Pos) []Pos {
	var out []Pos
	for _, δ := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		n := p.Add(Pos{δ[0], δ[1]})
		if !g.Valid(n) {
			continue
		}
		out = append(out, n)
	}
	return out

}
