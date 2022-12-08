package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Predicted score: %d\n", Score(data))
	fmt.Printf("Real score: %d\n", RealScore(data))
}

func ReadInput(r io.Reader) ([][2]int, error) {
	var out [][2]int
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		a, b, ok := strings.Cut(l, " ")
		if !ok || len(a) != 1 || len(b) != 1 {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		x, y := int(a[0]-'A'), int(b[0]-'X')
		if x < 0 || x > 2 || y < 0 || y > 2 {
			return nil, fmt.Errorf("invalid line %q", l)
		}
		out = append(out, [2]int{x, y})
	}
	return out, s.Err()
}

func Score(guide [][2]int) int {
	var score int
	for _, r := range guide {
		score += RoundScore(r[0], r[1])
	}
	return score
}

func RealScore(guide [][2]int) int {
	var score int
	for _, r := range guide {
		play := (3 + r[0] + (r[1] - 1)) % 3
		score += RoundScore(r[0], play)
	}
	return score
}

func RoundScore(they, we int) int {
	// 0: draw, 1: we lost, 2: we won
	s := (3 + they - we) % 3
	switch s {
	case 0:
		s = 3
	case 1:
		s = 0
	case 2:
		s = 6
	default:
		panic("unreachable")
	}
	return s + 1 + we
}
