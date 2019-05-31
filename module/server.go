package module

import (
	"net"
	"smtp/util"
)

type SMTPServer struct {
	conns chan net.Conn
}

func (ss *SMTPServer) ListenAndAccept() {
	l, err := net.Listen("tcp", ":25")
	if err != nil {
		util.SMTPLog.Fatal(err)
	}
	defer l.Close()

	for {
		util.SMTPLog.Println("等待连接!")
		conn, err := l.Accept()
		if err != nil {
			util.SMTPLog.Println(err)
			continue
		}

		go func() {
			util.SMTPLog.Println("接受新连接!")
			rer := NewReceiver(conn)
			rer.Start()
		}()
	}
}
