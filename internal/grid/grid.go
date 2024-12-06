package grid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/Merovius/AdventOfCode/internal/math"
	"golang.org/x/exp/constraints"
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

func (p Pos) Sub(δ Pos) Pos {
	return Pos{
		Row: p.Row - δ.Row,
		Col: p.Col - δ.Col,
	}
}

func (p Pos) Length() int {
	return math.Abs(p.Row) + math.Abs(p.Col)
}

func (p Pos) Dist(q Pos) int {
	return p.Sub(q).Length()
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row, p.Col)
}

// Enum return a function suitable to pass to Read, that accepts any of opts
// and turns it into its position in that list.
func Enum[T constraints.Integer](opts ...rune) func(rune) (T, error) {
	return func(r rune) (T, error) {
		for i, v := range opts {
			if v == r {
				return T(i), nil
			}
		}
		return *new(T), fmt.Errorf("invalid rune %q")
	}
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

func (g *Grid[T]) Bounds() Rectangle {
	return Rect(0, 0, g.H, g.W)
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
	out := make([]Pos, 0, 4)
	for _, δ := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		n := p.Add(Pos{δ[0], δ[1]})
		if !g.Valid(n) {
			continue
		}
		out = append(out, n)
	}
	return out
}

type Rectangle struct {
	Min Pos
	Max Pos
}

func Rect(rmin, cmin, rmax, cmax int) Rectangle {
	return Rectangle{Pos{rmin, cmin}, Pos{rmax, cmax}}.Canon()
}

func (r Rectangle) Add(p Pos) Rectangle {
	r.Min = r.Min.Add(p)
	r.Max = r.Max.Add(p)
	return r
}

func (r Rectangle) Canon() Rectangle {
	if r.Max.Row < r.Min.Row {
		r.Min.Row, r.Max.Row = r.Max.Row, r.Min.Row
	}
	if r.Max.Col < r.Min.Col {
		r.Min.Col, r.Max.Col = r.Max.Col, r.Min.Col
	}
	return r
}

func (r Rectangle) Contains(p Pos) bool {
	return p.Row >= r.Min.Row && p.Row < r.Max.Row && p.Col >= r.Min.Col && p.Col < r.Max.Col
}

func (r Rectangle) Width() int {
	return r.Max.Col - r.Min.Col
}

func (r Rectangle) Height() int {
	return r.Max.Row - r.Min.Row
}

func (r Rectangle) Empty() bool {
	return r.Max.Row <= r.Min.Row || r.Max.Col <= r.Min.Col
}

func (r Rectangle) Eq(s Rectangle) bool {
	return r == s || (r.Empty() && s.Empty())
}

func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}
	return s.Min.Row <= r.Min.Row && r.Max.Row <= s.Max.Row && s.Min.Col <= r.Min.Col && r.Max.Col <= s.Max.Col
}

func (r Rectangle) Inset(n int) Rectangle {
	if r.Height() < 2*n {
		r.Min.Row = (r.Min.Row + r.Max.Row) / 2
		r.Max.Row = r.Min.Row
	} else {
		r.Min.Row += n
		r.Max.Row -= n
	}
	if r.Width() < 2*n {
		r.Min.Col = (r.Min.Col + r.Max.Col) / 2
		r.Max.Col = r.Min.Col
	} else {
		r.Min.Col += n
		r.Max.Col -= n
	}
	return r
}

func (r Rectangle) Intersect(s Rectangle) Rectangle {
	r.Min.Row = max(r.Min.Row, s.Min.Row)
	r.Min.Col = max(r.Min.Col, s.Min.Col)
	r.Max.Row = min(r.Max.Row, s.Max.Row)
	r.Max.Col = min(r.Max.Col, s.Max.Col)
	if r.Empty() {
		return Rectangle{}
	}
	return r
}

func (r Rectangle) Hull(s Rectangle) Rectangle {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	r.Min.Row = min(r.Min.Row, s.Min.Row)
	r.Min.Col = min(r.Min.Col, s.Min.Col)
	r.Max.Row = max(r.Max.Row, s.Max.Row)
	r.Max.Col = max(r.Max.Col, s.Max.Col)
	return r
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() && r.Min.Col < s.Max.Col && s.Min.Col < r.Max.Col && r.Min.Row < s.Max.Row && s.Min.Row < r.Max.Row
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func (d Direction) Move(p Pos) Pos {
	switch d {
	case Up:
		p.Row--
	case Right:
		p.Col++
	case Down:
		p.Row++
	case Left:
		p.Col--
	default:
		panic("invalid Direction")
	}
	return p
}

func (d Direction) RotateRight() Direction {
	return (d + 1) % 4
}

func (d Direction) RotateLeft() Direction {
	return (d + 3) % 4
}

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
		return "?"
	}
}
