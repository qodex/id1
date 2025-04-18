package main

type Op int

const (
	Set Op = iota
	Add
	Get
	Del
	Mov
	List
)

var opName = map[Op]string{
	Set:  "set",
	Add:  "add",
	Get:  "get",
	Del:  "del",
	Mov:  "mov",
	List: "list",
}

var nameOp = map[string]Op{
	"set":  Set,
	"add":  Add,
	"get":  Get,
	"del":  Del,
	"mov":  Mov,
	"list": List,
}

func (t Op) String() string {
	return opName[t]
}

func op(s string) Op {
	return nameOp[s]
}
