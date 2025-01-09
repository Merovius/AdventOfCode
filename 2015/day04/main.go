package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	buf := append([]byte(nil), in...)
	for i := 0; ; i++ {
		buf = strconv.AppendInt(buf[:len(in)], int64(i), 10)
		h := md5.Sum(buf)
		if h[0] == 0 && h[1] == 0 && (h[2]>>4) == 0 {
			return i
		}
	}
}

func Part2(in string) int {
	buf := append([]byte(nil), in...)
	for i := 0; ; i++ {
		buf = strconv.AppendInt(buf[:len(in)], int64(i), 10)
		h := md5.Sum(buf)
		if h[0] == 0 && h[1] == 0 && h[2] == 0 {
			return i
		}
	}
}
