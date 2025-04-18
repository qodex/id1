package main

import (
	"log"
	"os"
	"path/filepath"
)

func (t *Command) move() error {
	oldKey := t.Key.String()
	newKey := string(t.Data)

	oldPath := filepath.Join(dbpath, oldKey)
	newPath := filepath.Join(dbpath, newKey)
	newDir := filepath.Dir(newPath)

	if _, err := os.Stat(oldPath); err != nil {
		return ErrNotFound
	}

	if _, err := os.Stat(newPath); err == nil {
		return ErrExists
	}

	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(newDir, 0770); mkdirErr != nil {
			return mkdirErr
		}
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		log.Printf("cmd: add error, %s", err)
		return err
	}

	return nil
}
