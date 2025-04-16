package main

import (
	"log"
	"testing"
)

func TestParseCommand(t *testing.T) {
	str := "set:/max/msg/1731664334195180?ttl=600\ndata..."
	cmd, err := ParseCommand([]byte("set:/max/msg/1731664334195180?ttl=600\ndata..."))

	fail := err != nil ||
		cmd.Op != CmdSet ||
		cmd.Key != "max/msg/1731664334195180" ||
		cmd.Args["ttl"] != "600" ||
		string(cmd.Data) != "data..." ||
		cmd.String() != str

	if fail {
		log.Printf("err=%s, op=%s, key=%s, args: %s, data: %s", err, cmd.Op, cmd.Key, cmd.Args, cmd.Data)
		log.Printf("String: %s", cmd.String())
		t.Fail()
	}
}

func TestCommandCRUD(t *testing.T) {
	dbpath = "test"
	if _, err := NewCommand(CmdSet, "test/one", map[string]string{}, []byte("one")).Exec(); err != nil {
		t.Fail()
	}
	if data, err := NewCommand(CmdGet, "test/one", map[string]string{}, []byte{}).Exec(); err != nil || string(data) != "one" {
		t.Fail()
	}
	if _, err := NewCommand(CmdDel, "test/one", map[string]string{}, []byte{}).Exec(); err != nil {
		t.Fail()
	}
	if data, err := NewCommand(CmdGet, "test/one", map[string]string{}, []byte{}).Exec(); err == nil || string(data) == "one" {
		t.Fail()
	}
}
