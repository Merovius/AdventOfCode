package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Merovius/AdventOfCode/internal/aoc"
	"github.com/Merovius/AdventOfCode/internal/math"

	"github.com/google/subcommands"
	"golang.org/x/exp/slices"
)

type boardCmd struct {
	events bool
}

func (cmd *boardCmd) Name() string {
	return "board"
}

func (cmd *boardCmd) Synopsis() string {
	return "Fetch a private leaderboard"
}

func (cmd *boardCmd) Usage() string {
	return "Fetch a private leaderboard."
}

func (cmd *boardCmd) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.events, "events", false, "list a set of changes to the board and poll for updates")
}

func (cmd *boardCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	log.SetFlags(0)
	cfg, client := args[0].(*config), args[1].(*aoc.Client)

	var board string
	if f.NArg() == 0 {
		board = cfg.DefaultBoard
		if cfg.DefaultBoard == "" {
			log.Println("usage: board [-events] <board>")
			log.Println("neither board given, nor default_board configured")
			return subcommands.ExitUsageError
		}
	} else if f.NArg() == 1 {
		board = f.Arg(0)
	} else {
		log.Println("usage: board [-events] <board>")
		return subcommands.ExitUsageError
	}
	var boardID int
	for k, v := range cfg.Boards {
		if board == k || board == strconv.Itoa(v) {
			if boardID > 0 {
				log.Printf("board %q is ambiguous", board)
				return subcommands.ExitUsageError
			}
			boardID = v
		}
	}
	if boardID <= 0 {
		log.Printf("board %q not configured", board)
		return subcommands.ExitUsageError
	}
	if cmd.events {
		if err := cmd.printEvents(ctx, client, boardID); err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}
	b, err := client.Leaderboard(ctx, boardID)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	if err := cmd.printBoard(b); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (cmd *boardCmd) printBoard(b *aoc.Leaderboard) error {
	var nameMax int
	for _, m := range b.Members {
		nameMax = math.Max(nameMax, len(m.Name))
	}
	pad := strings.Repeat(" ", nameMax)
	fmt.Println(pad, "         1111111111222222")
	fmt.Println(pad, "1234567890123456789012345")
	slices.SortStableFunc(b.Members, func(i, j aoc.LeaderboardMember) bool {
		switch {
		case i.LocalScore > j.LocalScore:
			return true
		case i.LocalScore < j.LocalScore:
			return false
		case i.GlobalScore > j.GlobalScore:
			return true
		case i.GlobalScore < j.GlobalScore:
			return false
		default:
			return i.Name < j.Name
		}
	})
	for _, m := range b.Members {
		stars := make([]int, 25)
		slices.SortStableFunc(m.Days, func(i, j aoc.LeaderboardMemberDay) bool {
			return i.Day < j.Day
		})
		for _, d := range m.Days {
			if d.Day < 1 || d.Day > 25 {
				return fmt.Errorf("day %d out of range", d.Day)
			}
			if d.Part1 != nil {
				stars[d.Day-1]++
			}
			if d.Part2 != nil {
				stars[d.Day-1]++
			}
		}
		fmt.Printf("% *s ", nameMax, m.Name)
		for _, s := range stars {
			switch s {
			case 0:
				fmt.Print("\033[90m☆\033[0m")
			case 1:
				fmt.Print("\033[37m★\033[0m")
			case 2:
				fmt.Print("\033[33m★\033[0m")
			}
		}
		fmt.Printf(" %d\n", m.LocalScore)
	}
	return nil
}

func (cmd *boardCmd) printEvents(ctx context.Context, client *aoc.Client, id int) error {
	type event struct {
		time time.Time
		name string
		day  int
		star int
	}
	print := func(e event) {
		fmt.Printf("%s %s got day %d star %d\n", e.time.Format(time.DateTime), e.name, e.day, e.star)
	}
	var events []event
	for {
		b, err := client.Leaderboard(ctx, id)
		if err != nil {
			return err
		}
		var evs []event
		for _, m := range b.Members {
			for _, d := range m.Days {
				if d.Part1 != nil {
					evs = append(evs, event{
						time: d.Part1.Got.Local(),
						name: m.Name,
						day:  d.Day,
						star: 1,
					})
				}
				if d.Part2 != nil {
					evs = append(evs, event{
						time: d.Part2.Got.Local(),
						name: m.Name,
						day:  d.Day,
						star: 2,
					})
				}
			}
		}
		slices.SortStableFunc(evs, func(i, j event) bool {
			return i.time.Before(j.time)
		})
		i, j := 0, 0
		for i < len(events) && j < len(evs) {
			if events[i] == events[j] {
				i++
				j++
				continue
			}
			if events[i].time.Before(evs[j].time) {
				i++
				continue
			}
			print(evs[j])
			j++
		}
		for _, e := range evs[j:] {
			print(e)
		}
		events = evs
		time.Sleep(15 * time.Minute)
	}
}
