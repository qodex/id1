package main

import (
	"slices"
	"testing"
)

func setup() {
	dbpath = "test"
	CmdSet(K("test0/pub/key"), []byte("...")).Exec()
	CmdSet(K("test0/pub/tags/Robot"), []byte("...")).Exec()
	CmdSet(K("test0/.get"), []byte("Reader")).Exec()
	CmdSet(K("test0/pub/tags/.set"), []byte("Tagger")).Exec()
	CmdSet(K(".roles/max"), []byte("Reader")).Exec()
	CmdSet(K("test0/.roles/max"), []byte("Admin")).Exec()
	CmdSet(K("test0/pub/tags/.roles/max"), []byte("Tagger")).Exec()
}

func TestAuthDotOp(t *testing.T) {
	setup()
	if !authDotOp("max", CmdSet(K("test0/pub/tags/Robot"), []byte{})) {
		t.Fail()
	}
	if authDotOp("max", CmdDel(K("test0/pub/tags/Robot"))) {
		t.Fail()
	}
	if !authDotOp("max", CmdGet(K("test0/token"))) {
		t.Fail()
	}
	if authDotOp("max", CmdDel(K("test0/pub/key"))) {
		t.Fail()
	}
}

func TestGetRoles(t *testing.T) {
	setup()
	roles := getRoles("max", K("test0/pub/tags/Robot"))
	if !slices.Contains(roles, "Reader") ||
		!slices.Contains(roles, "Admin") ||
		!slices.Contains(roles, "Tagger") ||
		!slices.Contains(roles, "max") ||
		!slices.Contains(roles, "*") {
		t.Fail()
	}
}
