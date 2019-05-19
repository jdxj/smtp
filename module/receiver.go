package module

import (
	"bufio"
	"log"
	"net"
)

func NewReceiver(conn net.Conn) *Receiver {
	rer := &Receiver{
		conn: conn,
		bfr: bufio.NewReader(conn),
		bfw: bufio.NewWriter(conn),
	}
	log.Println("Create a Receiver!")
	return rer
}

type Receiver struct {
	conn net.Conn
	bfr *bufio.Reader
	bfw *bufio.Writer
}


func (rer *Receiver) Start() {
	defer rer.conn.Close()

	// 问候
	greetings := rer.ReplyGreetings()
	n, err := rer.bfw.WriteString(greetings)
	if err != nil {
		log.Fatalf("write count: %d. err: %s", n, err)
	}

	rer.bfw.Flush()

	// 重复:
	//     1. 读命令
	//     2. 写回复
	// 收到 QUIT 就关闭连接.
	for {
		com := rer.ReadCommand()

		rer.WriteReply(com)


	}
}

func (rer *Receiver) ReplyGreetings() string {
	rep := &Reply{
		StateCode: 220,
		Text: "mail.aaronkir.xyz",
	}
	return rep.String()
}

// todo
func (rer *Receiver) ReadCommand() *Command {
	// line:
	// StateCode Param
	//line := rer.bfr.ReadString('\n')
	return nil
}

// todo
func (rer *Receiver) WriteReply(com *Command) {
	if com == nil {
		log.Println("com is nil")
		return
	}

	n, err := rer.bfw.WriteString(com.String())
	if err != nil {
		log.Fatalf("write count: %d. err: %s", n, err)
	}

	rer.bfw.Flush()
}
