package main

import (
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/google/renameio/v2"
)

var update = flag.Bool("update", false, "update dot files of testcases")

type TestResult struct {
	Value int
	Error string
}

type TestCase struct {
	name  string
	Want  TestResult
	Input Network
}

func ReadTestCases(t *testing.T) []TestCase {
	t.Helper()
	testdata := filepath.Join("testdata", filepath.FromSlash(t.Name()))
	d, err := os.ReadDir(testdata)
	if err != nil {
		t.Fatal(err)
	}
	var tcs []TestCase
	for _, e := range d {
		name, ok := strings.CutSuffix(e.Name(), ".txt")
		if !ok {
			continue
		}
		b, err := os.ReadFile(filepath.Join(testdata, e.Name()))
		if err != nil {
			t.Error(err)
			continue
		}
		tc, err := parse.Struct[TestCase](
			split.SplitN("\n", 2),
			func(in string) (TestResult, error) {
				v, err := strconv.Atoi(in)
				if err != nil {
					return TestResult{Error: in}, nil
				}
				return TestResult{Value: v}, nil
			},
			Parse,
		)(string(b))
		if err != nil {
			t.Errorf("parsing %s: %v", filepath.Join(testdata, e.Name()), err)
			continue
		}
		if *update {
			if _, err := validate(tc.Input); err == nil {
				if name == "input" {
					t.Log("Skipping input.dot due to size")
				} else if err = writeDot(filepath.Join(testdata, name+".dot"), tc.Input); err != nil {
					t.Error(err)
				}
			}
		}
		tc.name = name
		tcs = append(tcs, tc)
	}
	return tcs
}

func writeDot(fname string, net Network) error {
	f, err := renameio.NewPendingFile(fname, renameio.WithExistingPermissions())
	if err != nil {
		return err
	}
	defer f.Cleanup()
	if err := WriteDot(f, net); err != nil {
		return err
	}
	return f.CloseAtomicallyReplace()
}

func TestPart1(t *testing.T) {
	for _, tc := range ReadTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Part1(tc.Input)
			if err != nil {
				if tc.Want.Error == "" {
					t.Errorf("Part1(…) = %v, %v, want %v, <nil>", got, err, tc.Want.Value)
				} else if err.Error() != tc.Want.Error {
					t.Errorf("Part1(…) = %v, %v, want _, %s", got, err, tc.Want.Error)
				}
				return
			}
			if tc.Want.Error != "" {
				t.Errorf("Part1(…) = %v, <nil>, want _, %s", got, tc.Want.Error)
			} else if got != tc.Want.Value {
				t.Errorf("Part1(…) = %v, want %v", got, tc.Want.Value)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	for _, tc := range ReadTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Part2(tc.Input)
			if err != nil {
				if tc.Want.Error == "" {
					t.Errorf("Part2(…) = %v, %v, want %v, <nil>", got, err, tc.Want.Value)
				} else if err.Error() != tc.Want.Error {
					t.Errorf("Part2(…) = %v, %v, want _, %s", got, err, tc.Want.Error)
				}
				return
			}
			if tc.Want.Error != "" {
				t.Errorf("Part2(…) = %v, <nil>, want _, %s", got, tc.Want.Error)
			} else if got != tc.Want.Value {
				t.Errorf("Part2(…) = %v, want %v", got, tc.Want.Value)
			}
		})
	}
}
