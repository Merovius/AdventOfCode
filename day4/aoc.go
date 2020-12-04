package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	pp, err := ParsePassports(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var (
		valids         int
		strictlyValids int
	)
	for _, p := range pp {
		if p.IsValid() {
			valids++
		}
		if p.IsStrictlyValid() {
			strictlyValids++
		}
	}
	fmt.Printf("%d valid passports\n", valids)
	fmt.Printf("%d strictly valid passports\n", strictlyValids)
}

type passport map[string]string

func ParsePassports(r io.Reader) ([]passport, error) {
	s := bufio.NewScanner(r)
	pp := make(passport)
	var out []passport
	for s.Scan() {
		if s.Text() == "" {
			out = append(out, pp)
			pp = make(passport)
			continue
		}
		fields := strings.Split(s.Text(), " ")
		for _, f := range fields {
			kv := strings.SplitN(f, ":", 2)
			if len(kv) != 2 {
				return nil, errors.New("missing : in key/value")
			}
			if _, ok := pp[kv[0]]; ok {
				return nil, fmt.Errorf("duplicate key %q in passport", kv[0])
			}
			pp[kv[0]] = kv[1]
		}
	}
	out = append(out, pp)
	return out, s.Err()
}

func (p passport) IsValid() bool {
	hasAll := func(ks ...string) bool {
		for _, k := range ks {
			if _, ok := p[k]; !ok {
				return false
			}
		}
		return true
	}
	return hasAll("byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid")
}

func (p passport) IsStrictlyValid() bool {
	return p.byr() && p.iyr() && p.eyr() && p.hgt() && p.hcl() && p.ecl() && p.pid()
}

func (p passport) byr() bool {
	v, err := strconv.Atoi(p["byr"])
	if err != nil {
		return false
	}
	return v >= 1920 && v <= 2002
}

func (p passport) iyr() bool {
	v, err := strconv.Atoi(p["iyr"])
	if err != nil {
		return false
	}
	return v >= 2010 && v <= 2020
}

func (p passport) eyr() bool {
	v, err := strconv.Atoi(p["eyr"])
	if err != nil {
		return false
	}
	return v >= 2020 && v <= 2030
}

func (p passport) hgt() bool {
	s := p["hgt"]
	var cm bool
	if cm = strings.HasSuffix(s, "cm"); !cm {
		if !strings.HasSuffix(s, "in") {
			return false
		}
	}
	s = s[:len(s)-2]

	v, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return (cm && v >= 150 && v <= 193) || (!cm && v >= 59 && v <= 76)
}

func (p passport) hcl() bool {
	s := p["hcl"]
	if !strings.HasPrefix(s, "#") {
		return false
	}
	s = s[1:]
	_, err := hex.DecodeString(s)
	if err != nil {
		return false
	}
	return strings.ToLower(s) == s
}

var validEcl = map[string]bool{
	"amb": true,
	"blu": true,
	"brn": true,
	"gry": true,
	"grn": true,
	"hzl": true,
	"oth": true,
}

func (p passport) ecl() bool {
	return validEcl[p["ecl"]]
}

func (p passport) pid() bool {
	s := p["pid"]
	_, err := strconv.Atoi(s)
	return err == nil && len(s) == 9
}
