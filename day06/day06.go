package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Merovius/aoc_2022/internal/set"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	input = bytes.TrimSpace(input)
	m, ok := Marker(input, 4)
	if ok {
		fmt.Printf("Packet marker after %d characters\n", m)
	} else {
		fmt.Println("No packet marker found")
	}
	m, ok = Marker(input, 14)
	if ok {
		fmt.Printf("Message marker after %d characters\n", m)
	} else {
		fmt.Println("No message marker found")
	}
}

func Marker(buf []byte, l int) (int, bool) {
	for i := l; i < len(buf); i++ {
		if len(set.Make(buf[i-l:i]...)) == l {
			return i, true
		}
	}
	return 0, false
}
