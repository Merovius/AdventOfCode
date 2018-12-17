package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"
)

func main() {
	entries := scanEntries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Time.Before(entries[j].Time) })
	maxG := maxGuard(entries)
	data := make([]hist, maxG+1)
	var (
		g      int
		asleep bool
		minute int
	)
	for _, e := range entries {
		if e.GuardID != 0 && e.GuardID != g {
			g = e.GuardID
			continue
		}
		for ; minute != e.Time.Minute(); minute = (minute + 1) % 60 {
			if asleep {
				data[g][minute]++
			}
		}
		asleep = e.Sleep
	}

	var (
		bestGuard  int
		bestMinute int
	)
	for g, h := range data {
		for m, v := range h {
			if v > data[bestGuard][bestMinute] {
				fmt.Printf("Guard #%d is %d times asleep during minute %d\n", g, v, m)
				bestGuard, bestMinute = g, m
			}
		}
	}
	fmt.Printf("Guard %d during minute %d: %d\n", bestGuard, bestMinute, bestGuard*bestMinute)
}

type hist [60]int

func (h hist) total() int {
	var total int
	for _, v := range h {
		total += v
	}
	return total
}

type entry struct {
	Time    time.Time
	GuardID int
	Sleep   bool
}

func scanEntries() []entry {
	var out []entry
	re := regexp.MustCompile(`\[(.*)\] (.*)`)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := re.FindStringSubmatch(s.Text())
		if len(m) == 0 {
			log.Fatalf("%q didn't match", s.Text())
		}
		var e entry
		t, err := time.Parse("2006-01-02 15:04", m[1])
		if err != nil {
			log.Fatalf("Couldn't parse %q: %v", m[1], err)
		}
		e.Time = t
		switch m[2] {
		case "falls asleep":
			e.Sleep = true
		case "wakes up":
		default:
			_, err = fmt.Sscanf(m[2], "Guard #%d begins shift", &e.GuardID)
			if err != nil {
				log.Fatalf("Unrecognized entry %q", m[2])
			}
		}
		out = append(out, e)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return out
}

func maxGuard(es []entry) int {
	var g int
	for _, e := range es {
		if e.GuardID > g {
			g = e.GuardID
		}
	}
	return g
}
