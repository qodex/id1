package main

import (
	"testing"
)

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

func TestAuth(t *testing.T) {
	dbpath = "test"
	if _, err := NewCommand(CmdSet, "testid1/pub/key", map[string]string{}, []byte("..........")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}

	if !auth("testid1", NewCommand(CmdSet, "testid1/pub/key", map[string]string{}, []byte{})) {
		t.Errorf("owner should be able to set own pub key")
	}

	if auth("testid2", NewCommand(CmdSet, "testid1/pub/key", map[string]string{}, []byte{})) {
		t.Errorf("non-owner shouldn't be able to set other pub keys")
	}

	if !auth("", NewCommand(CmdGet, "testid1/pub/key", map[string]string{}, []byte{})) {
		t.Errorf("anyone should be able to get pub keys")
	}
}

func TestParseClaims(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0aWQiLCJpYXQiOjE1MTYyMzkwMjJ9.m7GbsjZeOBZhdFfaU1_ulqeaogLi5gduLXqfLhyxH5w"
	if claims, err := parseToken(token, "test"); err != nil {
		t.Errorf("%s", err)
	} else if claims.Subject != "testid" {
		t.Errorf("expected 'testid' got %s", claims.Subject)
	}
}

func TestIdExists(t *testing.T) {
	dbpath = "test"
	if _, err := NewCommand(CmdSet, "testid/pub/key", map[string]string{}, []byte("..........")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}

	if !idExists("testid") {
		t.Errorf("expected exists")
	}

	if idExists("testid123") {
		t.Errorf("expected not exists")
	}
}

var testPubKey1 = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAM1W6T4DqjOTFHvsAUTebVF+NofSA3qJW7SF7gJTPh3IE0W6hkT0XSMP
Ue6eyS+2vITfmX5gShkm7z/HHpUS2Kho+Rj8HjRu0Ng68qbdpCkYcgkrrEJneX7U
WqmD6zw8RKkLA4Rsfu+wrTjf0ijxpS2vS0fzghyB9TcbsFzCo573AgMBAAE=
-----END RSA PUBLIC KEY-----`

var testPubKey2 = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBANDSsusgXGowG2Dsm2pCWyGbIEEGwsRgoKbUPx2JuVI0NWEvTrEmPfqa
H23ACLwetp4XMgZEYLmuS3PkA/HuQiUkYPElKEmfuO2jQ6F4/mHy6UkOsP9PMXwl
ff02vCJ43hBFIJdgchDSywHIb4F1hv6ap6PlrYMGwvIJ6gln9GIdAgMBAAE=
-----END RSA PUBLIC KEY-----`
