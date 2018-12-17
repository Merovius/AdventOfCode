package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	var arr []string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		arr = append(arr, s.Text())
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if hamming(arr[i], arr[j]) == 1 {
				fmt.Println(common(arr[i], arr[j]))
				return
			}
		}
	}
}

func hamming(s, t string) int {
	if len(s) != len(t) {
		log.Fatal("len(s) != len(t)")
	}
	var dist int
	for i := range s {
		if s[i] != t[i] {
			dist += 1
		}
	}
	return dist
}

func common(s, t string) string {
	var out []byte
	for i := range s {
		if s[i] == t[i] {
			out = append(out, s[i])
		}
	}
	return string(out)
}
