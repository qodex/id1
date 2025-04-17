package main

import "strings"

type Id1Key struct {
	Id        string
	Pub       bool
	Last      string
	Timestamp int64
	Segments  []string
}

func (t Id1Key) String() string {
	return strings.Join(t.Segments, "/")
}

func K(s string) Id1Key {
	k := Id1Key{}

	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.Trim(s, "/")
	s = strings.ToLower(s)

	split := strings.Split(s, "/")
	k.Segments = split

	if len(split) > 0 {
		k.Id = split[0]
		k.Last = split[len(split)-1]
	}

	if len(split) > 1 {
		k.Pub = split[1] == "pub"
	}

	return k
}
