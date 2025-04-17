package main

import (
	"os"
	"path/filepath"
)

func (t *Command) get() ([]byte, error) {
	filePath := filepath.Join(dbpath, t.Key)

	if info, err := os.Stat(filePath); os.IsNotExist(err) {
		return []byte{}, ErrNotFound
	} else if info.IsDir() {
		return []byte{}, ErrNotFound
	}

	if data, err := os.ReadFile(filePath); err != nil {
		return []byte{}, err
	} else {
		return data, nil
	}
}
