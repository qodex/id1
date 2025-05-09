package id1

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type Command struct {
	Op   Op
	Key  Id1Key
	Args map[string]string
	Data []byte
}

func NewCommand(op Op, key Id1Key, args map[string]string, data []byte) Command {
	cmd := Command{
		Op:   op,
		Key:  key,
		Args: args,
		Data: data,
	}
	return cmd
}

func CmdSet(key Id1Key, args map[string]string, data []byte) Command {
	return NewCommand(Set, key, args, data)
}

func CmdGet(key Id1Key) Command {
	return NewCommand(Get, key, map[string]string{}, []byte{})
}

func CmdList(key Id1Key, opt ListOptions) Command {
	return NewCommand(List, key, opt.Map(), []byte{})
}

func CmdDel(key Id1Key) Command {
	return NewCommand(Del, key, map[string]string{}, []byte{})
}

func CmdMov(src Id1Key, tgt Id1Key) Command {
	return NewCommand(Mov, src, map[string]string{}, []byte(tgt.String()))
}

func (t Command) Bytes() []byte {
	bytes := fmt.Appendf(nil, "%s\n%s", t.String(), t.Data)
	return bytes
}

func (t Command) String() string {
	args := url.Values{}
	for arg := range t.Args {
		args.Set(arg, t.Args[arg])
	}
	url := url.URL{
		Scheme:   t.Op.String(),
		Path:     t.Key.String(),
		RawQuery: args.Encode(),
	}
	cmdStr := strings.ReplaceAll(url.String(), "//", "/")
	return cmdStr
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
	command.Op = op(url.Scheme)
	command.Key = K(url.Path)
	command.Args = map[string]string{}
	for k := range url.Query() {
		command.Args[k] = url.Query().Get(k)
	}
	return command, nil
}

func (t Command) Exec() ([]byte, error) {
	switch t.Op {
	case Set:
		return []byte{}, t.set()
	case Add:
		return []byte{}, t.add()
	case Get:
		return t.get()
	case Del:
		return []byte{}, t.del()
	case Mov:
		return []byte{}, t.move()
	case List:
		return t.list()
	default:
		return []byte{}, fmt.Errorf("not supported")
	}
}
