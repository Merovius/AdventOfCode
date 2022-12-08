package main

import (
	"container/heap"
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

	const (
		NWorkers = 5
		FixTime  = 60
	)

	var (
		h = new(intHeap)
		T int
	)
	done := make(map[int]bool)
	started := make(map[int]bool)
	for len(done) < N {
		for h.Len() < NWorkers {
			var i int
			for ; i < N; i++ {
				if !started[i] && len(in[i]) == 0 {
					break
				}
			}
			if i == N {
				break
			}
			started[i] = true
			heap.Push(h, heapItem{T + FixTime + i + 1, i})
		}
		if h.Len() == 0 {
			log.Fatal("All elves are idle, but still work to do")
		}
		doneWork := []heapItem{heap.Pop(h).(heapItem)}
		for h.Len() > 0 {
			tmp := heap.Pop(h).(heapItem)
			if tmp.t > doneWork[0].t {
				heap.Push(h, tmp)
				break
			}
			doneWork = append(doneWork, tmp)
		}
		T = doneWork[0].t
		for _, it := range doneWork {
			i := it.step
			fmt.Printf("T = %d, Done = %s\n", T, string(byte(i)+'A'))
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
	}
}

type heapItem struct {
	t    int
	step int
}

type intHeap []heapItem

func (h *intHeap) Len() int {
	return len(*h)
}

func (h *intHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *intHeap) Less(i, j int) bool {
	hi, hj := (*h)[i], (*h)[j]
	if hi.t < hj.t {
		return true
	}
	if hj.t < hi.t {
		return false
	}
	return hi.step < hj.step
}

func (h *intHeap) Push(v interface{}) {
	*h = append(*h, v.(heapItem))
}

func (h *intHeap) Pop() (v interface{}) {
	n := h.Len() - 1
	*h, v = (*h)[:n], (*h)[n]
	return v
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
