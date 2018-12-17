package main

import (
	"fmt"
	"io"
	"log"
)

func main() {
	edges, N := readEdges()
	in := make(map[int][]int)
	out := make(map[int][]int)
	for _, e := range edges {
		in[e[1]] = append(in[e[1]], e[0])
		out[e[0]] = append(out[e[0]], e[1])
	}
	done := make(map[int]bool)
	for len(done) < N {
		var i int
		for ; i < N; i++ {
			if !done[i] && len(in[i]) == 0 {
				break
			}
		}
		if i == N {
			log.Fatal("No nodes without dependencies")
		}
		fmt.Print(string(byte(i) + 'A'))
		done[i] = true
		for _, j := range out[i] {
			l := in[j]
			for k := range l {
				if l[k] == i {
					in[j] = append(l[:k], l[k+1:]...)
					break
				}
			}
		}
	}
	fmt.Println()
}

type edge [2]int

func readEdges() (out []edge, N int) {
	for {
		var from, to string
		_, err := fmt.Scanf("Step %s must be finished before step %s can begin.\n", &from, &to)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		f, t := int(from[0]-'A'), int(to[0]-'A')
		if f > N {
			N = f
		}
		if t > N {
			N = t
		}
		out = append(out, edge{f, t})
	}
	return out, N + 1
}
