package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (t *Command) get() ([]byte, error) {
	filePath := filepath.Join(dbpath, t.Key)

	if info, err := os.Stat(filePath); os.IsNotExist(err) {
		return []byte{}, fmt.Errorf("not found")
	} else if info.IsDir() {
		return []byte{}, fmt.Errorf("not found")
	}

	if data, err := os.ReadFile(filePath); err != nil {
		return []byte{}, fmt.Errorf("not found")
	} else {
		return data, nil
	}
}
