package main

import (
	"fmt"
	"io"
	"log"
)

func main() {
	var total int
	for {
		var δ int
		_, err := fmt.Scan(&δ)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		total += δ
	}
	fmt.Println(total)
}
