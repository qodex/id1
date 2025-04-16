package main

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ListOptions struct {
	Limit     int
	SizeLimit int
	Keys      bool
	Recursive bool
	Children  bool
}

func (t *ListOptions) Parse(args map[string]string) {
	if i, err := strconv.ParseInt(args["limit"], 10, 64); err == nil {
		t.Limit = int(i)
	} else {
		t.Limit = 999999
	}
	if i, err := strconv.ParseInt(args["size-limit"], 10, 64); err == nil {
		t.SizeLimit = int(i)
	} else {
		t.SizeLimit = 999999999999
	}
	t.Keys = args["keys"] == "true"

	t.Recursive = true
	if args["recursive"] == "false" {
		t.Recursive = false
	}

	t.Children = args["children"] == "true"
	if t.Children {
		t.Recursive = false
	}
}

func (t *Command) list() ([]byte, error) {
	opt := ListOptions{}
	opt.Parse(t.Args)
	key := strings.TrimSuffix(t.Key, "*")
	dirPath := filepath.Join(dbpath, key)

	if stat, err := os.Stat(dirPath); err != nil || !stat.IsDir() {
		return []byte{}, fmt.Errorf("not found")
	}

	var results map[string][]byte

	if opt.Recursive {
		results = walkDir(dirPath, opt)
	} else {
		results = listDir(dirPath, opt)
	}

	list := []string{}
	for k, v := range results {
		line := fmt.Sprintf("%s=%s", k, base64.StdEncoding.EncodeToString(v))
		list = append(list, line)

	}
	return []byte(strings.Join(list, "\n")), nil
}

func listDir(path string, opt ListOptions) map[string][]byte {
	results := map[string][]byte{}
	if entries, err := os.ReadDir(path); err != nil {
		log.Printf("error listing dir %s: %s", path, err)
	} else {
		for _, e := range entries {
			if len(results) >= opt.Limit {
				break
			}
			itemPath := filepath.Join(path, e.Name())
			if stat, err := os.Stat(itemPath); err != nil {
				log.Printf("cmd: list error, %s", err)
			} else if stat.Size() > int64(opt.SizeLimit) {
				continue
			}

			var key string
			if opt.Children {
				key = e.Name()
			} else {
				key = strings.TrimPrefix(itemPath, dbpath)
				key = strings.TrimPrefix(key, "/")
			}

			if opt.Keys {
				results[key] = []byte{}
				continue
			}

			if data, err := os.ReadFile(itemPath); err != nil {
				log.Printf("cmd: list error, %s", err)
			} else {
				results[key] = data
			}
		}
	}
	return results
}

func walkDir(path string, opt ListOptions) map[string][]byte {
	results := map[string][]byte{}
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if len(results) >= opt.Limit {
			return nil
		}

		if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		if stat, err := os.Stat(path); err != nil {
			log.Printf("cmd: list error, %s", err)
			return nil
		} else if stat.Size() > int64(opt.SizeLimit) {
			return nil
		}

		key := strings.TrimPrefix(path, dbpath)
		key = strings.TrimPrefix(key, "/")

		if opt.Keys {
			results[key] = []byte{}
		} else if data, err := os.ReadFile(path); err != nil {
			log.Printf("cmd: list error, %s", err)
		} else {
			results[key] = data
		}
		return nil
	})

	return results
}
