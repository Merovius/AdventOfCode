package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/aoc"
	"github.com/google/renameio/v2"
	"github.com/google/subcommands"
)

type inputCmd struct {
	day     int
	outName string
}

func (*inputCmd) Name() string {
	return "input"
}

func (*inputCmd) Synopsis() string {
	return "Fetch the puzzle input"
}

func (*inputCmd) Usage() string {
	return "Fetch the puzzle input for a given day. Defaults to the most recent day."
}

func (cmd *inputCmd) SetFlags(fs *flag.FlagSet) {
	fs.IntVar(&cmd.day, "day", 0, "day to fetch")
	fs.StringVar(&cmd.outName, "out", "-", "file to write the output to")
}

func (cmd *inputCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	log.SetFlags(0)
	if f.NArg() > 0 {
		log.Println("usage: input [-day <day>] [-out <outname>]")
		return subcommands.ExitUsageError
	}
	var w io.WriteCloser
	if cmd.outName == "-" {
		w = os.Stdout
	} else {
		pf, err := renameio.NewPendingFile(cmd.outName)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
		defer pf.Cleanup()
		w = (*pendingFile)(pf)
	}
	_, client := args[0].(*config), args[1].(*aoc.Client)
	if err := client.Input(ctx, cmd.day, w); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	if err := w.Close(); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

type pendingFile renameio.PendingFile

func (f *pendingFile) Write(p []byte) (int, error) {
	return (*renameio.PendingFile)(f).Write(p)
}

func (f *pendingFile) Close() error {
	return (*renameio.PendingFile)(f).CloseAtomicallyReplace()
}
