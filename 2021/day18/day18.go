package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func main() {
	ns, err := parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Magnitude of sum:", Sum(ns...).Magnitude())
	var max = math.MinInt
	for _, n := range ns {
		for _, m := range ns {
			if v := Add(n, m).Magnitude(); v > max {
				max = v
			}
		}
	}
	fmt.Println("Maximum magnitude:", max)
}

type Num struct {
	Val int
	L   *Num
	R   *Num
}

func (n *Num) isRegular() bool {
	return n.L == nil && n.R == nil
}

func (n *Num) addLeft(v int) {
	if n.isRegular() {
		n.Val += v
		return
	}
	n.L.addLeft(v)
}

func (n *Num) addRight(v int) {
	if n.isRegular() {
		n.Val += v
		return
	}
	n.R.addRight(v)
}

func (n *Num) explodeRec(level int) (exploded bool, l, r int) {
	if n.isRegular() {
		return false, -1, -1
	}
	if level == 4 {
		l, r = n.L.Val, n.R.Val
		n.Val, n.L, n.R = 0, nil, nil
		return true, l, r
	}
	if x, l, r := n.L.explodeRec(level + 1); x {
		if r >= 0 {
			n.R.addLeft(r)
		}
		return true, l, -1
	}
	if x, l, r := n.R.explodeRec(level + 1); x {
		if l >= 0 {
			n.L.addRight(l)
		}
		return true, -1, r
	}
	return false, -1, -1
}

func (n *Num) explode() bool {
	x, _, _ := n.explodeRec(0)
	return x
}

func (n *Num) split() bool {
	if !n.isRegular() {
		return n.L.split() || n.R.split()
	}
	if n.Val < 10 {
		return false
	}
	l := n.Val / 2
	n.L = &Num{Val: l}
	n.R = &Num{Val: n.Val - l}
	n.Val = 0
	return true
}

func (n *Num) reduce() {
	for {
		if n.explode() {
			continue
		}
		if n.split() {
			continue
		}
		break
	}
}

func (n *Num) copy() *Num {
	if n == nil {
		return nil
	}
	return &Num{Val: n.Val, L: n.L.copy(), R: n.R.copy()}
}

func Add(a, b *Num) *Num {
	n := &Num{L: a.copy(), R: b.copy()}
	n.reduce()
	return n
}

func Sum(ns ...*Num) *Num {
	out, ns := ns[0], ns[1:]
	for _, n := range ns {
		m := Add(out, n)
		out = m
	}
	return out
}

func (n *Num) Magnitude() int {
	if n.isRegular() {
		return n.Val
	}
	return 3*n.L.Magnitude() + 2*n.R.Magnitude()
}

func (n *Num) String() string {
	if n == nil {
		return "<nil>"
	}
	if n.L == nil && n.R == nil {
		return fmt.Sprint(n.Val)
	}
	return fmt.Sprintf("(%v,%v)", n.L, n.R)
}

func (n *Num) UnmarshalJSON(b []byte) error {
	if b[0] != '[' {
		return json.Unmarshal(b, &n.Val)
	}
	var l [2]Num
	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}
	n.L, n.R = &l[0], &l[1]
	return nil
}

func parse(r io.Reader) ([]*Num, error) {
	var out []*Num

	s := bufio.NewScanner(r)
	for s.Scan() {
		n := new(Num)
		if err := json.Unmarshal(s.Bytes(), n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, s.Err()
}
