package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
	"github.com/Merovius/AdventOfCode/internal/math"
	"golang.org/x/exp/slices"
)

// Example command to render animation (scales up image by 4, using nearest
// neighbor interpolation, which is fine for integer scaling):
// 	ffmpeg -pattern_type glob -i 'frames/part1_*.png' -vf "scale=iw*4:-1" -sws_flags neighbor part1.mp4

func main() {
	animate := flag.String("animate", "", "output animation frames to folder")
	flag.Parse()
	var anim1, anim2 string
	if *animate != "" {
		anim1 = filepath.Join(*animate, "part1_%.5d.png")
		anim2 = filepath.Join(*animate, "part2_%.5d.png")
	}
	data, err := parse.Lines(
		parse.Slice(
			split.On(" -> "),
			parse.Array[[2]int](
				split.On(","),
				parse.Signed[int],
			),
		),
	).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	spawn, max := Remap(data)

	g := TracePaths(data, max)
	g2 := &image.Paletted{
		Pix:     slices.Clone(g.Pix),
		Stride:  g.Stride,
		Rect:    g.Rect,
		Palette: g.Palette,
	}
	vid := frameWriter{path: anim1}
	fmt.Println("Part 1:", SimulateFalling(g, spawn, false, vid.Write), "units of sand come to rest")
	vid = frameWriter{path: anim2}
	fmt.Println("Part 2:", SimulateFalling(g2, spawn, true, vid.Write), "units of sand come to rest")
}

const (
	Empty = iota
	Rock
	Sand
)

var palette = color.Palette{
	Empty: color.Black,
	Rock:  color.RGBA{0x80, 0x80, 0x80, 0xFF},
	Sand:  color.RGBA{0xFF, 0xFF, 0x00, 0xFF},
}

func Corners(paths [][][2]int) (minX, minY, maxX, maxY int) {
	minX, minY, maxX, maxY = math.MaxInt, math.MaxInt, math.MinInt, math.MinInt
	for _, path := range paths {
		for _, p := range path {
			minX, minY, maxX, maxY = math.Min(minX, p[0]), math.Min(minY, p[1]), math.Max(maxX, p[0]), math.Max(maxY, p[1])
		}
	}
	return minX, minY, maxX, maxY
}

func Remap(paths [][][2]int) (spawn, max image.Point) {
	spawnX, spawnY := 500, 0
	minX, minY, maxX, maxY := spawnX, spawnY, spawnX, spawnY
	for _, path := range paths {
		for _, p := range path {
			minX, minY, maxX, maxY = math.Min(minX, p[0]), math.Min(minY, p[1]), math.Max(maxX, p[0]), math.Max(maxY, p[1])
		}
	}
	minX, maxX = math.Min(minX, spawnX-maxY), math.Max(maxX, spawnX+maxY)
	const border = 2
	minX -= border
	maxX += border
	minY -= border
	maxY += border

	spawnX -= minX
	spawnY -= minY
	for i := range paths {
		for j := range paths[i] {
			paths[i][j][0] -= minX
			paths[i][j][1] -= minY
		}
	}
	maxX -= minX
	maxY -= minY
	return image.Pt(spawnX, spawnY), image.Pt(maxX, maxY)
}

func TracePaths(paths [][][2]int, max image.Point) *image.Paletted {
	img := image.NewPaletted(image.Rectangle{image.Point{}, max}, palette)
	for _, path := range paths {
		from := path[0]
		for _, to := range path[1:] {
			δx, δy := math.Sgn(to[0]-from[0]), math.Sgn(to[1]-from[1])
			for x, y := from[0], from[1]; x != to[0] || y != to[1]; x, y = x+δx, y+δy {
				img.SetColorIndex(x, y, Rock)
			}
			img.SetColorIndex(to[0], to[1], Rock)
			from = to
		}
	}
	return img
}

func SimulateFalling(img *image.Paletted, spawn image.Point, part2 bool, writeFrame func(image.Image) error) int {
	if writeFrame == nil {
		writeFrame = func(image.Image) error { return nil }
	}
	n := 0
	for {
		p := spawn
		if img.ColorIndexAt(p.X, p.Y) == Sand {
			return n
		}
	grainLoop:
		for {
			if p.Y == img.Rect.Max.Y-1 {
				if part2 {
					img.SetColorIndex(p.X, p.Y, Sand)
					writeFrame(img)
					n++
					break
				} else {
					return n
				}
			}
			for _, δ := range []image.Point{{0, 1}, {-1, 1}, {1, 1}} {
				if q := p.Add(δ); img.ColorIndexAt(q.X, q.Y) == Empty {
					p = p.Add(δ)
					continue grainLoop
				}
			}
			img.SetColorIndex(p.X, p.Y, Sand)
			writeFrame(img)
			n++
			break
		}
	}
}

type frameWriter struct {
	path  string
	frame int
}

func (w *frameWriter) Write(i image.Image) error {
	if w.path == "" {
		return nil
	}
	f, err := os.Create(fmt.Sprintf(w.path, w.frame))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, i); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	w.frame++
	return nil
}
