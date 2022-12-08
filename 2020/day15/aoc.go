package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	bs := bytes.Split(bytes.TrimSpace(buf), []byte{','})
	var ns []int
	for _, b := range bs {
		n, err := strconv.Atoi(string(b))
		if err != nil {
			log.Fatal(err)
		}
		ns = append(ns, n)
	}
	if len(ns) == 0 {
		log.Fatal("No numbers given")
	}
	fmt.Println("2020th number:", findNthNumber(ns, 2020))
	fmt.Println("3000000th number:", findNthNumber(ns, 30000000))
}

func findNthNumber(ns []int, n int) int {
	seen := make(map[int]int)
	for i := 0; i < len(ns)-1; i++ {
		seen[ns[i]] = i
	}
	N := ns[len(ns)-1]

	for i := len(ns); i < n; i++ {
		s, ok := seen[N]
		seen[N] = i - 1
		if ok {
			N = i - 1 - s
		} else {
			N = 0
		}
	}
	return N
}
