package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	log.SetFlags(log.Lshortfile)

	p1, p2, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	r1, r2 := PlayGame(cp(p1), cp(p2))
	if len(r1) == 0 {
		fmt.Printf("Player 2 won with a score of %d\n", ScoreDeck(r2))
	} else {
		fmt.Printf("Player 1 won with a score of %d\n", ScoreDeck(r1))
	}
	r1, r2 = PlayGameRecursive(cp(p1), cp(p2))
	log.Println(r1)
	log.Println(r2)
	if len(r1) == 0 {
		fmt.Printf("Player 2 won the recursive game with a score of %d\n", ScoreDeck(r2))
	} else {
		fmt.Printf("Player 1 won the recursive game with a score of %d\n", ScoreDeck(r1))
	}
}

func ReadInput(r io.Reader) (p1, p2 []int, err error) {
	var player2 bool
	s := bufio.NewScanner(r)
	if !s.Scan() {
		return nil, nil, io.ErrUnexpectedEOF
	}
	if s.Text() != "Player 1:" {
		return nil, nil, fmt.Errorf("unexpected line %q want %q", s.Text(), "Player 1:")
	}
	for s.Scan() {
		l := s.Text()
		if l == "" {
			if !s.Scan() {
				return nil, nil, io.ErrUnexpectedEOF
			}
			if s.Text() != "Player 2:" {
				return nil, nil, fmt.Errorf("unexpected line %q want %q", s.Text(), "Player 2:")
			}
			player2 = true
			continue
		}
		n, err := strconv.Atoi(l)
		if err != nil {
			return nil, nil, fmt.Errorf("unexpected line %q, want number", l)
		}
		if player2 {
			p2 = append(p2, n)
		} else {
			p1 = append(p1, n)
		}
	}
	return p1, p2, s.Err()
}

func PlayRound(p1, p2 []int) (p1n, p2n []int) {
	c1, p1 := p1[0], p1[1:]
	c2, p2 := p2[0], p2[1:]
	if c1 > c2 {
		p1 = append(p1, c1, c2)
	} else if c1 < c2 {
		p2 = append(p2, c2, c1)
	} else {
		panic("tied cards in the deck")
	}
	return p1, p2
}

func PlayGame(p1, p2 []int) (p1n, p2n []int) {
	for len(p1) > 0 && len(p2) > 0 {
		p1, p2 = PlayRound(p1, p2)
	}
	return p1, p2
}

func ScoreDeck(deck []int) int {
	var score int
	for i, v := range deck {
		score += v * (len(deck) - i)
	}
	return score
}

func cp(v []int) []int {
	return append([]int(nil), v...)
}

func draw(deck []int) (c int, rest []int) {
	c, rest = deck[0], deck[1:]
	return c, rest
}

func PlayGameRecursive(p1, p2 []int) ([]int, []int) {
	pastStates := make(map[string]bool)
	for len(p1) > 0 && len(p2) > 0 {
		k := fmt.Sprintf("%v%v", p1, p2)
		if pastStates[k] {
			return append(p1, p2...), nil
		}
		pastStates[k] = true

		var c1, c2 int
		c1, p1 = draw(p1)
		c2, p2 = draw(p2)
		if len(p1) < c1 || len(p2) < c2 {
			if c1 > c2 {
				p1 = append(p1, c1, c2)
			} else if c1 < c2 {
				p2 = append(p2, c2, c1)
			} else {
				panic("tied cards in the deck")
			}
			continue
		}
		sp1, sp2 := PlayGameRecursive(cp(p1[:c1]), cp(p2[:c2]))
		if len(sp1) == 0 {
			p2 = append(p2, c2, c1)
		} else if len(sp2) == 0 {
			p1 = append(p1, c1, c2)
		}
	}
	return p1, p2
}
