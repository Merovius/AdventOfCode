package main

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	players, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total winnings:", Part1(players))
	fmt.Println("Total winnings with jokers:", Part2(players))
}

func Parse(b []byte) ([]Player, error) {
	players, err := parse.Lines(
		parse.Struct[Player](
			split.Fields,
			parse.Array[Hand, byte](split.Bytes, parse.Byte),
			parse.Signed[int],
		),
	)(string(bytes.TrimSpace(b)))
	if err != nil {
		return nil, err
	}
	return players, nil
}

type Player struct {
	Hand Hand
	Bid  int
}

type Hand [5]byte

func Part1(players []Player) int {
	players = slices.Clone(players)
	slices.SortFunc(players, func(a, b Player) int {
		return a.Hand.Cmp(b.Hand)
	})
	var total int
	for r, p := range players {
		total += (r + 1) * p.Bid
	}
	return total
}

func Part2(players []Player) int {
	players = slices.Clone(players)
	slices.SortFunc(players, func(a, b Player) int {
		return a.Hand.Cmp2(b.Hand)
	})
	var total int
	for r, p := range players {
		total += (r + 1) * p.Bid
	}
	return total
	return 0
}

func CmpCards(a, b byte) int {
	const strength = "23456789TJQKA"
	return cmp.Compare(strings.IndexByte(strength, a), strings.IndexByte(strength, b))
}

func CmpCards2(a, b byte) int {
	const strength = "J23456789TQKA"
	return cmp.Compare(strings.IndexByte(strength, a), strings.IndexByte(strength, b))
}

type Type int

const (
	HighCard Type = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func (h Hand) String() string {
	slices.SortFunc(h[:], CmpCards)
	return string(h[:])
}

func (h Hand) Type() Type {
	// look at sorted hand. Every type we check for rules out the lower types.
	slices.Sort(h[:])
	if h[0] == h[4] {
		return FiveOfAKind
	}
	if h[0] == h[3] || h[1] == h[4] {
		return FourOfAKind
	}
	if (h[0] == h[2] && h[3] == h[4]) || (h[0] == h[1] && h[2] == h[4]) {
		return FullHouse
	}
	if h[0] == h[2] || h[1] == h[3] || h[2] == h[4] {
		return ThreeOfAKind
	}
	// count pairs
	var np int
	for i := 1; i < 5; i++ {
		if h[i] == h[i-1] {
			np++
		}
	}
	switch np {
	case 2:
		return TwoPair
	case 1:
		return OnePair
	default:
		return HighCard
	}
}

func (h Hand) Type2() Type {
	slices.Sort(h[:])

	// remove jokers
	r := strings.ReplaceAll(string(h[:]), "J", "")

	if len(r) < 2 {
		// 4 or 5 jokers
		return FiveOfAKind
	}
	if len(r) == 2 {
		// 3 jokers
		if r[0] == r[1] {
			return FiveOfAKind
		}
		return FourOfAKind
	}
	if len(r) == 3 {
		// 2 jokers
		if r[0] == r[2] {
			return FiveOfAKind
		}
		if r[0] == r[1] || r[1] == r[2] {
			return FourOfAKind
		}
		// Full House can be re-jokered into Four of a kind
		return ThreeOfAKind
	}
	if len(r) == 4 {
		// 1 joker
		if r[0] == r[3] {
			return FiveOfAKind
		}
		if r[0] == r[2] || r[1] == r[3] {
			return FourOfAKind
		}
		// count pairs
		var np int
		for i := 1; i < len(r); i++ {
			if r[i] == r[i-1] {
				np++
			}
		}
		switch np {
		case 2:
			return FullHouse
		case 1:
			return ThreeOfAKind
		default:
			return OnePair
		}
	}
	// no jokers
	return h.Type()
}

func (h Hand) Cmp(g Hand) int {
	ht, gt := h.Type(), g.Type()
	if c := cmp.Compare(ht, gt); c != 0 {
		return c
	}
	return slices.CompareFunc(h[:], g[:], CmpCards)
}

func (h Hand) Cmp2(g Hand) int {
	ht, gt := h.Type2(), g.Type2()
	if c := cmp.Compare(ht, gt); c != 0 {
		return c
	}
	return slices.CompareFunc(h[:], g[:], CmpCards2)
}
