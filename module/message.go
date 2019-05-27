package module

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/mail"
	"net/textproto"
	"strings"
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

type MailMsg struct {
	msg *mail.Message
}

func (m *MailMsg) String() string {
	str := ""
	for k, v := range m.msg.Header {
		str += fmt.Sprintf("%s: %s\n", k, v)
	}
	body, _ := ioutil.ReadAll(m.msg.Body)
	return str + string(body)
}

// ParseMail 用于解析 multipart/alternative 邮件部分
func (m *MailMsg) ParseMail() string {
	v, ok := m.msg.Header["Content-Type"]
	if !ok {
		return ""
	}

	n := strings.Index(v[0], "boundary=")
	if n < 0 {
		fmt.Println("have no boundary=")
		return ""
	}

	nextPart := v[0][n+9:]
	nextPart = strings.ReplaceAll(nextPart, "\"", "")
	nextParts := strings.Split(nextPart, "=")
	if len(nextParts) < 2 {
		fmt.Println("have no enough len!")
		return ""
	}
	nextPart = nextParts[1]

	lines, err := textproto.NewReader(bufio.NewReader(m.msg.Body)).ReadDotLines()
	if err != nil {
		log.Println("err: ", err)
		return ""
	}

	var parts []string
	part := ""
	for _, v := range lines {
		if strings.Index(v, nextPart) >= 0 {
			parts = append(parts, part)
			part = ""
			continue
		}

		part += v
	}

	var realPats []string
	for _, v := range parts {
		i := strings.Index(v, "base64")
		if i < 0 {
			continue
		}
		data, err := base64.StdEncoding.DecodeString(string(v[i+6:]))
		if err != nil {
			fmt.Println("decode err: ", err)
			return ""
		}
		sRd := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
		data, err = ioutil.ReadAll(sRd)
		if err != nil {
			fmt.Println("err: ", err)
			return ""
		}

		realPats = append(realPats, string(data))
	}

	res := ""
	for i, v := range realPats {
		res += fmt.Sprintf("第 %d 部分\n", i)
		res += v
		res += "\n"
	}
	return res
}
