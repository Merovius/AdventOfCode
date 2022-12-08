package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Rows = 5
	Cols = 5
)

type Data struct {
	Nums   []int
	Boards []Board
}

type Board struct {
	Won   bool
	Score int

	vals   [Rows][Cols]int
	marked [Rows][Cols]bool
}

func (b *Board) Mark(v int) {
	if b.Won {
		return
	}
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if b.vals[r][c] != v {
				continue
			}
			b.marked[r][c] = true
			b.check(r, c)
		}
	}
}

func (b *Board) check(r, c int) {
	wr, wc := true, true
	for r := 0; r < Rows; r++ {
		if !b.marked[r][c] {
			wr = false
		}
	}
	for c := 0; c < Cols; c++ {
		if !b.marked[r][c] {
			wc = false
		}
	}
	if !wr && !wc {
		return
	}
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if !b.marked[r][c] {
				b.Score += b.vals[r][c]
			}
		}
	}
	b.Score *= b.vals[r][c]
	b.Won = true
}

func (b Board) String() string {
	buf := new(strings.Builder)
	for r := range b.vals {
		for c := range b.vals[r] {
			fmt.Fprintf(buf, "%2.1d", b.vals[r][c])
			if b.marked[r][c] {
				buf.WriteString("* ")
			} else {
				buf.WriteString("  ")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func readBoard(lines []string) (Board, error) {
	var b Board
	if lines[0] != "" {
		return Board{}, errors.New("want empty line")
	}
	for r := 0; r < Rows; r++ {
		fields := strings.Fields(lines[1+r])
		if len(fields) != 5 {
			return Board{}, errors.New("wrong number of fields")
		}
		for c, s := range fields {
			v, err := strconv.Atoi(s)
			if err != nil {
				return Board{}, err
			}
			b.vals[r][c] = v
		}
	}
	return b, nil
}

func readData(r io.Reader) (*Data, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	buf = bytes.TrimSpace(buf)
	lines := strings.Split(string(buf), "\n")
	if len(lines)%(Rows+1) != 1 {
		log.Println(len(lines))
		return nil, errors.New("wrong number of lines")
	}
	data := new(Data)
	for _, s := range strings.Split(lines[0], ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		data.Nums = append(data.Nums, v)
	}
	for i := 1; i < len(lines); i += 6 {
		b, err := readBoard(lines[i : i+6])
		if err != nil {
			return nil, fmt.Errorf("board %d: %v", len(data.Boards), err)
		}
		data.Boards = append(data.Boards, b)
	}
	return data, nil
}

func play(data *Data) (winners []Board) {
	for _, v := range data.Nums {
		for i := range data.Boards {
			if data.Boards[i].Won {
				continue
			}
			data.Boards[i].Mark(v)
			if data.Boards[i].Won {
				winners = append(winners, data.Boards[i])
			}
		}
	}
	return winners
}

func main() {
	log.SetFlags(0)

	data, err := readData(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	winners := play(data)
	fmt.Printf("First winner's score: %d\n", winners[0].Score)
	fmt.Printf("Last winner's score: %d\n", winners[len(winners)-1].Score)
}
