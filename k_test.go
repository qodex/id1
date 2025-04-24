package id1

import (
	"testing"
)

func TestK(t *testing.T) {
	k := K("TestId/pub/key")
	if k.Id != "TestId" || !k.Pub || k.Name != "key" || len(k.Segments) != 3 {
		t.Fail()
	}

	k = K("/TestId/")
	if k.Id != "TestId" || len(k.Segments) != 1 {
		t.Fail()
	}
}
