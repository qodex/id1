package id1

import (
	"testing"
)

func TestPublicKeyEnc(t *testing.T) {
	secret := generateSecret("test1")
	if data, err := encrypt(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC3gMw5zAsvBuJ+swdVW8Ec9r1zu42Z+m7TsgoaV6yas58hxrPCeBUoNhFmz380yBpXjB7jwX1f5nGrZA9FWt2hmtJNLCvr6U1ZMZeERbPWjFIE02BWK0p+qZKByjpNv+LYMr8YM/JfYmqhhVbhqno15vVFyfNmaVIB6y1yJtn7xQIDAQAB
-----END PUBLIC KEY-----`, secret); err != nil || len(data) == 0 {
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
