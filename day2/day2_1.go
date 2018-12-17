package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	two, three := 0, 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		hist := make(map[rune]int)
		for _, r := range l {
			hist[r] += 1
		}
		for _, c := range hist {
			if c == 2 {
				two++
				break
			}
		}
		for _, c := range hist {
			if c == 3 {
				three++
				break
			}
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println(two * three)
}
