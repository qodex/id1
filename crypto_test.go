package id1

import (
	"testing"
)

func TestPublicKeyEnc(t *testing.T) {
	msg := ""
	for range 117 { // max msg size
		msg += "a"
	}
	if _, err := encryptWithPubKey(testPubKey1, msg); err != nil {
		t.Errorf("encrypt failed: %s", err)
	}
}

func TestGenerateChallenge(t *testing.T) {
	challenge1, err1 := generateChallenge("testid1", testPubKey1)
	challenge2, err2 := generateChallenge("testid2", testPubKey2)
	if err1 != nil || err2 != nil {
		t.Errorf("error generating challenge")
	}
	if challenge1 == challenge2 {
		t.Errorf("invalid challenge")
	}
}
