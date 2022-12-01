package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func main() {
	c, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	a := Aggregate(c)
	fmt.Printf("The most heavily loaded elf carries %d calories\n", Max(a))
	if len(a) < 3 {
		log.Fatal("Less than 3 total elves in input")
	}
	slices.Sort(a)
	var t int
	for i := 1; i <= 3; i++ {
		t += a[len(a)-i]
	}
	fmt.Printf("The three most heavily loaded elves carry %d calories\n", t)
}

func ReadInput(r io.Reader) ([][]int, error) {
	var out [][]int

	s := bufio.NewScanner(r)
	for s.Scan() {
		t := strings.TrimSpace(s.Text())
		if t == "" {
			out = append(out, nil)
			continue
		}
		n, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		if len(out) == 0 {
			out = append(out, nil)
		}
		out[len(out)-1] = append(out[len(out)-1], n)
	}
	return out, s.Err()
}

func Aggregate(c [][]int) []int {
	out := make([]int, len(c))
	for i, s := range c {
		for _, n := range s {
			out[i] += n
		}
	}
	return out
}

func Max[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		panic("Max called on empty slice")
	}
	v := s[0]
	for _, w := range s[1:] {
		if w > v {
			v = w
		}
	}
	return v
}
