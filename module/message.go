package module

import "fmt"

// Command 用于描述 SMTP 中的命令.
type Command struct {
	Cmd string
	Params map[string][]string
}

// todo: 实现
func (com *Command) String() string {
	return ""
}

// Reply 用于描述 SMTP 中的回复.
type Reply struct {
	StateCode int
	Text string
}

// String 将 Reply 变成 line.
func (rep *Reply) String() string {
	return fmt.Sprintf("%d %s\r\n", rep.StateCode, rep.text)
}