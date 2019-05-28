package module

import (
	"log"
	"net"
)

func init() {
	log.SetPrefix("[SMTPServer]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type SMTPServer struct {
	conns chan net.Conn
}

func (ss *SMTPServer) ListenAndAccept() {
	l, err := net.Listen("tcp", ":25")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		log.Println("等待连接!")
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			log.Println("接受新连接!")
			rer := NewReceiver(conn)
			rer.Start()
			log.Printf("server store addr: %p", Store)
			log.Printf("server store list addr: %p", Store.list)
			log.Printf("server store len: %d", Store.Len())
		}()
	}
}
