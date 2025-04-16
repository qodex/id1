package main

type CmdOp int

const (
	CmdSet CmdOp = iota
	CmdGet
	CmdDel
	CmdMov
	CmdList
)

var opName = map[CmdOp]string{
	CmdSet:  "set",
	CmdGet:  "get",
	CmdDel:  "del",
	CmdMov:  "mov",
	CmdList: "list",
}

var nameOp = map[string]CmdOp{
	"set":  CmdSet,
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
