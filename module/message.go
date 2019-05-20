package module

import (
	"fmt"
	"net/textproto"
)

// Command 用于描述 SMTP 中的命令.
type Command struct {
	Cmd    string
	Params map[string][]string
}

// todo: 实现
func (com *Command) String() string {
	return fmt.Sprintln(com)
}

// Reply 用于描述 SMTP 中的回复.
type Reply struct {
	StateCode int
	Text      string
}

// String 将 Reply 变成 line.
func (rep *Reply) String() string {
	return fmt.Sprintf("%d %s\r\n", rep.StateCode, rep.Text)
}

type Mail struct {
	mime textproto.MIMEHeader
	data []string
}

func (m *Mail) String() string {
	str := ""
	for k, v := range m.mime {
		str += fmt.Sprintf("%s: %s\n", k, v)
	}
	for i, v := range m.data {
		str += fmt.Sprintf("%s", v)
		if i != len(m.data)-1 {
			str += "\n"
		}
	}
	return str
}
