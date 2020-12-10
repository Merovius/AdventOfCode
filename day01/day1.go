// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	var in []int
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		v, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		in = append(in, v)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	pair, ok := Pair2020(in)
	if !ok {
		log.Fatal("Could not find pair summing to 2020")
	}
	triple, ok := Triple2020(in)
	if !ok {
		log.Fatal("Could not find triple summing to 2020")
	}
	fmt.Printf("%d+%d=2020, %d*%d = %d\n", pair[0], pair[1], pair[0], pair[1], pair[0]*pair[1])
	fmt.Printf("%d+%d+%d=2020, %d*%d*%d = %d\n", triple[0], triple[1], triple[2], triple[0], triple[1], triple[2], triple[0]*triple[1]*triple[2])
}
