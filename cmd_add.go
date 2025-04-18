package main

import (
	"os"
	"path/filepath"
)

func (t *Command) add() error {
	filePath := filepath.Join(dbpath, t.Key.String())
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(dir, 0770); mkdirErr != nil {
			return mkdirErr
		}
	}

	if f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	} else if _, err := f.WriteString(string(t.Data)); err != nil {
		return err
	} else {
		f.Close()
		return nil
	}
}
