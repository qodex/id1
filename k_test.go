package main

import (
	"testing"
)

func TestK(t *testing.T) {
	k := K("TestId/pub/key")
	if k.Id != "testid" || !k.Pub || k.Last != "key" || len(k.Segments) != 3 {
		t.Fail()
	}

	k = K("/TestId/")
	if k.Id != "testid" || len(k.Segments) != 1 {
		t.Fail()
	}
}
