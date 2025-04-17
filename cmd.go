package main

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type Command struct {
	Op   CmdOp
	Key  string
	Args map[string]string
	Data []byte
}

func NewCommand(op CmdOp, key string, args map[string]string, data []byte) Command {
	cmd := Command{
		Op:   op,
		Key:  key,
		Args: args,
		Data: data,
	}
	return cmd
}

func (t Command) Bytes() []byte {
	args := url.Values{}
	for arg := range t.Args {
		args.Set(arg, t.Args[arg])
	}
	url := url.URL{
		Scheme:   t.Op.String(),
		Path:     t.Key,
		RawQuery: args.Encode(),
	}
	command := strings.ReplaceAll(url.String(), "//", "/")
	bytes := slices.Concat([]byte(command), []byte("\n"), t.Data)
	return bytes
}

func (t Command) String() string {
	return string(t.Bytes())
}

/*
First line is command args, the rest is data
<op>:/<key>?<args... k1=v1&k2=k2..>
[data]

Examples:
set:/max/msg/1744641370068551?ttl=60
Hi...

get:/max/msg/1744641370068551
*/
func ParseCommand(data []byte) (Command, error) {
	command := Command{}
	firstLineEnd := slices.Index(data, byte('\n'))
	if firstLineEnd < 0 {
		firstLineEnd = len(data)
		data = append(data, byte('\n'))
	}
	firstLine := string(data[0:firstLineEnd])
	command.Data = data[firstLineEnd+1:]

	if strings.HasPrefix(firstLine, "#") {
		return command, fmt.Errorf("not a command")
	}

	url, err := url.Parse(firstLine)
	if err != nil {
		return command, err
	}
	command.Op = NewCmdOp(url.Scheme)
	command.Key = strings.TrimPrefix(url.Path, "/")
	command.Args = map[string]string{}
	for k := range url.Query() {
		command.Args[k] = url.Query().Get(k)
	}
	return command, nil
}

func (t Command) Exec() ([]byte, error) {
	switch t.Op {
	case CmdSet:
		return []byte{}, t.set()
	case CmdAdd:
		return []byte{}, t.add()
	case CmdGet:
		return t.get()
	case CmdDel:
		return []byte{}, t.del()
	case CmdMov:
		return []byte{}, t.move()
	case CmdList:
		return t.list()
	default:
		return []byte{}, fmt.Errorf("not supported")
	}
}
