package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in := string(bytes.TrimSpace(buf))
	out := react(in)

	fmt.Println(len(in), len(out))
}

func react(p string) string {
	var in, out []byte
	in = []byte(p)
	for {
		out = out[:0]
		for i := 1; i <= len(in); i++ {
			if i < len(in) && reduce(in[i-1], in[i]) {
				i++
				continue
			}
			out = append(out, in[i-1])
		}
		if bytes.Compare(in, out) == 0 {
			return string(in)
		}
		in, out = out, in
	}
}

func reduce(a, b byte) bool {
	if (a >= 'a') == (b >= 'a') {
		return false
	}
	return element(a) == element(b)
}

func element(b byte) byte {
	if b >= 'a' {
		return b
	}
	return b - 'A' + 'a'
}
