package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePassport(t *testing.T) {
	tcs := []struct {
		in   string
		want []passport
	}{
		{
			"ecl:gry",
			[]passport{{"ecl": "gry"}},
		},
		{
			"ecl:gry pid:860033327",
			[]passport{{"ecl": "gry", "pid": "860033327"}},
		},
		{
			"ecl:gry pid:860033327\nbyr:1937",
			[]passport{{"ecl": "gry", "pid": "860033327", "byr": "1937"}},
		},
		{
			"ecl:gry pid:860033327\nbyr:1937 iyr:2017",
			[]passport{{"ecl": "gry", "pid": "860033327", "byr": "1937", "iyr": "2017"}},
		},
		{
			"ecl:gry pid:860033327\nbyr:1937 iyr:2017\n\niyr:2013 ecl:amb\nhcl:#cfa07d",
			[]passport{
				{"ecl": "gry", "pid": "860033327", "byr": "1937", "iyr": "2017"},
				{"iyr": "2013", "ecl": "amb", "hcl": "#cfa07d"},
			},
		},
	}

	for _, tc := range tcs {
		got, err := ParsePassports(strings.NewReader(tc.in))
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("ParsePassport(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsValid(t *testing.T) {
	tcs := []struct {
		in   string
		want bool
	}{
		{"ecl:gry pid:860033327 eyr:2020 hcl:#fffffd\nbyr:1937 iyr:2017 cid:147 hgt:183cm", true},
		{"iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884\nhcl:#cfa07d byr:1929", false},
		{"hcl:#ae17e1 iyr:2013\neyr:2024\necl:brn pid:760753108 byr:1931\nhgt:179cm", true},
		{"hcl:#cfa07d eyr:2025 pid:166559648\niyr:2011 ecl:brn hgt:59in", false},
	}
	for _, tc := range tcs {
		pp, err := ParsePassports(strings.NewReader(tc.in))
		if err != nil {
			t.Errorf("ParsePassports(%q) = _, %v, want <nil>", tc.in, err)
			continue
		}
		if len(pp) != 1 {
			t.Errorf("len(ParsePassports(%q)) = %v, want 1", tc.in, len(pp))
			continue
		}
		if got := pp[0].IsValid(); got != tc.want {
			t.Errorf("ParsePassports(%q).IsValid() = %v, want %v", tc.in, got, tc.want)
		}
	}
}
