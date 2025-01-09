package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in := string(buf)
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Part1(in string) int {
	var n int
	for _, b := range in {
		switch b {
		case '(':
			n++
		case ')':
			n--
		}
	}
	return n
}

func Part2(in string) int {
	var n int
	for i, b := range in {
		switch b {
		case '(':
			n++
		case ')':
			n--
		}
		if n < 0 {
			return i + 1
		}
	}
	return 0
}
