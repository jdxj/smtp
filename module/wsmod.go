package module

import "github.com/gorilla/websocket"

type WSConn struct {
	ID string
	Conn *websocket.Conn
}

func (wsConn *WSConn) Push(messageType int, p []byte) error {
	return wsConn.Conn.WriteMessage(messageType, p)
}

func (wsConn *WSConn) Close() error {
	return wsConn.Conn.Close()
}
