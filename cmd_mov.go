package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (t *Command) move() error {
	oldKey := t.Key
	oldPath := filepath.Join(dbpath, oldKey)

	newKey := string(t.Data)
	newPath := filepath.Join(dbpath, newKey)
	newDir := filepath.Dir(newPath)

	if _, err := os.Stat(oldPath); err != nil {
		return fmt.Errorf("not found")
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("exists")
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
