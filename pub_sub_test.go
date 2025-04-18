package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	pubsub := NewPubSub()
	var mu sync.Mutex
	received := 0
	subCount := 3
	for sub := range subCount {
		subId := fmt.Sprintf("subscriber%d", sub)
		ch := pubsub.Subscribe(subId)
		go func() {
			for {
				cmd := <-ch
				if cmd.Op != Set {
					t.Fail()
				}
				mu.Lock()
				received++
				mu.Unlock()
			}
		}()
	}

	pubCount := 3
	for pub := range pubCount {
		go func() {
			for sub := range subCount {
				key := KK(fmt.Sprintf("subscriber%d", sub), pub)
				cmd := CmdSet(key, []byte{})
				pubsub.Publish(cmd)
				time.Sleep(time.Millisecond)
			}
		}()
	}
	time.Sleep(time.Second * 1)
	if received != pubCount*subCount {
		t.Fail()
	}
}
