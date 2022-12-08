package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

func main() {
	flag.Parse()
	serial, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	var grid [300][300]int
	for x := 0; x < 300; x++ {
		rack := x + 10
		for y := 0; y < 300; y++ {
			pl := (rack*y + serial) % 1000
			pl = (pl * rack) % 1000
			pl = pl / 100
			grid[x][y] = pl - 5
		}
	}

	var mx, my, mv int
	for x := 0; x < 298; x++ {
		for y := 0; y < 298; y++ {
			pl := 0
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					pl += grid[x+i][y+j]
				}
			}
			if pl > mv {
				mx, my, mv = x, y, pl
			}
		}
	}
	fmt.Printf("%d,%d -> %d\n", mx, my, mv)
}
