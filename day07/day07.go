package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gonih.org/stack"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(buf)), "\n")
	root, err := Collect(lines)
	if err != nil {
		log.Fatal(err)
	}
	du := DiskUsages(root)
	for k, v := range du {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Printf("Total disk usage of small-ish dirs: %d\n", DiskUsageSmall(du))
	d, n := FreeSpace(du, 70000000, 30000000)
	fmt.Printf("Freeing %q frees up %d bytes\n", d, n)
}

type File struct {
	Name     string
	Size     int
	Children map[string]*File
}

func Collect(lines []string) (*File, error) {
	var dirs stack.Stack[*File]
	cwd := &File{
		Children: make(map[string]*File),
	}
	root := cwd
	dirs.Push(cwd)
	for i, l := range lines {
		sp := strings.Split(l, " ")
		if sp[0] != "$" {
			if len(sp) != 2 {
				return nil, fmt.Errorf("%d: invalid line %q", i+1, l)
			}
			if sp[0] == "dir" {
				cwd.Children[sp[1]] = &File{
					Name:     sp[1],
					Children: make(map[string]*File),
				}
				continue
			}
			n, err := strconv.Atoi(sp[0])
			if err != nil {
				return nil, fmt.Errorf("%d: invalid line %q: %w", i+1, l, err)
			}
			cwd.Children[sp[1]] = &File{Size: n}
			continue
		}
		if sp[1] == "ls" {
			if len(sp) != 2 {
				return nil, fmt.Errorf("%d: invalid line %q", i+1, l)
			}
			continue
		}
		if sp[1] != "cd" || len(sp) != 3 {
			return nil, fmt.Errorf("%d: invalid line %q", i+1, l)
		}
		switch d := sp[2]; d {
		case "..":
			dirs.Pop()
			cwd = dirs.Top()
		case "/":
			cwd = dirs[0]
			dirs = dirs[:1]
		default:
			cwd = cwd.Children[d]
			dirs.Push(cwd)
		}
	}
	_ = root
	return dirs[0], nil
}

func DiskUsages(root *File) map[string]int {
	var (
		walk func(*File) int
		cwd  stack.Stack[string]
		out  = make(map[string]int)
		pwd  = func() string {
			d := strings.Join(cwd, "/")
			if d == "" {
				d = "/"
			}
			return d
		}
	)
	walk = func(f *File) int {
		cwd.Push(f.Name)
		defer cwd.Pop()
		if f.Children == nil {
			return f.Size
		}
		var n int
		for _, c := range f.Children {
			n += walk(c)
		}
		out[pwd()] = n
		return n
	}
	walk(root)
	return out
}

func DiskUsageSmall(du map[string]int) int {
	var total int
	for _, v := range du {
		if v < 100000 {
			total += v
		}
	}
	return total
}

func FreeSpace(du map[string]int, total, need int) (string, int) {
	free := total - du["/"]
	bestD := "/"
	bestN := du["/"]
	for k, v := range du {
		if free+v >= need && v < bestN {
			bestD, bestN = k, v
		}
	}
	return bestD, bestN
}
