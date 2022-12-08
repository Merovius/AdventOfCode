package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	in, err := parseInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	invalids := findInvalidValues(in)
	fmt.Printf("Invalid values: %v (Σ = %d)\n", invalids, sum(invalids))
	fmt.Printf("%d total tickets\n", len(in.nearby))
	discardInvalidTickets(&in)
	fmt.Printf("%d valid tickets\n", len(in.nearby))
	fo, err := fieldOrder(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Your ticket:")
	for i, f := range fo {
		fmt.Printf("\t%s: %d\n", f, in.yours[i])
	}
	type field struct {
		name  string
		value int
	}

	var depFields []field
	product := 1
	for i, f := range fo {
		if strings.HasPrefix(f, "departure") {
			depFields = append(depFields, field{f, in.yours[i]})
			product *= in.yours[i]
		}
	}
	for _, f := range depFields {
		fmt.Printf("%s: %v\n", f.name, f.value)
	}
	fmt.Printf("Product of departure fields: %v\n", product)
}

func fieldOrder(in input) ([]string, error) {
	all := (uint32(1) << len(in.yours)) - 1
	pos := make([]uint32, len(in.rules))
	for i := range pos {
		pos[i] = all
	}
	for _, t := range in.nearby {
		for i, v := range t {
		ruleLoop:
			for j, r := range in.rules {
				for _, rg := range r.ranges {
					if v >= rg.min && v <= rg.max {
						continue ruleLoop
					}
				}
				pos[j] &= ^(1 << i)
			}
		}
	}
	for i, p := range pos {
		fmt.Printf("%.20b %q\n", p, in.rules[i].field)
	}

	out := make([]string, len(in.rules))
outerLoop:
	for {
		for i, p := range pos {
			if bits.OnesCount32(p) == 1 {
				idx := bits.TrailingZeros32(p)
				if out[idx] != "" {
					return nil, fmt.Errorf("both %q and %q only fit in %d", out[idx], in.rules[i].field, idx)
				}
				out[idx] = in.rules[i].field
				for j, p := range pos {
					pos[j] = p &^ (1 << idx)
				}
				continue outerLoop
			}
		}
		for _, p := range pos {
			if p != 0 {
				return nil, errors.New("can't find unique valid field order")
			}
		}
		break
	}
	return out, nil
}

func findInvalidValues(in input) []int {
	var out []int
	for _, vs := range in.nearby {
	valLoop:
		for _, v := range vs {
			for _, r := range in.rules {
				for _, rg := range r.ranges {
					if v >= rg.min && v <= rg.max {
						continue valLoop
					}
				}
			}
			out = append(out, v)
		}
	}
	return out
}

func discardInvalidTickets(in *input) {
	valids := in.nearby[:0]
	for _, t := range in.nearby {
		valid := true
	ticketLoop:
		for _, v := range t {
			for _, r := range in.rules {
				for _, rg := range r.ranges {
					if v >= rg.min && v <= rg.max {
						continue ticketLoop
					}
				}
			}
			valid = false
			break
		}
		if valid {
			valids = append(valids, t)
		}
	}
	in.nearby = valids
}

func sum(vs []int) int {
	var total int
	for _, v := range vs {
		total += v
	}
	return total
}

type _range struct {
	min int
	max int
}

type rule struct {
	field  string
	ranges []_range
}

type input struct {
	rules  []rule
	yours  []int
	nearby [][]int
}

func parseInput(r io.Reader) (input, error) {
	var (
		in  input
		err error
	)
	s := bufio.NewScanner(r)
	if in.rules, err = parseRules(s); err != nil {
		return input{}, err
	}
	if in.yours, err = parseYourTicket(s); err != nil {
		return input{}, err
	}
	if in.nearby, err = parseNearby(s); err != nil {
		return input{}, err
	}
	N := len(in.yours)
	for _, t := range in.nearby {
		if len(t) != N {
			return input{}, errors.New("inconsistent ticket length")
		}
	}
	return in, s.Err()
}

func parseRules(s *bufio.Scanner) ([]rule, error) {
	var rs []rule
	for s.Scan() {
		l := s.Text()
		if l == "" {
			return rs, nil
		}
		i := strings.Index(l, ": ")
		if i < 0 {
			return nil, fmt.Errorf("can't parse rule %q", l)
		}
		var r rule
		r.field = l[:i]
		sp := strings.Split(l[i+2:], " or ")
		for _, s := range sp {
			sp := strings.Split(s, "-")
			if len(sp) != 2 {
				return nil, fmt.Errorf("can't parse interval %q", s)
			}
			min, err := strconv.Atoi(sp[0])
			if err != nil {
				return nil, fmt.Errorf("can't parse interval %q: %w", s, err)
			}
			max, err := strconv.Atoi(sp[1])
			if err != nil {
				return nil, fmt.Errorf("can't parse interval %q: %w", s, err)
			}
			r.ranges = append(r.ranges, _range{min, max})
		}
		rs = append(rs, r)
	}
	return rs, s.Err()
}

func parseYourTicket(s *bufio.Scanner) ([]int, error) {
	if !s.Scan() || s.Text() != "your ticket:" {
		return nil, errors.New(`can't find "your ticket:" line`)
	}
	if !s.Scan() {
		return nil, io.ErrUnexpectedEOF
	}
	out, err := splitNumbers(s.Text())
	if err != nil {
		return nil, err
	}
	if !s.Scan() || s.Text() != "" {
		return nil, errors.New("unexpected EOF or garbage after your ticket")
	}
	return out, nil
}

func parseNearby(s *bufio.Scanner) ([][]int, error) {
	if !s.Scan() || s.Text() != "nearby tickets:" {
		return nil, errors.New(`can't find "nearby tickets:" line`)
	}
	var out [][]int
	for s.Scan() {
		vs, err := splitNumbers(s.Text())
		if err != nil {
			return nil, err
		}
		out = append(out, vs)
	}
	return out, s.Err()
}

func splitNumbers(s string) ([]int, error) {
	sp := strings.Split(s, ",")
	var out []int
	for _, s := range sp {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
