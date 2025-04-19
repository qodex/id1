package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDotAfter(t *testing.T) {
	dbpath = "test"
	ttlKey := K("testda/1sec/qqqqbbbb")
	NewCommand(Set, ttlKey, map[string]string{"ttl": "1"}, []byte("...")).Exec()
	time.Sleep(time.Second)
	dotAfter(dbpath)
	if _, err := CmdGet(ttlKey).Exec(); err == nil {
		t.Fail()
	}
	if _, err := CmdGet(K(fmt.Sprintf("%s/.ttl.%s", ttlKey.Parent, ttlKey.Name))).Exec(); err == nil {
		t.Fail()
	}
}
