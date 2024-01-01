package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/Merovius/AdventOfCode/internal/aoc"
	"github.com/google/subcommands"
)

func main() {
	log.SetFlags(0)

	var (
		cfgFile string
		event   string
	)
	if ucd, err := os.UserConfigDir(); err == nil {
		cfgFile = filepath.Join(ucd, "aoc", "config.json")
	}
	flag.StringVar(&cfgFile, "cfg", cfgFile, "config file to load")
	flag.StringVar(&event, "event", "", "year to use (defaults to current year)")
	flag.Parse()
	if cfgFile == "" {
		log.Fatal("-cfg is required")
	}
	cfg, err := loadconfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	client := &aoc.Client{
		SessionCookie: cfg.SessionCookie,
		Cache:         aoc.DirCache(""),
		Event:         event,
	}
	subcommands.Register(new(inputCmd), "")
	subcommands.Register(new(boardCmd), "")
	os.Exit(int(subcommands.Execute(context.Background(), cfg, client)))
}

type config struct {
	SessionCookie string         `json:"session_cookie"`
	Boards        map[string]int `json:"boards"`
	DefaultBoard  string         `json:"default_board"`
}

type optString struct {
	Valid bool
	Value string
}

func (s *optString) String() string {
	return s.Value
}

func (s *optString) Set(v string) error {
	s.Valid, s.Value = true, v
	return nil
}

func loadconfig(s string) (*config, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := new(config)
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	dec.UseNumber()
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, errors.New("trailing data after JSON config")
	}
	return cfg, nil
}
