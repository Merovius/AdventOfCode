package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
	"strings"
)

func main() {
	v1 := make(map[uint64]uint64)
	f1 := func(addr, m0, m1, mx, v uint64) {
		v1[addr] = (v | m1) &^ m0
	}
	v2 := make(map[uint64]uint64)
	f2 := func(addr, m0, m1, mx, v uint64) {
		addr |= m1
		if mx == 0 {
			v2[addr] = v
			return
		}
		rangeMask(mx, func(m uint64) {
			v2[addr^m] = v
		})
	}
	if err := Exec(os.Stdin, f1, f2); err != nil {
		log.Fatal(err)
	}
	var t1, t2 uint64
	for _, v := range v1 {
		t1 += v
	}
	for _, v := range v2 {
		t2 += v
	}
	fmt.Println("Total (V1):", t1)
	fmt.Println("Total (V2):", t2)
}

func Exec(r io.Reader, fs ...func(addr uint64, m0, m1, mx, v uint64)) error {
	var (
		m0, m1, mx uint64
	)
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if strings.HasPrefix(l, "mask = ") {
			m0, m1, mx = parseBitmask(strings.TrimPrefix(l, "mask = "))
			continue
		}
		var (
			addr uint64
			v    uint64
		)
		if _, err := fmt.Sscanf(l, "mem[%d] = %d", &addr, &v); err != nil {
			return fmt.Errorf("could not parse %q: %w", l, err)
		}
		for _, f := range fs {
			f(addr, m0, m1, mx, v)
		}
	}
	return s.Err()
}

func parseBitmask(s string) (m0, m1, mx uint64) {
	for i := 0; i < len(s); i++ {
		switch s[len(s)-1-i] {
		case '0':
			m0 |= 1 << i
		case '1':
			m1 |= 1 << i
		case 'X':
			mx |= 1 << i
		}
	}
	return m0, m1, mx
}

func rangeMask(m uint64, f func(uint64)) {
	popcnt := bits.OnesCount64(m)
	for i := uint64(0); i < (uint64(1) << popcnt); i++ {
		var (
			v uint64
			s int
		)
		for n, j := m, i; n != 0; n >>= 1 {
			if n&1 == 1 {
				v |= (j & 1) << s
				j >>= 1
			}
			s++
		}
		f(v)
	}
}
