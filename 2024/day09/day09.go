package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/Merovius/AdventOfCode/internal/interval"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	buf = bytes.TrimSpace(buf)

	fmt.Println(Part1(buf))
	fmt.Println(Part2(buf))
}

func Part1(input []byte) int {
	var blocks []int
	for i, b := range input {
		for range int(b - '0') {
			if i%2 == 0 {
				blocks = append(blocks, i/2)
			} else {
				blocks = append(blocks, -1)
			}
		}
	}
	i, j := 0, len(blocks)-1
	for {
		// find next empty block
		for ; i <= j; i++ {
			if blocks[i] == -1 {
				break
			}
		}
		// find last file block
		for ; j >= i; j-- {
			if blocks[j] != -1 {
				break
			}
		}
		if i >= j {
			break
		}
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
	var check int
	for i, v := range blocks {
		if v != -1 {
			check += i * v
		}
	}
	return check
}

type Interval = interval.CO[int]

func Part2(input []byte) int {
	var (
		empty = make([]Interval, 0, len(input)/2)
		files = make([]Interval, 0, (len(input)+1)/2)
		offs  int
	)
	for i, b := range input {
		f := Interval{offs, offs + int(b-'0')}
		if i%2 == 0 {
			files = append(files, f)
		} else {
			empty = append(empty, f)
		}
		offs = f.Max
	}
	for j, f := range slices.Backward(files) {
		for i, e := range empty {
			if e.Len() >= f.Len() {
				files[j].Min, files[j].Max = e.Min, e.Min+f.Len()
				empty[i].Min = files[j].Max
				break
			}
			if e.Min >= f.Min {
				break
			}
		}
	}
	var check int
	for i, f := range files {
		// we add all block-indices of the file, which can be
		// simplified using triangle numbers.
		t1 := f.Min * (f.Min - 1) / 2
		t2 := f.Max * (f.Max - 1) / 2
		check += (t2 - t1) * i
	}
	return check
}
