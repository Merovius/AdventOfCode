package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	P = 20201227
	G = 7
)

func main() {
	var (
		pk_door int
		pk_card int
	)

	log.SetFlags(log.Lshortfile)
	flag.IntVar(&pk_door, "door", 0, "The public key of the door")
	flag.IntVar(&pk_card, "card", 0, "The public key of the card")
	flag.Parse()
	if pk_door == 0 || pk_card == 0 {
		log.Fatal("Both -door and -card must be provided")
	}
	ls_door := Crack(G, pk_door, P)
	ls_card := Crack(G, pk_card, P)
	s_door := Transform(pk_card, ls_door, P)
	s_card := Transform(pk_door, ls_card, P)
	fmt.Println("Secret of door:", s_door)
	fmt.Println("Secret of card:", s_card)
}

// Transform calculates sn^ls mod p. (sn = subject number, ls = loop size, p is a prime)
func Transform(sn, ls, p int) int {
	v := 1
	for i := 0; i < ls; i++ {
		v = (v * sn) % p
	}
	return v
}

// Crack tries to find ls such that sn^ls mod p == x.
func Crack(sn, x, p int) (ls int) {
	v := 1
	for {
		if x == v {
			return ls
		}
		v = (v * sn) % p
		ls++
	}
}
