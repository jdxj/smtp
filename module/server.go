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

	util.SMTPLog.Println("SMTP Server started!")
	for {
		conn, err := l.Accept()
		if err != nil {
			util.SMTPLog.Println(err)
			continue
		}

		util.SMTPLog.Println("Receive a connection, Remote addr: ", conn.RemoteAddr())
		util.IPLog.Println(conn.RemoteAddr())

		task := func() error {
			rer := NewReceiver(conn)
			rer.Start()
			return nil
		}
		util.WorkerPool.Submit(task)
	}
}
