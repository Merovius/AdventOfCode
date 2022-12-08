package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
)

func main() {
	groups, err := ParseGroups(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var (
		any int
		all int
	)
	for _, g := range groups {
		any += bits.OnesCount32(CollapseGroupAny(g))
		all += bits.OnesCount32(CollapseGroupAll(g))
	}
	fmt.Println("Total of any yes answers:", any)
	fmt.Println("Total of all yes answers:", all)
}

func ParseGroups(r io.Reader) ([][]uint32, error) {
	s := bufio.NewScanner(r)
	var group []uint32
	var out [][]uint32
	for s.Scan() {
		if s.Text() == "" {
			out = append(out, group)
			group = nil
			continue
		}
		var v uint32
		for _, b := range s.Bytes() {
			v |= 1 << (b - 'a')
		}
		group = append(group, v)
	}
	out = append(out, group)
	return out, s.Err()
}

func CollapseGroupAny(g []uint32) uint32 {
	var out uint32
	for _, v := range g {
		out |= v
	}
	return out
}

func CollapseGroupAll(g []uint32) uint32 {
	var out uint32 = 0x3ffffff
	for _, v := range g {
		out &= v
	}
	return out
}
