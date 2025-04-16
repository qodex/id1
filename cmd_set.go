package main

import (
	"os"
	"path/filepath"
)

func (t *Command) set() error {
	keyPath := filepath.Join(dbpath, t.Key)
	dir := filepath.Dir(keyPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(dir, 0770); mkdirErr != nil {
			return mkdirErr
		}
	}

	if err := os.WriteFile(keyPath, t.Data, 0644); err != nil {
		return err
	}

	return nil
}
