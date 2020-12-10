package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	var ids []SeatID

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ids = append(ids, ParseSeatID(s.Text()))
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	fmt.Printf("Maximum seat ID: %v\n", ids[len(ids)-1])
	for i := 1; i < len(ids); i++ {
		if ids[i] != ids[i-1]+1 {
			fmt.Printf("Your seat ID: %v\n", ids[i-1]+1)
			return
		}
	}
}

type SeatID uint16

func ParseSeatID(s string) SeatID {
	var id SeatID
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'B', 'R':
			id = (id << 1) | 1
		case 'F', 'L':
			id = (id << 1) | 0
		default:
			panic("invalid seat spec")
		}
	}
	return id
}
