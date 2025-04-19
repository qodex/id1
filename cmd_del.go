package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (t *Command) del() error {
	path := filepath.Join(dbpath, t.Key.String())
	if stat, err := os.Stat(path); err != nil {
		return ErrNotFound
	} else if stat.IsDir() {
		return os.RemoveAll(path)
	} else {
		dotTtlPath := filepath.Join(dbpath, t.Key.Parent, fmt.Sprintf(".ttl.%s", t.Key.Name))
		os.Remove(dotTtlPath)
		return os.Remove(path)
	}
}
