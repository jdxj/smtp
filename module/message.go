package module

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
)

// Command 用于描述 SMTP 中的命令.
type Command struct {
	Cmd    string
	Params map[string][]string
}

// todo: 实现
func (com *Command) String() string {
	return fmt.Sprintln("test string")
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

// todo: 更好的格式化输出

type MailMsg struct {
	msg   *mail.Message
	parts []*multipart.Part
	// contents 是与 parts 对应的数据内容.
	contents []string
}

func (m *MailMsg) String() string {
	str := ""
	for k, v := range m.msg.Header {
		str += fmt.Sprintf("%s: %s\n", k, v)
	}

	for i, p := range m.parts {
		str += fmt.Sprintf("Part %d:\n", i)
		str += fmt.Sprintf("\tHeader: %s\n", p.Header)
		str += fmt.Sprintf("\tContent: %s\n", m.contents[i])
	}
	return str
}

// ParseMail 用于解析 multipart/alternative 邮件部分
func (m *MailMsg) ParseMail() {
	boundary := m.ExtractBoundary()
	if boundary == "" {
		return
	}

	pReader := multipart.NewReader(m.msg.Body, boundary)
	var res []*multipart.Part
	var contents []string
	for part, err := pReader.NextPart(); err == nil; part, err = pReader.NextPart() {
		str := Decode(part)
		contents = append(contents, str)
		res = append(res, part)
	}

	m.parts = res
	m.contents = contents
}

func (m *MailMsg) ExtractBoundary() string {
	ct := m.msg.Header["Content-Type"][0]
	media, param, err := mime.ParseMediaType(ct)
	if err != nil || media != "multipart/alternative" {
		return ""
	}

	return param["boundary"]
}

// todo: 使用 part.Header 指定编码解码.
func Decode(part *multipart.Part) string {
	if part == nil {
		return ""
	}

	data, err := ioutil.ReadAll(part)
	if err != nil {
		log.Println(err)
		return ""
	}

	data, err = base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Println("decode err: ", err)
		return ""
	}

	sRd := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	data, err = ioutil.ReadAll(sRd)
	if err != nil {
		log.Println("err: ", err)
		return ""
	}
	return string(data)
}
