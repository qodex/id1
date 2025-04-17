package main

import (
	"log"
	"strings"
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
	NewCommand(CmdDel, "test", map[string]string{}, []byte{}).Exec()

	if _, err := NewCommand(CmdSet, "test/one", map[string]string{}, []byte("1")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}
	if data, err := NewCommand(CmdGet, "test/one", map[string]string{}, []byte{}).Exec(); err != nil || string(data) != "1" {
		t.Errorf("get err: %s, data: '%s'", err, string(data))
	}

	if _, err := NewCommand(CmdSet, "test/one", map[string]string{}, []byte("11")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}
	if data, err := NewCommand(CmdGet, "test/one", map[string]string{}, []byte{}).Exec(); err != nil || string(data) != "11" {
		t.Errorf("get err: %s, data: '%s'", err, string(data))
	}

	if _, err := NewCommand(CmdDel, "test/one", map[string]string{}, []byte{}).Exec(); err != nil {
		t.Errorf("del err: %s", err)
	}
	if data, err := NewCommand(CmdGet, "test/one", map[string]string{}, []byte{}).Exec(); err == nil || string(data) == "1" {
		t.Errorf("del get err: %s, data: %s", err, string(data))
	}
}

func TestCommandMov(t *testing.T) {
	dbpath = "test"
	NewCommand(CmdDel, "test", map[string]string{}, []byte{}).Exec()

	if _, err := NewCommand(CmdSet, "test/one", map[string]string{}, []byte("1")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}
	if _, err := NewCommand(CmdMov, "test/one", map[string]string{}, []byte("test/two")).Exec(); err != nil {
		t.Errorf("mov err: %s", err)
	}
	if data, err := NewCommand(CmdGet, "test/one", map[string]string{}, []byte{}).Exec(); err == nil {
		t.Errorf("mov err: %s, data: %s", err, string(data))
	}
	if data, err := NewCommand(CmdGet, "test/two", map[string]string{}, []byte{}).Exec(); err != nil || string(data) != "1" {
		t.Errorf("mov err: %s, data: '%s'", err, string(data))
	}
}

func TestCommandList(t *testing.T) {
	dbpath = "test"
	NewCommand(CmdDel, "test", map[string]string{}, []byte{}).Exec()

	if _, err := NewCommand(CmdSet, "test/one", map[string]string{}, []byte("1")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}
	if _, err := NewCommand(CmdSet, "test/two", map[string]string{}, []byte("22")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}
	if _, err := NewCommand(CmdSet, "test/three", map[string]string{}, []byte("333")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}
	if _, err := NewCommand(CmdSet, "test/sub/one", map[string]string{}, []byte("sub/1")).Exec(); err != nil {
		t.Errorf("set err: %s", err)
	}

	if data, err := NewCommand(CmdList, "test", map[string]string{"recursive": "true"}, []byte{}).Exec(); err != nil || len(strings.Split(string(data), "\n")) != 4 {
		t.Errorf("err: %s, data: %s", err, string(data))
	}
	if data, err := NewCommand(CmdList, "test", map[string]string{}, []byte{}).Exec(); err != nil || len(strings.Split(string(data), "\n")) != 3 {
		t.Errorf("err: %s, data: %s", err, string(data))
	}
	if data, err := NewCommand(CmdList, "test", map[string]string{"limit": "2"}, []byte{}).Exec(); err != nil || len(strings.Split(string(data), "\n")) != 2 {
		t.Errorf("err: %s, data: %s", err, string(data))
	}
	if data, err := NewCommand(CmdList, "test", map[string]string{"keys": "true"}, []byte{}).Exec(); err != nil || strings.Contains(string(data), "22") {
		t.Errorf("err: %s, data: %s", err, string(data))
	}
	if data, err := NewCommand(CmdList, "test", map[string]string{"size-limit": "1"}, []byte{}).Exec(); err != nil || strings.Contains(string(data), "22") {
		t.Errorf("err: %s, data: %s", err, string(data))
	}
}
