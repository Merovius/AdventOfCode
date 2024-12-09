package main

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"log"
	"math"
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

type Interval = interval.CO[int]

func checksum(id int, i Interval) int {
	// We add all block-indices of the file, which can be
	// simplified using triangle numbers.
	t1 := i.Min * (i.Min - 1)
	t2 := i.Max * (i.Max - 1)
	return (t2 - t1) / 2 * id
}

func Part1(input []byte) int {
	var (
		// Keep track of free blocks, removing them as they are used.
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
	var check int
	for j, f := range slices.Backward(files) {
		for f.Len() > 0 {
			e, l := empty[0], f.Len()
			if e.Min >= f.Max {
				// First free block starts after f. Rest of f
				// has to stay where it is
				check += checksum(j, f)
				break
			}
			// Use as much as possible of the free block for this
			// file
			l = min(e.Len(), l)
			i := Interval{e.Min, e.Min + l}
			check += checksum(j, i)

			// Free space gets used from the start, files get
			// consumed from the end.
			e.Min += l
			f.Max -= l

			// If this free block is used up, discard it.
			if e.Empty() {
				empty = empty[1:]
			} else {
				empty[0] = e
			}
		}
	}
	return check
}

func Part2(input []byte) int {
	var (
		// per-size list of free sections, ordered by start
		empty = make([][]Interval, 10)
		files = make([]Interval, 0, (len(input)+1)/2)
		offs  int
	)
	for i, b := range input {
		f := Interval{offs, offs + int(b-'0')}
		if i%2 == 0 {
			files = append(files, f)
		} else if l := f.Len(); l > 0 {
			empty[l] = append(empty[l], f)
		}
		offs = f.Max
	}
	var check int
	for j, f := range slices.Backward(files) {
		// find the earliest free block with enough capacity for f
		e, l := Interval{math.MaxInt, math.MaxInt}, 0
		for i := f.Len(); i < len(empty); i++ {
			s := empty[i]
			if len(empty[i]) == 0 || s[0].Min >= f.Min {
				continue
			}
			if s[0].Min < e.Min {
				e, l = s[0], i
			}
		}
		if e.Min == math.MaxInt {
			// all free blocks of sufficient size are after f
			check += checksum(j, f)
			continue
		}
		// remove free block from its list and move file into it
		empty[l] = empty[l][1:]
		f.Min, f.Max = e.Min, e.Min+f.Len()
		check += checksum(j, f)
		e.Min = f.Max

		// add free block to the list for its new size
		if l := e.Len(); l > 0 {
			i, _ := slices.BinarySearchFunc(empty[l], e.Min, func(i Interval, n int) int {
				return cmp.Compare(i.Min, n)
			})
			empty[l] = slices.Insert(empty[l], i, e)
		}
	}
	return check
}
