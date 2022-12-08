package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	m, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d trees are visible\n", CountVisible(m))
	fmt.Printf("Best scenic score is %d\n", BestScenicScore(m))
}

func ReadInput(r io.Reader) (*Map, error) {
	var out [][]int
	s := bufio.NewScanner(r)
	for s.Scan() {
		var row []int
		l := s.Text()
		for i := 0; i < len(l); i++ {
			if l[i] < '0' || l[i] > '9' {
				return nil, fmt.Errorf("invalid line %q", l)
			}
			row = append(row, int(l[i]-'0'))
		}
		out = append(out, row)
	}
	if len(out) == 0 {
		return nil, errors.New("empty input")
	}
	m := &Map{
		Width:  len(out[0]),
		Height: len(out),
		Cells:  out,
	}
	for _, r := range m.Cells[1:] {
		if len(r) != m.Width {
			return nil, errors.New("number of columns is inconsistent")
		}
	}
	return m, nil
}

type Map struct {
	Width  int
	Height int
	Cells  [][]int
}

func (m *Map) Valid(i, j int) bool {
	return i >= 0 && i < m.Height && j >= 0 && j < m.Width
}

type Ray struct {
	i  int
	j  int
	δi int
	δj int
}

func Rays(i, j int) [4]Ray {
	return [4]Ray{
		{i, j, -1, 0},
		{i, j, 1, 0},
		{i, j, 0, -1},
		{i, j, 0, 1},
	}
}

// Ray visits the ray from (i,j) in the direction (δi,δj), calling f with the
// height of each tree found. If f returns false, iteration is stopped.
func (m *Map) Visit(r Ray, f func(int) bool) {
	if (r.δi == 0) == (r.δj == 0) {
		panic("exactly one of r.δi or r.δj must be 0")
	}
	if r.δi < -1 || r.δi > 1 || r.δj < -1 || r.δj > 1 {
		panic("r.δi and r.δj must be in [-1,1]")
	}
	for i, j := r.i+r.δi, r.j+r.δj; m.Valid(i, j); i, j = i+r.δi, j+r.δj {
		if !f(m.Cells[i][j]) {
			return
		}
	}
}

func (m *Map) Visible(i, j int) bool {
	h := m.Cells[i][j]
	for _, r := range Rays(i, j) {
		visible := true
		m.Visit(r, func(hh int) bool {
			visible = visible && h > hh
			return visible
		})
		if visible {
			return true
		}
	}
	return false
}

func (m *Map) ViewingDistance(r Ray) int {
	d := 0
	h := m.Cells[r.i][r.j]
	m.Visit(r, func(hh int) bool {
		d++
		return hh < h
	})
	return d
}

func (m *Map) ScenicScore(i, j int) int {
	s := 1
	for _, r := range Rays(i, j) {
		s *= m.ViewingDistance(r)
	}
	return s
}

func CountVisible(m *Map) int {
	var n int
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			if m.Visible(i, j) {
				n++
			}
		}
	}
	return n
}

func BestScenicScore(m *Map) int {
	var best int
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			if s := m.ScenicScore(i, j); s > best {
				best = s
			}
		}
	}
	return best
}
