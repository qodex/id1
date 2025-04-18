package main

import (
	"os"
	"path/filepath"
)

func (t *Command) del() error {
	keyPath := filepath.Join(dbpath, t.Key.String())
	if stat, err := os.Stat(keyPath); err != nil {
		return ErrNotFound
	} else if stat.IsDir() {
		return os.RemoveAll(keyPath)
	} else {
		return os.Remove(keyPath)
	}
}
