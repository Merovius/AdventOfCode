package main

import (
	"fmt"
	"io"
	"log"
)

func main() {
	claims := readClaims()
	w, h := totalDim(claims)
	mark := make([]int, w*h)
	for _, c := range claims {
		for i := c.Top; i < c.Top+c.Height; i++ {
			for j := c.Left; j < c.Left+c.Width; j++ {
				mark[i*w+j] += 1
			}
		}
	}
	var total int
	for _, m := range mark {
		if m > 1 {
			total++
		}
	}
	fmt.Println(total)
}

type claim struct {
	ID     int
	Left   int
	Top    int
	Width  int
	Height int
}

func readClaims() []claim {
	var out []claim
	for {
		var c claim
		_, err := fmt.Scanf("#%d @ %d,%d: %dx%d", &c.ID, &c.Left, &c.Top, &c.Width, &c.Height)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, c)
	}
	return out
}

func totalDim(claims []claim) (w, h int) {
	for _, c := range claims {
		if c.Left+c.Width > w {
			w = c.Left + c.Width
		}
		if c.Top+c.Height > h {
			h = c.Top + c.Height
		}
	}
	return w, h
}
