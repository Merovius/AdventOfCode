package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	r := repeat(os.Stdin)

	seen := make(map[int]bool)

	var total int
	for {
		var δ int
		_, err := fmt.Fscan(r, &δ)
		if err != nil {
			log.Fatal(err)
		}
		total += δ
		if seen[total] {
			break
		}
		seen[total] = true
	}
	fmt.Println(total)
}

func repeat(r io.Reader) io.Reader {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return &repeatReader{buf, bytes.NewReader(buf)}
}

type repeatReader struct {
	buf []byte
	r   *bytes.Reader
}

func (r *repeatReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == io.EOF {
		r.r.Reset(r.buf)
		return r.r.Read(p)
	}
	return n, err
}
