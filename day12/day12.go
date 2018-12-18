package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
)

func main() {
	N := flag.Int("n", 20, "Number of steps")
	flag.Parse()

	st, p := readProg()
	for j := 1; j <= *N; j++ {
		st.step(p)
	}
	fmt.Println(st.Sum())
}

type prog uint32

func (p prog) print() {
	for i := uint8(0); i < 32; i++ {
		fmt.Print(i2s(i), " => ")
		if p&(1<<i) != 0 {
			fmt.Println("#")
		} else {
			fmt.Println(".")
		}
	}
}

func i2s(i uint8) string {
	var buf [5]byte
	for j := 0; j < 5; j++ {
		if i&(0x10>>uint(j)) == 0 {
			buf[j] = '.'
		} else {
			buf[j] = '#'
		}
	}
	return string(buf[:])
}

type state struct {
	a  []int
	a2 []int
}

func (s *state) step(p prog) {
	const (
		infty = int(^uint(0) >> 1)
	)

	if p&1 != 0 {
		panic("program not well-defined")
	}
	inst := uint8(0)
	tape := s.a
	idx := 0
	for inst != 0 || len(tape) > 0 {
		next := infty
		if len(tape) > 0 {
			next = tape[0]
		}
		if inst == 0 {
			idx = next - 2
		}
		if next-idx < 3 {
			inst |= 1 << uint(2+idx-next)
			tape = tape[1:]
		}
		if p&(1<<inst) != 0 {
			s.a2 = append(s.a2, idx)
		}
		idx++
		inst = (inst << 1) % 32
	}
	s.a, s.a2 = s.a2, s.a[:0]
}

func (s *state) Sum() int {
	total := 0
	for _, v := range s.a {
		total += v
	}
	return total
}

func (s *state) String() string {
	var b strings.Builder
	l := -2
	h := 3
	if len(s.a) > 0 {
		l = s.a[0] - 2
		h = s.a[len(s.a)-1] + 3
	}
	tape := s.a
	for i := l; i < h; i++ {
		if i == 0 {
			b.WriteString("\x1B[1;31m")
		}
		if len(tape) > 0 && i == tape[0] {
			b.WriteString("#")
			tape = tape[1:]
		} else {
			b.WriteString(".")
		}
		if i == 0 {
			b.WriteString("\x1B[m")
		}
	}
	return b.String()
}

func readProg() (*state, prog) {
	st := new(state)

	var i string
	fmt.Scanf("initial state: %s\n", &i)
	fmt.Scanf("\n")
	for j, r := range i {
		if r == '#' {
			st.a = append(st.a, j)
		}
	}

	var p prog
	for i := 0; i < 32; i++ {
		var match, alive string
		_, err := fmt.Scanf("%s => %s\n", &match, &alive)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var m uint8
		for _, r := range match {
			m <<= 1
			if r == '#' {
				m |= 1
			}
		}
		if alive == "#" {
			p |= (1 << m)
		}
	}
	return st, p
}
