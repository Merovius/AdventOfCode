// Package aoc implements a simple client for the AoC website.
package aoc

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"github.com/Merovius/AdventOfCode/internal/sync"
	"golang.org/x/exp/slices"
	"golang.org/x/net/publicsuffix"
)

// AoCWebsite is the default website of Advent of Code.
const AoCWebsite = "https://adventofcode.com/"

// Client to the AoC Website. The zero value is valid and uses sensible
// defaults.
type Client struct {
	// SessionCookie to use. Must be set to access user data.
	SessionCookie string

	// BaseURL of the AoC website. Defaults to AoCWebsite.
	BaseURL string

	// Event to use, where needed. Defaults to current year.
	Event string

	// Cache for temporary data. If nil, no cache is used.
	Cache Cache

	Jar http.CookieJar

	initOnce sync.OnceValue[error]
	client   *http.Client
	base     *url.URL
	baseID   string
	cache    Cache
}

func (c *Client) init() error {
	return c.initOnce.Do(func() (err error) {
		if c.BaseURL != "" {
			c.base, err = url.Parse(c.BaseURL)
		} else {
			c.base, err = url.Parse(AoCWebsite)
		}
		if err != nil {
			return err
		}
		h := sha1.New()
		io.WriteString(h, c.BaseURL)
		c.baseID = base64.RawURLEncoding.EncodeToString(h.Sum(nil))

		jar := c.Jar
		if jar == nil {
			// TODO: this should probably be loaded from a cookies.txt file
			jar, err = cookiejar.New(&cookiejar.Options{
				PublicSuffixList: publicsuffix.List,
			})
			if err != nil {
				return err
			}
		}
		if c.SessionCookie != "" {
			jar.SetCookies(c.base, []*http.Cookie{
				{Name: "session", Value: c.SessionCookie},
			})
		}
		if c.Event == "" {
			c.Event = strconv.Itoa(time.Now().UTC().Year())
		}
		c.client = &http.Client{
			Jar: jar,
		}
		c.cache = c.Cache
		if c.cache == nil {
			c.cache = nopCache{}
		}
		return nil
	})
}

// Leaderboard fetches the private leaderboard with the given id.
func (c *Client) Leaderboard(ctx context.Context, id int) (*Leaderboard, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	if c.SessionCookie == "" {
		return nil, errors.New("c.SessionCookie must be set to fetch private leaderboards")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", c.base.JoinPath(c.Event, "leaderboard", "private", "view", strconv.Itoa(id)+".json").String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read at most 1MB
	buf, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	b := new(Leaderboard)
	if err := json.Unmarshal(buf, b); err != nil {
		return nil, err
	}
	return b, nil
}

// Input fetches the input for the given day and writes it to w. If day is 0,
// it defaults to the most recent day.
func (c *Client) Input(ctx context.Context, day int, w io.Writer) error {
	if err := c.init(); err != nil {
		return err
	}
	if c.SessionCookie == "" {
		return errors.New("c.SessionCookie must be set to fetch puzzle inputs")
	}
	if day == 0 {
		if now := time.Now().UTC(); now.Hour() < 5 {
			day = now.Day() - 1
		} else {
			day = now.Day()
		}
		if day > 25 {
			day = 25
		}
	}
	req, err := http.NewRequestWithContext(ctx, "GET", c.base.JoinPath(c.Event, "day", strconv.Itoa(day), "input").String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return resp.Body.Close()
}

// Leaderboard is the data for a private leaderboard.
type Leaderboard struct {
	Owner   int
	Event   string
	Members []LeaderboardMember
}

// LeaderboardMember is the data for the member of a private leaderboard.
type LeaderboardMember struct {
	ID          int
	Name        string
	Stars       int
	LastStar    time.Time
	LocalScore  int
	GlobalScore int
	Days        []LeaderboardMemberDay
}

// LeaderboardMemberDay represents the stars a given LeaderboardMember got on a
// given day.
type LeaderboardMemberDay struct {
	Day   int
	Part1 *LeaderboardStar
	Part2 *LeaderboardStar
}

// LeaderboardStar represents a star gotten by a LeaderboardMember on a
// LeaderboardMemberDay.
type LeaderboardStar struct {
	Index int
	Got   time.Time
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Leaderboard) UnmarshalJSON(buf []byte) error {
	var err error
	type completion struct {
		StarIndex int   `json:"star_index"`
		GetStarTS int64 `json:"get_start_ts"`
	}
	type member struct {
		Stars              int                              `json:"stars"`
		LastStarTS         int64                            `json:"last_star_ts"`
		CompletionDayLevel map[string]map[string]completion `json:"completion_day_level"`
		LocalScore         int                              `json:"local_score"`
		Name               string                           `json:"name"`
		ID                 int                              `json:"id"`
		GlobalScore        int                              `json:"global_score"`
	}
	type leaderboard struct {
		Owner   int               `json:"owner_id"`
		Event   string            `json:"event"`
		Members map[string]member `json:"members"`
	}
	var board leaderboard
	if err := json.Unmarshal(buf, &board); err != nil {
		return err
	}
	b.Owner = board.Owner
	b.Event = board.Event
	for _, m := range board.Members {
		var days []LeaderboardMemberDay
		for d, c := range m.CompletionDayLevel {
			var day LeaderboardMemberDay
			day.Day, err = strconv.Atoi(d)
			if err != nil {
				return fmt.Errorf("invalid completion_day_level key %q", d)
			}
			if star, ok := c["1"]; ok {
				day.Part1 = &LeaderboardStar{
					Index: star.StarIndex,
					Got:   time.Unix(star.GetStarTS, 0).Local(),
				}
			}
			if star, ok := c["2"]; ok {
				day.Part2 = &LeaderboardStar{
					Index: star.StarIndex,
					Got:   time.Unix(star.GetStarTS, 0).Local(),
				}
			}
			days = append(days, day)
		}
		slices.SortFunc(days, func(a, b LeaderboardMemberDay) bool {
			return a.Day < b.Day
		})
		b.Members = append(b.Members, LeaderboardMember{
			ID:          m.ID,
			Name:        m.Name,
			Stars:       m.Stars,
			LastStar:    time.Unix(m.LastStarTS, 0).Local(),
			LocalScore:  m.LocalScore,
			GlobalScore: m.GlobalScore,
			Days:        days,
		})
	}
	slices.SortFunc(b.Members, func(a, b LeaderboardMember) bool {
		return a.ID < b.ID
	})
	return nil
}
