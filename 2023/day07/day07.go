package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/slices"
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
	return parse.Lines(
		parse.Struct[Player](
			split.Fields,
			parse.Array[Hand, byte](split.Bytes, parse.Byte),
			parse.Signed[int],
		),
	)(string(bytes.TrimSpace(b)))
}

type Player struct {
	Hand Hand
	Bid  int
}

type Hand [5]byte

func Part1(players []Player) int {
	return winnings(players, false)
}

func Part2(players []Player) int {
	return winnings(players, true)
}

func winnings(players []Player, joker bool) int {
	vals := make([]uint32, len(players))
	for i, p := range players {
		vals[i] = p.Hand.Value(joker)
	}
	slices.SortBy(players, vals)
	var total int
	for r, p := range players {
		total += (r + 1) * p.Bid
	}
	return total
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
	return string(h[:])
}

const (
	Order1 = "23456789TJQKA"
	Order2 = "J23456789TQKA"
)

// handType maps the highest two counts of card values to the type of hand.
// The second highest count is at most 2, so we can encode them into a uint8
// as (m1<<2)+m2
var handType = [...]Type{
	(5 << 2) + 0: FiveOfAKind,
	(4 << 2) + 1: FourOfAKind,
	(3 << 2) + 2: FullHouse,
	(3 << 2) + 1: ThreeOfAKind,
	(2 << 2) + 2: TwoPair,
	(2 << 2) + 1: OnePair,
	(1 << 2) + 1: HighCard,
}

// Value returns a value for h that can be used to rank Hands. A better Hand
// will always have a higher Value.
func (h Hand) Value(joker bool) uint32 {
	// Find the two most common card values. They determine the hand type.
	var (
		njokers uint8
		counts  [13]uint8
		maxC    [2]uint8
		ord     = Order1
	)
	if joker {
		ord = Order2
	}
	for i, b := range h {
		h[i] = uint8(strings.IndexByte(ord, b))
		counts[h[i]]++
	}
	for i, c := range counts {
		if joker && i == 0 {
			njokers = c
			continue
		}
		if c > maxC[1] {
			if c > maxC[0] {
				maxC[0], maxC[1] = c, maxC[0]
			} else {
				maxC[1] = c
			}
		}
	}
	if joker {
		// jokers should always boost the most common card value, as we don't
		// have straights
		maxC[0] += njokers
	}

	// Encode the card values, in lexicographic order, as a Base 13 integer.
	// The hand type is the most significant digit, so it dominates the
	// magnitude of the value.
	v := uint32(handType[(maxC[0]<<2)+maxC[1]])
	for _, b := range h {
		v *= uint32(len(ord))
		v += uint32(b)
	}
	return v
}
