package module

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"
)

func NewReceiver(conn net.Conn) *Receiver {
	rer := &Receiver{
		conn: conn,
		bfr:  bufio.NewReader(conn),
		bfw:  bufio.NewWriter(conn),
	}
	log.Println("Create a Receiver!")
	return rer
}

type Receiver struct {
	conn net.Conn
	bfr  *bufio.Reader
	bfw  *bufio.Writer
}

func (rer *Receiver) Start() {
	defer rer.conn.Close()
	mail := &Mail{}

	// 问候
	rer.WriteReply(rer.ReplyGreetings())

	// 重复:
	//     1. 读命令
	//     2. 写回复
	// 收到 QUIT 就关闭连接.
	for {
		// todo: 邮件事务监控
		com := rer.ReadCommand()
		if com == nil {
			break
		}
		log.Println("cmd is: ", com.Cmd)

		switch com.Cmd {
		case "ehlo":
			rer.WriteReply(rer.ReplyEHLO())
		case "mail":
			rer.WriteReply(rer.ReplyMAIL())
		case "rcpt":
			rer.WriteReply(rer.ReplyRCPT())
		case "data":
			rer.WriteReply(rer.ReplyDATA())
			ok := rer.ReadData(mail)
			if ok {
				rer.WriteReply(rer.ReplyDataEnd())
			} else {
				rer.WriteReply(rer.ReplyDataFailure())
			}
		case "quit":
			rer.WriteReply(rer.ReplyQUIT())
			break
		default:
			log.Printf("Unresolved command: %s, data: %s", com.Cmd, com.String())
		}
	}
	fmt.Println(mail.ParseMail())
	log.Println("session is over!")
}

func (rer *Receiver) ReadCommand() *Command {
	for {
		line, err := rer.bfr.ReadString('\n')
		if err == io.EOF {
			log.Println("read eof, err:", err)
			return nil
		} else if err != nil {
			log.Println("err when read Command: ", err)
			time.Sleep(time.Second)
			continue
		}
		// todo: 对命令以及其参数的更详细的解析
		params := strings.Split(line, " ")
		count := len(params)
		if count == 0 {
			log.Println("read a bare line")
			continue
		} else if count > 0 {
			cmd := params[0]
			cmd = strings.TrimSpace(cmd)
			cmd = strings.ToLower(cmd)

			com := &Command{
				Cmd: cmd,
			}
			return com
		}
	}
}

func (rer *Receiver) ReadData(mail *Mail) bool {
	tr := textproto.NewReader(rer.bfr)
	mime, err := tr.ReadMIMEHeader()
	if err != nil {
		log.Println("read header err: ", err)
		return false
	}
	mail.mime = mime

	data, err := tr.ReadDotLines()
	if err == io.EOF {
		log.Println("read data accepted!")
	} else if err != nil {
		log.Println("read header err: ", err)
		return false
	}
	mail.data = data
	return true
}

func (rer *Receiver) WriteReply(rep *Reply) {
	if rep == nil {
		log.Println("rep is nil")
		return
	}

	n, err := rer.bfw.WriteString(rep.String())
	if err != nil {
		log.Printf("write count: %d. err: %s\n", n, err)
		return
	}

	rer.bfw.Flush()
}

func (rer *Receiver) ReplyGreetings() *Reply {
	rep := &Reply{
		StateCode: 220,
		Text:      "mail.aaronkir.xyz",
	}
	return rep
}

func (rer *Receiver) ReplyEHLO() *Reply {
	rep := &Reply{
		StateCode: 250,
		Text:      "mail.aaronkir.xyz",
	}
	return rep
}

func (rer *Receiver) ReplyMAIL() *Reply {
	return rer.ReplyEHLO()
}

func (rer *Receiver) ReplyRCPT() *Reply {
	return rer.ReplyEHLO()
}

func (rer *Receiver) ReplyDATA() *Reply {
	rep := &Reply{
		StateCode: 354,
		Text:      "Enter mail, end with '.' on a line by itself.",
	}
	return rep
}

func (rer *Receiver) ReplyDataEnd() *Reply {
	rep := &Reply{
		StateCode: 250,
		Text:      "Mail accepted",
	}
	return rep
}

func (rer *Receiver) ReplyDataFailure() *Reply {
	rep := &Reply{
		StateCode: 554,
		Text:      "Transaction failed.",
	}
	return rep
}

func (rer *Receiver) ReplyQUIT() *Reply {
	rep := &Reply{
		StateCode: 221,
		Text:      "QUIT.",
	}
	return rep
}
