package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/math"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(in []byte) ([]int, error) {
	return parse.Slice(split.Lines, func(s string) (int, error) {
		if len(s) == 0 {
			return 0, errors.New("empty line")
		}
		n, err := strconv.Atoi(s[1:])
		if err != nil {
			return 0, err
		}
		switch s[0] {
		case 'L':
			return -n, nil
		case 'R':
			return n, nil
		default:
			return 0, fmt.Errorf("invalid direction %q", s[0])
		}
	}).Parse(bytes.NewReader(in))
}

type Rotation struct {
	Direction byte
	Clicks    int
}

func Part1(in []int) int {
	total := 50
	N := 0
	for _, d := range in {
		total = math.Mod(total+d, 100)
		if total == 0 {
			N++
		}
	}
	return N
}

func Part2(in []int) int {
	pos := 50
	N := 0
	for _, d := range in {
		N += math.Abs(d / 100)
		if pos == 0 {
			pos = math.Mod(d, 100)
			continue
		}
		pos += d % 100
		if pos <= 0 || pos >= 100 {
			N++
		}
		pos = math.Mod(pos, 100)
	}
	return N
}
