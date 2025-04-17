package main

import (
	"log"
	"os"
	"path/filepath"
)

func (t *Command) del() error {
	keyPath := filepath.Join(dbpath, t.Key)
	if stat, err := os.Stat(keyPath); err != nil {
		log.Printf("cmd: del error, %s", err)
		return ErrNotFound
	} else if stat.IsDir() {
		return os.RemoveAll(keyPath)
	} else {
		return os.Remove(keyPath)
	}
}
