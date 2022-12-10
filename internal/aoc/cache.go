package aoc

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Merovius/AdventOfCode/internal/sync"
	"github.com/google/renameio/v2"
)

// Cache stores temporary data.
type Cache interface {
	// Get the value associated with key.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set the value associated with key to val.
	Set(ctx context.Context, key string, val []byte) error
}

// DirCache stores data in a directory. If empty, defaults to
// filepath.Join(os.UserCacheDir(), "aoc", "cache")
type DirCache string

var userCacheDir sync.OnceValues[string, error]

func (c DirCache) dir() (string, error) {
	if c != "" {
		return string(c), nil
	}
	ucd, err := userCacheDir.Do(os.UserCacheDir)
	if err != nil || ucd == "" {
		return "", fmt.Errorf("could not determine default cache dir: %w", err)
	}
	d := filepath.Join(ucd, "aoc", "cache")
	if err := os.MkdirAll(d, 0700); err != nil {
		return "", err
	}
	return d, nil
}

// Get implements Cache.
func (c DirCache) Get(ctx context.Context, key string) ([]byte, error) {
	d, err := c.dir()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(d, filepath.FromSlash(key)))
}

// Set implements Cache.
func (c DirCache) Set(ctx context.Context, key string, val []byte) error {
	d, err := c.dir()
	if err != nil {
		return err
	}
	return renameio.WriteFile(filepath.Join(d, filepath.FromSlash(key)), val, 0600)
}

// Purge the cache.
func (c DirCache) Purge(ctx context.Context) error {
	d, err := c.dir()
	if err != nil {
		return err
	}
	de, err := os.ReadDir(d)
	if err != nil {
		return err
	}
	for _, e := range de {
		if err := os.RemoveAll(filepath.Join(d, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

type nopCache struct{}

func (nopCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, fs.ErrNotExist
}

func (nopCache) Set(ctx context.Context, key string, val []byte) error {
	return nil
}
