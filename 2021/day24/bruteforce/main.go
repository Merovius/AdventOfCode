package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	log.SetFlags(0)
	out, err := os.Create("serials.txt")
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	blocks := make(chan Block)
	serials := make(chan string)
	done := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go Worker(wg, blocks, serials)
	}
	go WriteSerials(out, serials, done)
	var nBlocks uint64
	go func() {
		tick := time.Tick(10 * time.Second)
		last := uint64(0)
		for {
			select {
			case <-tick:
				n := atomic.LoadUint64(&nBlocks)
				log.Printf("%d blocks generated (%f blocks/s, %f inputs/s)", n, float64(n-last)/10, float64(n-last)/10*math.Pow(9, blockSize))
				last = n
			case <-done:
				return
			}
		}
	}()

	var b Block
	for i := range b {
		b[i] = 1
	}
	for {
		blocks <- b
		atomic.AddUint64(&nBlocks, 1)
		if inc(b[:]) {
			break
		}
	}
	close(blocks)
	wg.Wait()
	close(serials)
	<-done
}

const (
	blockSize = 6
	nInput    = 14
)

type Block [nInput - blockSize]int

func Worker(wg *sync.WaitGroup, ch <-chan Block, serials chan<- string) {
	defer wg.Done()

	buf := new(strings.Builder)

	in := make([]int, nInput)
	for i := range in {
		in[i] = 1
	}
	for b := range ch {
		copy(in, b[:])
		for {
			if eval(in) == 0 {
				for _, v := range in {
					buf.WriteByte(byte(v) + '0')
				}
				serials <- buf.String()
				buf.Reset()
			}
			if inc(in[nInput-blockSize:]) {
				break
			}
		}
	}
}

// inc increments the block by b, returning true if a wrap-around occured.
func inc(b []int) bool {
	for i := range b {
		b[i] += 1
		if b[i] < 10 {
			return false
		}
		b[i] = 1
	}
	return true
}

func WriteSerials(w io.Writer, serials <-chan string, done chan<- struct{}) {
	bw := bufio.NewWriter(w)
	for s := range serials {
		if _, err := bw.WriteString(s); err != nil {
			log.Fatal(err)
		}
		if err := bw.WriteByte('\n'); err != nil {
			log.Fatal(err)
		}
	}
	if err := bw.Flush(); err != nil {
		log.Fatal(err)
	}
	close(done)
}
