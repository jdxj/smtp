package module

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"smtp/util"
)

// Command 用于描述 SMTP 中的命令.
type Command struct {
	Cmd   string
	Param string
}

func (com *Command) String() string {
	return fmt.Sprintf("Cmd: %s, Param: %s\n", com.Cmd, com.Param)
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

// todo: json格式
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
func (m *MailMsg) ParseMail() error {
	// todo: 如果不是 multipart, 需要重新构造 part.
	boundary := m.ExtractBoundary()
	if boundary == "" {
		return fmt.Errorf("%s\n", "Cann't find boundary!")
	}

	pReader := multipart.NewReader(m.msg.Body, boundary)
	var res []*multipart.Part
	var contents []string
	var part *multipart.Part
	var err error
	for part, err = pReader.NextPart(); err == nil; part, err = pReader.NextPart() {
		str := Decode(part)
		contents = append(contents, str)
		res = append(res, part)
	}

	m.parts = res
	m.contents = contents
	return err
}

func (m *MailMsg) ExtractBoundary() string {
	if _, ok := m.msg.Header["Content-Type"]; !ok {
		// 可能不是 multipart
		return ""
	}

	ct := m.msg.Header["Content-Type"][0]
	media, param, err := mime.ParseMediaType(ct)
	if err != nil || media != "multipart/alternative" {
		return ""
	}

	return param["boundary"]
}

func (m *MailMsg) FromAddr() string {
	preAddr := m.msg.Header["From"][0]
	addr, err := mail.ParseAddress(preAddr)
	if err != nil {
		return ""
	}
	return addr.Address
}

func (m *MailMsg) ToAddr() string {
	preAddr := m.msg.Header["To"][0]
	addr, err := mail.ParseAddress(preAddr)
	if err != nil {
		return ""
	}
	return addr.Address
}

func (m *MailMsg) Json() ([]byte, error) {
	mj := &MailJson{
		Header: m.msg.Header,
		Addr: m.ToAddr(),
		Desc: "Receive a mail!",
	}

	for i, v := range m.parts {
		pj := &PartJson{
			Header:  v.Header,
			Content: m.contents[i],
		}
		mj.Contents = append(mj.Contents, pj)
	}

	return json.Marshal(mj)
}

// todo: 使用 part.Header 指定编码解码.
func Decode(part *multipart.Part) string {
	if part == nil {
		return ""
	}

	data, err := ioutil.ReadAll(part)
	if err != nil {
		util.SMTPLog.Println(err)
		return ""
	}

	// 解析传输编码
	if cte, ok := part.Header["Content-Transfer-Encoding"]; ok {
		media, _, err := mime.ParseMediaType(cte[0])
		if err != nil {
			util.SMTPLog.Println(err)
			return ""
		}

		switch media {
		case "base64":
			data, err = base64.StdEncoding.DecodeString(string(data))
			if err != nil {
				util.SMTPLog.Println("decode err: ", err)
				return ""
			}
		default:
			util.SMTPLog.Println("Not resolved Content-Transfer-Encoding: ", media)
		}
	}

	// 解析字符集编码
	if ct, ok := part.Header["Content-Type"]; ok {
		_, param, err := mime.ParseMediaType(ct[0])
		if err != nil {
			util.SMTPLog.Println(err)
			return ""
		}

		charset := param["charset"]
		switch charset {
		case "GBK", "gb18030", "gb2310":
			sRd := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
			data, err = ioutil.ReadAll(sRd)
			if err != nil {
				util.SMTPLog.Println("err: ", err)
				return ""
			}
		default:
			util.SMTPLog.Println("Not resolved charset: ", charset)
		}
	}
	return string(data)
}

// todo
func ParseTransferEncoding(data []byte) []byte {
	return nil
}
