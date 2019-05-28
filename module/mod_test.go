package module

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
	"testing"
)

var data = `DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed; d=qq.com; s=s201512;
	t=1558324854; bh=0V8UZaxOGRf8p8BdTAENiHY3Gk9n78qA7jAwMIvjQIQ=;
	h=From:To:Subject:Mime-Version:Date:Message-ID;
	b=Zg6T7oJqFUFyJDmtOG1cQ8WXHBSs95auDd65nC+AjWbSL0lmslBigfGLm0Ra0KeN1
	 HXM+0FYRc+ZBc4Qlu2TLEAt3QlgiYTSN1xIXbZpo0B0LcFCzfiD5ATbMIEtepVeqsS
	 +BFfaLw8hZ7NtihKtiNyGYfm6u+Hcp+xvhTQrDx0=
X-QQ-FEAT: tKEuZPfUcq1fqgpBV/yslW7iEcNVLuXICwGz9zgIO17nKLhbYI1FzQmHY7et8
	veDLanY7TgUyXC2/zTeIcObUC1tiktpngNdO7gB7+K142C+TQF0+F+KPj0l8aj5mYZpmFOf
	VVjRVc36rt4iJ3aslZhjRmzF6YjTQtsHngQCRt73f2tJY5zaCz5wmeaySU3DljQup3pvkVR
	Kqvb5/BHQGYLiqKiLK96GuZnjFHlu1QdbLs61wtExqbV/Y+Bok4Lh4YFViG5s+apYitJ8S/
	ysG68YcdjnOJHo9S5fZiV+8F0=
X-QQ-SSF: 00000000000000F000000000000000D
X-HAS-ATTACH: no
X-QQ-BUSINESS-ORIGIN: 2
X-Originating-IP: 122.224.122.173
X-QQ-STYLE: 
X-QQ-mid: webmail416t1558324853t2679444
From: "=?gb18030?B?zfXT0ce/?=" <985759262@qq.com>
To: "=?gb18030?B?amR4ag==?=" <jdxj@mail.aaronkir.xyz>
Subject: fa
Mime-Version: 1.0
Content-Type: multipart/alternative;
	boundary="----=_NextPart_5CE22675_0B105A80_6C22D71A"
Content-Transfer-Encoding: 8Bit
Date: Mon, 20 May 2019 12:00:52 +0800
X-Priority: 3
Message-ID: <tencent_17880353F60DD8C8D4DDF75A57F873DF9A06@qq.com>
X-QQ-MIME: TCMime 1.0 by Tencent
X-Mailer: QQMail 2.x
X-QQ-Mailer: QQMail 2.x
X-QQ-SENDSIZE: 520
Received: from qq.com (unknown [127.0.0.1])
	by smtp.qq.com (ESMTP) with SMTP
	id ; Mon, 20 May 2019 12:00:53 +0800 (CST)
Feedback-ID: webmail:qq.com:bgforeign:bgforeign2
X-QQ-Bgrelay: 1

This is a multi-part message in MIME format.

------=_NextPart_5CE22675_0B105A80_6C22D71A
Content-Type: text/plain;
	charset="gb18030"
Content-Transfer-Encoding: base64

ZGYNCg0KDQotLS0tLS0tLS0tLS0tLS0tLS0NCrLiytQx

------=_NextPart_5CE22675_0B105A80_6C22D71A
Content-Type: text/html;
	charset="gb18030"
Content-Transfer-Encoding: base64

PGRpdj5kZjwvZGl2PjxkaXY+PGJyPjwvZGl2PjxkaXY+PGRpdiBzdHlsZT0iY29sb3I6Izkw
OTA5MDtmb250LWZhbWlseTpBcmlhbCBOYXJyb3c7Zm9udC1zaXplOjEycHgiPi0tLS0tLS0t
LS0tLS0tLS0tLTwvZGl2PjxkaXYgc3R5bGU9ImZvbnQtc2l6ZToxNHB4O2ZvbnQtZmFtaWx5
OlZlcmRhbmE7Y29sb3I6IzAwMDsiPrLiytQxPC9kaXY+PC9kaXY+PGRpdj4mbmJzcDs8L2Rp
dj4=

------=_NextPart_5CE22675_0B105A80_6C22D71A--


.

`

func TestReceiver_ReadData(t *testing.T) {
	tr := textproto.NewReader(bufio.NewReader(strings.NewReader(data)))
	mime, err := tr.ReadMIMEHeader()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("count: ", len(mime))
	for k, v := range mime {
		fmt.Println("k: ", k)
		fmt.Printf("v: %s\n", v)
	}

	fmt.Println("---------------abc-------------------------------")

	lines, err := tr.ReadDotLines()
	if err != nil {
		fmt.Println("err when read dot: ", err)
		return
	}

	fmt.Println("lines's count: ", len(lines))
	for _, v := range lines {
		fmt.Println(v)
	}
}

func TestNewReceiver_Mail(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(data))
	tr := textproto.NewReader(r)
	mime, err := tr.ReadMIMEHeader()
	if err != nil {
		fmt.Println(err)
		return
	}

	v, ok := mime["Content-Type"]
	if !ok {
		fmt.Println("have no Content-Tye")
		return
	}

	n := strings.Index(v[0], "boundary=")
	if n < 0 {
		fmt.Println("have no boundary=")
		return
	}

	nextPart := v[0][n+9:]
	nextPart = strings.ReplaceAll(nextPart, "\"", "")
	nextParts := strings.Split(nextPart, "=")
	if len(nextParts) < 2 {
		fmt.Println("have no enough len!")
		return
	}
	nextPart = nextParts[1]

	lines, err := tr.ReadDotLines()
	if err != nil {
		fmt.Println("err: ", err)
		return
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
			return
		}
		sRd := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
		data, err = ioutil.ReadAll(sRd)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}

		realPats = append(realPats, string(data))
	}

	for _, v := range realPats {
		fmt.Println(v)
	}

}

func TestGoMail(t *testing.T) {
	fmt.Println("Test Go Mail ...")
	r := strings.NewReader(data)
	mailMsg, err := mail.ReadMessage(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range mailMsg.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	for k, v := range mailMsg.Header {
		fmt.Println("k: ", k)
		for i, v2 := range v {
			fmt.Printf("%d, %s\n", i, v2)
		}
		fmt.Println("-----------------")
	}
}

func TestSplitLine(t *testing.T) {
	r := strings.NewReader(data)
	msg, err := mail.ReadMessage(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	res := make(map[string]string)
	// v 的长度可能大于1, 大部分为1.
	for _, v1 := range msg.Header["Content-Type"] {
		parts := strings.Split(v1, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)

			if idx := strings.Index(part, "="); idx < 0 {
				res[part] = ""
			} else {
				res[part[:idx]] = part[idx+1:]
			}
		}
	}

	for k, v := range res {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func TestMultipart(t *testing.T) {
	r := strings.NewReader(data)
	msg, err := mail.ReadMessage(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	body := msg.Body
	partReader := multipart.NewReader(body, "----=_NextPart_5CE22675_0B105A80_6C22D71A")
	part, err := partReader.NextPart()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(part.Header)
	fmt.Println(part.FileName())
	fmt.Println(part.FormName())
}

func TestNewReceiver_Mail2(t *testing.T) {
	msg, err := mail.ReadMessage(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	boundary := "----=_NextPart_5CE22675_0B105A80_6C22D71A"
	pReader := multipart.NewReader(msg.Body, boundary)

	i := 0
	for part, err := pReader.NextPart(); err == nil; part, err = pReader.NextPart() {
		fmt.Printf("%d: %v\n", i, part.Header)
		buf := make([]byte, 4096)
		n, err := part.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("err 1: ", err)
			return
		}
		fmt.Printf("\tlen: %d, content: %s", n, buf[:n])
		i++
	}
}

func TestGoParseMediaType(t *testing.T) {
	msg, err := mail.ReadMessage(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	media, param, err := mime.ParseMediaType(msg.Header["Content-Type"][0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("media: ", media)
	for k, v := range param {
		fmt.Println("k:", k)
		fmt.Println(v)
	}
}

func TestStore(t *testing.T) {
	s := &Store{}
	s.Add(2)
	fmt.Println("ok")
}
