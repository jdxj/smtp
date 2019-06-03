package api

import (
	"net/mail"
	"net/textproto"
)

type MailJson struct {
	Header   mail.Header `json:"header"`
	Contents []*PartJson `json:"contents"`
}

type PartJson struct {
	Header  textproto.MIMEHeader `json:"header"`
	Content string               `json:"content"`
}
