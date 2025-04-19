package main

import (
	"testing"
)

func TestAuth(t *testing.T) {
	dbpath = "test"
	testid1PubKey := K("testid1/pub/key")
	if _, err := NewCommand(Set, K("testid1/pub/key"), map[string]string{}, []byte("..........")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}

	if !auth("testid1", NewCommand(Set, testid1PubKey, map[string]string{}, []byte{})) {
		t.Errorf("owner should be able to set own pub key")
	}

	if auth("testid2", NewCommand(Set, testid1PubKey, map[string]string{}, []byte{})) {
		t.Errorf("non-owner shouldn't be able to set other pub keys")
	}

	if !auth("", NewCommand(Get, testid1PubKey, map[string]string{}, []byte{})) {
		t.Errorf("anyone should be able to get pub keys")
	}
}

func TestParseClaims(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0aWQiLCJpYXQiOjE1MTYyMzkwMjJ9.m7GbsjZeOBZhdFfaU1_ulqeaogLi5gduLXqfLhyxH5w"

	if claims, err := validateToken(token, "test"); err != nil {
		t.Errorf("err: %s", err)
	} else if claims.Subject != "testid" {
		t.Errorf("expected 'testid' got %s", claims.Subject)
	}
}

func TestIdExists(t *testing.T) {
	dbpath = "test"
	testidPubKey := K("testid1/pub/key")
	if _, err := NewCommand(Set, testidPubKey, map[string]string{}, []byte("..........")).Exec(); err != nil {
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
