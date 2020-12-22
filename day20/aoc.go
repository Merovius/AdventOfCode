package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	tiles, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if err := ReorderTiles(tiles); err != nil {
		log.Fatal(err)
	}
	N := sqrt(len(tiles))
	M := len(tiles[0].Tile)
	for i := 0; i < N; i++ {
		hdr := make([]string, N)
		for j := 0; j < N; j++ {
			hdr[j] = pad(tiles[i*N+j].Index, M)
		}
		fmt.Println(strings.Join(hdr, " "))
		for r := 0; r < len(tiles[0].Tile); r++ {
			for j := 0; j < N; j++ {
				fmt.Print(tiles[i*N+j].Tile[r])
				fmt.Print(" ")
			}
			fmt.Println()
		}
		fmt.Println()
	}

	c := []int{
		tiles[0].Index,
		tiles[N-1].Index,
		tiles[N*N-N].Index,
		tiles[N*N-1].Index,
	}
	fmt.Println("Corners:", c)
	fmt.Println("Product of corners:", c[0]*c[1]*c[2]*c[3])

	m := AssembleMap(tiles)
	for _, r := range m.Tile {
		fmt.Println(r)
	}
	ms := FindMapOrientations(m)
	if len(ms) == 0 {
		log.Fatal("Could not find seamonsters")
	}
	if len(ms) > 1 {
		log.Fatal("More than one orientation has seamonsters")
	}
	m = ms[0]
	fmt.Println("Correct orientation:")
	for _, r := range m.Tile {
		fmt.Println(r)
	}
	sm := CountSeamonsters(m)
	hashes := CountHashes(m)
	fmt.Println("Number of non-seamonster hashes:", hashes-sm*15)
}

func FindMapOrientations(m Tile) []Tile {
	var ms []Tile
	for i := 0; i < 4; i++ {
		if FindSeamonsters(m) {
			ms = append(ms, m)
		}
		m = rotateTile(m, 1)
	}
	m = flipTile(m)
	for i := 0; i < 4; i++ {
		if FindSeamonsters(m) {
			ms = append(ms, m)
		}
		m = rotateTile(m, 1)
	}
	return ms
}

func pad(n, w int) string {
	s := strconv.Itoa(n)
	if len(s) >= w {
		return s
	}
	return strings.Repeat(" ", w-len(s)) + s
}

type Tile struct {
	Index int
	Tile  []string
}

func ReadInput(r io.Reader) ([]Tile, error) {
	var (
		tiles   []Tile
		current Tile
	)
	s := bufio.NewScanner(r)
	for s.Scan() {
		if current.Index == 0 {
			if _, err := fmt.Sscanf(s.Text(), "Tile %d:", &current.Index); err != nil {
				return nil, fmt.Errorf("parsing %q: %w", s.Text(), err)
			}
			continue
		}
		l := s.Text()
		if l == "" {
			tiles = append(tiles, current)
			current = Tile{}
			continue
		}
		current.Tile = append(current.Tile, l)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(tiles) == 0 {
		return nil, errors.New("no tiles in input")
	}
	if len(tiles[0].Tile) == 0 || len(tiles[0].Tile[0]) == 0 {
		return nil, errors.New("tiles are empty")
	}
	for _, t := range tiles {
		if len(t.Tile) != len(tiles[0].Tile) {
			return nil, fmt.Errorf("tile %d has different number of rows", t.Index)
		}
		if len(t.Tile) != len(t.Tile[0]) {
			return nil, fmt.Errorf("tile %d is not square", t.Index)
		}
		for _, r := range t.Tile {
			if len(r) != len(tiles[0].Tile[0]) {
				return nil, errors.New("tiles have different number of columns")
			}
		}
	}
	if _, ok := sqrt2(len(tiles)); !ok {
		return nil, fmt.Errorf("number of tiles (%d) is not a square", len(tiles))
	}

	return tiles, nil
}

func ReorderTiles(ts []Tile) error {
	N := sqrt(len(ts))

	type edge struct {
		tile   int
		result Tile
	}
	edges := make(map[string][]edge)
	for i, t := range ts {
		for j := 0; j < 4; j++ {
			edges[t.Tile[0]] = append(edges[t.Tile[0]], edge{
				tile:   i,
				result: t,
			})
			t = rotateTile(t, 1)
		}
		t = flipTile(t)
		for j := 0; j < 4; j++ {
			edges[t.Tile[0]] = append(edges[t.Tile[0]], edge{
				tile:   i,
				result: t,
			})
			t = rotateTile(t, 1)
		}
	}

	var (
		n           int
		uniqueEdges = make([]int, len(ts))
	)
	for k, es := range edges {
		if len(es) > 2 {
			return fmt.Errorf("edge %q appears more than twice", k)
		}
		if len(es) == 1 {
			n++
			uniqueEdges[es[0].tile]++
		}
	}
	if n != 8*N {
		return fmt.Errorf("%d edges are unique", n)
	}
	for i, ue := range uniqueEdges {
		fmt.Printf("%d: %v\n", i, ue)
	}

	// Find corner tile
	var corner int
	for i, n := range uniqueEdges {
		// 4 unique edges, as they appear double (once flipped)
		if n == 4 {
			corner = i
			break
		}
	}
	fmt.Printf("Selected corner: %d (id %d)\n", corner, ts[corner].Index)

	// Find correct rotation of corner tile. ue holds the unique edges of the
	// corner tile, in consecutive (clockwise) order.
	var ue []int
	for i := 0; i < 4; i++ {
		k := extractEdge(ts[corner], i)
		fmt.Printf("%q: %v\n", k, len(edges[k]))
		if len(edges[k]) == 1 {
			ue = append(ue, i)
		}
	}
	if ue[0] == 0 && ue[1] == 3 {
		ue[0], ue[1] = 3, 0
	}

	removeTile := func(idx int) {
		n := 0
		for k, es := range edges {
			es2 := es[:0]
			for _, e := range es {
				if e.tile != idx {
					es2 = append(es2, e)
				}
			}
			n += len(es) - len(es2)
			edges[k] = es2
		}
		log.Printf("removed %d edges for tile %d", n, idx)
	}

	printTile := func(t Tile) {
		for _, t := range t.Tile {
			fmt.Println(t)
		}
		fmt.Println()
	}

	// Assemble puzzle. Start with the corner tile, with one unique edge left
	// and one on top.
	ordered := []Tile{rotateTile(ts[corner], ue[1])}
	printTile(rotateTile(ts[corner], ue[1]))
	removeTile(corner)
	for len(ordered) < len(ts) {
		if len(ordered)%N != 0 {
			last := ordered[len(ordered)-1]
			// The connecting piece must have the reverse edge orientation.
			k := reverse(extractEdge(last, 1))
			if len(edges[k]) != 1 {
				return fmt.Errorf("%d edges with %q remaining", len(edges[k]), k)
			}
			e := edges[k][0]
			t := rotateTile(e.result, 1)
			printTile(t)
			ordered = append(ordered, t)
			removeTile(e.tile)
			continue
		}
		last := ordered[len(ordered)-N]
		// We connect to the bottom edge now.
		k := reverse(extractEdge(last, 2))
		if len(edges[k]) != 1 {
			return fmt.Errorf("%d edges with %q remaining", len(edges[k]), k)
		}
		e := edges[k][0]
		// the needed edge is already at the top - no rotation.
		printTile(e.result)
		ordered = append(ordered, e.result)
		removeTile(e.tile)
		continue
	}
	copy(ts, ordered)
	return nil
}

func AssembleMap(ts []Tile) Tile {
	var m Tile
	N := sqrt(len(ts))
	var w strings.Builder
	for i := 0; i < N; i++ {
		for r := 1; r < len(ts[0].Tile)-1; r++ {
			for j := 0; j < N; j++ {
				s := ts[i*N+j].Tile[r]
				w.WriteString(s[1 : len(s)-1])
			}
			m.Tile = append(m.Tile, w.String())
			w.Reset()
		}
	}
	return m
}

// dedup removes duplicates from s, as well as r.
func dedup(s []int, r int) []int {
	sort.Ints(s)
	o := s[:0]
	for i := 0; i < len(s); i++ {
		if s[i] == r {
			continue
		}
		if len(o) > 0 && o[len(o)-1] == s[i] {
			continue
		}
		o = append(o, s[i])
	}
	return o
}

func reverse(s string) string {
	var w strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		w.WriteByte(s[i])
	}
	return w.String()
}

func extractEdge(t Tile, e int) string {
	switch e {
	case 0:
		return t.Tile[0]
	case 1:
		var b []byte
		for j := 0; j < len(t.Tile); j++ {
			b = append(b, t.Tile[j][len(t.Tile[j])-1])
		}
		return string(b)
	case 2:
		return reverse(t.Tile[len(t.Tile)-1])
	case 3:
		var b []byte
		for j := 0; j < len(t.Tile); j++ {
			b = append(b, t.Tile[j][0])
		}
		return reverse(string(b))
	default:
		panic(fmt.Errorf("invalid edge index %d", e))
	}
}

// rotateTile rotates t, such that edge e is at the top.
func rotateTile(t Tile, e int) Tile {
	o := Tile{
		Index: t.Index,
	}
	if e == 0 {
		o.Tile = append(o.Tile, t.Tile...)
		return o
	}
	if e == 2 {
		o.Tile = make([]string, len(t.Tile))
		for i := 0; i < len(t.Tile); i++ {
			o.Tile[i] = reverse(t.Tile[len(t.Tile)-1-i])
		}
		return o
	}
	for j := len(t.Tile[0]) - 1; j >= 0; j-- {
		var b []byte
		for i := 0; i < len(t.Tile); i++ {
			b = append(b, t.Tile[i][j])
		}
		o.Tile = append(o.Tile, string(b))
	}
	if e == 1 {
		return o
	}
	return rotateTile(o, 2)
}

func flipTile(t Tile) Tile {
	o := Tile{
		Index: t.Index,
	}
	o.Tile = make([]string, len(t.Tile))
	for i := 0; i < len(t.Tile); i++ {
		o.Tile[i] = t.Tile[len(t.Tile)-1-i]
	}
	return o
}

func (t Tile) String() string {
	return strings.Join(t.Tile, "\n")
}

func sqrt(n int) int {
	m, ok := sqrt2(n)
	if !ok {
		panic(fmt.Errorf("%d is not a square", n))
	}
	return m
}

func sqrt2(n int) (m int, ok bool) {
	m = int(math.Round(math.Sqrt(float64(n))))
	return m, m*m == n
}

var seamonster []*regexp.Regexp

func init() {
	seamonster = append(seamonster, regexp.MustCompile(`..................#.`))
	seamonster = append(seamonster, regexp.MustCompile(`#....##....##....###`))
	seamonster = append(seamonster, regexp.MustCompile(`.#..#..#..#..#..#...`))
}

func FindSeamonsters(m Tile) bool {
	for i := 1; i < len(m.Tile)-1; i++ {
		loc := seamonster[1].FindStringIndex(m.Tile[i])
		if len(loc) == 0 {
			continue
		}
		if !seamonster[0].MatchString(m.Tile[i-1][loc[0]:loc[1]]) {
			continue
		}
		if !seamonster[2].MatchString(m.Tile[i+1][loc[0]:loc[1]]) {
			continue
		}
		return true
	}
	return false
}

// CountSeamonsters counts the seamonsters in m, assuming they don't overlap
// (which might not be true in general, but it is in our input :)
func CountSeamonsters(m Tile) int {
	var pos [][2]int

	for i := 1; i < len(m.Tile)-1; i++ {
		for j := 0; j < len(m.Tile[i]); {
			loc := seamonster[1].FindStringIndex(m.Tile[i][j:])
			if loc == nil {
				break
			}
			pos = append(pos, [2]int{i, loc[0] + j})
			j += loc[0] + 1
		}
	}
	var n int
	for _, p := range pos {
		if !seamonster[0].MatchString(m.Tile[p[0]-1][p[1] : p[1]+20]) {
			continue
		}
		if !seamonster[0].MatchString(m.Tile[p[0]+1][p[1] : p[1]+20]) {
			continue
		}
		fmt.Println("Confirmed sea monster at", p)
		n++
	}
	return n
}

func CountHashes(m Tile) int {
	var n int
	for _, r := range m.Tile {
		n += strings.Count(r, "#")
	}
	return n
}
