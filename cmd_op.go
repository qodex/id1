package main

type CmdOp int

const (
	CmdSet CmdOp = iota
	CmdAdd
	CmdGet
	CmdDel
	CmdMov
	CmdList
)

var opName = map[CmdOp]string{
	CmdSet:  "set",
	CmdAdd:  "add",
	CmdGet:  "get",
	CmdDel:  "del",
	CmdMov:  "mov",
	CmdList: "list",
}

var nameOp = map[string]CmdOp{
	"set":  CmdSet,
	"add":  CmdAdd,
	"get":  CmdGet,
	"del":  CmdDel,
	"mov":  CmdMov,
	"list": CmdList,
}

func (t CmdOp) String() string {
	return opName[t]
}

func NewCmdOp(s string) CmdOp {
	return nameOp[s]
}
