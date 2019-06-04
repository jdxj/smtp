package module

import (
	"net/mail"
	"net/textproto"
)

type MailJson struct {
	Header   mail.Header `json:"header,omitempty"`
	Contents []*PartJson `json:"contents,omitempty"`
	Addr string `json:"addr,omitempty"`
	// 可以添加一些描述
	Desc string `json:"desc,omitempty"`
}

type PartJson struct {
	Header  textproto.MIMEHeader `json:"header"`
	Content string               `json:"content"`
}
