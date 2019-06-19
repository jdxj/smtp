package module

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/mail"
	"smtp/util"
	"strings"
	"time"
)

func NewReceiver(conn net.Conn) *Receiver {
	util.SMTPLog.Println("Create a Receiver!")
	rer := &Receiver{
		conn: conn,
		bfr:  bufio.NewReader(conn),
		bfw:  bufio.NewWriter(conn),
	}
	return rer
}

type Receiver struct {
	conn net.Conn
	bfr  *bufio.Reader
	bfw  *bufio.Writer
}

func (rer *Receiver) Start() {
	defer rer.conn.Close()
	defer util.SMTPLog.Println("Session is over!")
	// 问候
	rer.WriteReply(rer.ReplyGreetings())

	// 重复:
	//     1. 读命令
	//     2. 写回复
	// 收到 QUIT 就关闭连接.
	for {
		// todo: 邮件事务监控
		com, err := rer.ReadCommand()
		if err == io.EOF {
			return
		} else if err != nil {
			util.SMTPLog.Println(err)
			continue
		}

		switch com.Cmd {
		case "ehlo":
			rer.WriteReply(rer.ReplyEHLO())
		case "mail":
			rer.WriteReply(rer.ReplyMAIL())
		case "rcpt":
			rer.WriteReply(rer.ReplyRCPT())
		case "data":
			rer.WriteReply(rer.ReplyDATA())
			mailMsg, err := rer.ReadMail()
			if err != nil {
				rer.WriteReply(rer.ReplyDataFailure())
				return
			}

			err = mailMsg.ParseMail()
			if err != io.EOF {
				rer.WriteReply(rer.ReplyDataFailure())
				return
			}
			rer.WriteReply(rer.ReplyDataEnd())

			// 存储邮件
			//Store.M.Store(mailMsg.ToAddr(), mailMsg)
			//Store.DelMail(util.Dur, mailMsg.ToAddr())
			data, ok := Store.M.Load(mailMsg.ToAddr())
			if !ok {
				// 没找到对应的 WebSocketConn
				break
			}
			userInfo, _ := data.(*UserInfo)
			userInfo.Mail = mailMsg
			userInfo.PushMail()
		case ".":
			rer.WriteReply(rer.ReplyDataEnd())
		case "quit":
			rer.WriteReply(rer.ReplyQUIT())
			return
		case "reset":
			rer.WriteReply(rer.ReplyRESET())
		default:
			rer.WriteReply(rer.ReplyWongCmd())
			util.SMTPLog.Printf("Unresolved command: %s, data: %s", com.Cmd, com.String())
		}
	}

}

const readDur = 5 * time.Second

func (rer *Receiver) ReadCommand() (*Command, error) {
	// 用于超时检测
	lineChan := make(chan string)
	var eofErr error

	go func() {
		line, err := rer.bfr.ReadString('\n')
		if err != nil {
			eofErr = err
			return
		}
		lineChan <- line
	}()

	var line string
	select {
	case line = <-lineChan:
	case <-time.After(readDur):
		util.SMTPLog.Println("Read command timeout!")
		return &Command{Cmd: "quit"}, eofErr
	}

	line = strings.TrimSuffix(line, "\r\n")
	if line == "" {
		return nil, fmt.Errorf("%s\n", "Read a blank line!")
	}

	line = strings.ToLower(line)
	util.SMTPLog.Println("Command line is: ", line)

	params := strings.Split(line, " ")
	com := &Command{
		Cmd:   params[0],
		Param: strings.Join(params[1:], " "),
	}
	return com, nil
}

func (rer *Receiver) ReadMail() (*MailMsg, error) {
	mailMsg := &MailMsg{}
	msg, err := mail.ReadMessage(rer.bfr)
	if err != nil {
		return nil, err
	}

	mailMsg.msg = msg
	return mailMsg, nil
}

func (rer *Receiver) WriteReply(rep *Reply) {
	if rep == nil {
		util.SMTPLog.Println("rep is nil")
		return
	}

	n, err := rer.bfw.WriteString(rep.String())
	if err != nil {
		util.SMTPLog.Printf("write count: %d. err: %s\n", n, err)
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

func (rer *Receiver) ReplyWongCmd() *Reply {
	rep := &Reply{
		StateCode: 550,
		Text:      "Syntax error, command unrecognized.",
	}
	return rep
}

func (rer *Receiver) ReplyRESET() *Reply {
	return rer.ReplyEHLO()
}
