package main

import (
	"fmt"
	"image"
	"io"
	"log"
)

func main() {
	claims := readClaims()

outerLoop:
	for i := 0; i < len(claims); i++ {
		for j := 0; j < len(claims); j++ {
			if i == j {
				continue
			}
			if intersect(claims[i], claims[j]) {
				continue outerLoop
			}
		}
		fmt.Println(claims[i].ID)
	}
}

type claim struct {
	ID     int
	Left   int
	Top    int
	Width  int
	Height int
}

func intersect(c1, c2 claim) bool {
	r1 := image.Rect(c1.Left, c1.Top, c1.Left+c1.Width, c1.Top+c1.Height)
	r2 := image.Rect(c2.Left, c2.Top, c2.Left+c2.Width, c2.Top+c2.Height)
	return r1.Overlaps(r2)
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
