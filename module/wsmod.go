package module

import "github.com/gorilla/websocket"

type UserInfo struct {
	MailAddr string
	Mail     *MailMsg
	Upgrader *websocket.Upgrader
	WSConn   *websocket.Conn
}

func (userInfo *UserInfo) PushMail() {
	data, _ := userInfo.Mail.Json()
	userInfo.WSConn.WriteMessage(websocket.TextMessage, data)
}

func (userInfo *UserInfo) PushMailAddr() {
	data := []byte(userInfo.MailAddr)
	userInfo.WSConn.WriteMessage(websocket.TextMessage, data)
}
