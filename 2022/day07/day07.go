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
	SetSizes(root)
	fmt.Printf("Total disk usage of small-ish dirs: %d\n", SumSmallDirs(root))
	fmt.Printf("Freeing up %d bytes\n", FreeSpace(root, 70000000, 30000000))
}

type File struct {
	Name     string
	Size     int
	Children map[string]*File
}

func Collect(lines []string) (*File, error) {
	var dirs stack.Stack[*File]
	dirs.Push(&File{
		Children: make(map[string]*File),
	})
	for i, l := range lines {
		sp := strings.Split(l, " ")
		if sp[0] != "$" {
			if len(sp) != 2 {
				return nil, fmt.Errorf("%d: invalid line %q", i+1, l)
			}
			if sp[0] == "dir" {
				dirs.Top().Children[sp[1]] = &File{
					Name:     sp[1],
					Children: make(map[string]*File),
				}
				continue
			}
			n, err := strconv.Atoi(sp[0])
			if err != nil {
				return nil, fmt.Errorf("%d: invalid line %q: %w", i+1, l, err)
			}
			dirs.Top().Children[sp[1]] = &File{Size: n}
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
		case "/":
			dirs = dirs[:1]
		default:
			dirs.Push(dirs.Top().Children[d])
		}
	}
	return dirs[0], nil
}

// Walk root in post-order.
func Walk(root *File, f func(*File)) {
	defer f(root)
	for _, c := range root.Children {
		Walk(c, f)
	}
}

func SetSizes(root *File) {
	Walk(root, func(f *File) {
		for _, c := range f.Children {
			f.Size += c.Size
		}
	})
}

func SumSmallDirs(root *File) (total int) {
	Walk(root, func(f *File) {
		if f.Children != nil && f.Size < 100000 {
			total += f.Size
		}
	})
	return total
}

func FreeSpace(root *File, total, need int) (freed int) {
	free := total - root.Size
	best := root.Size
	Walk(root, func(f *File) {
		if free+f.Size >= need && f.Size < best {
			best = f.Size
		}
	})
	return best
}
