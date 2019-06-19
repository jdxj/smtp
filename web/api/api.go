package api

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"smtp/module"
	"smtp/proto/test"
	"smtp/util"
	"smtp/web/tpldata"
	"time"
)

func WriteJsonMail(w http.ResponseWriter, r *http.Request) {
	RecordURL(r)
	addr, ok := module.Store.M.Load(r.RemoteAddr)
	fmt.Println(r.RemoteAddr)
	if ok { // 找到 user 标识
		addrStr := addr.(string)
		mailMsgI, ok := module.Store.M.Load(addrStr)
		if ok {
			mailMsg := mailMsgI.(*module.MailMsg)
			mj, err := mailMsg.Json()
			if err != nil {
				w.WriteHeader(500)

				mj := &module.MailJson{
					Desc: "Internal error!",
				}
				data, _ := json.Marshal(mj)

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(mj)
		} else {
			mj := &module.MailJson{
				Addr: addrStr,
				Desc: "Did not receive mail!",
			}
			data, _ := json.Marshal(mj)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		}
	} else {
		pfx := util.IDGen.GetID()
		addrStr := pfx + tpldata.AddrSuf

		module.Store.M.Store(r.RemoteAddr, addrStr)
		module.Store.DelUser(util.Dur, r.RemoteAddr)

		mj := &module.MailJson{
			Addr: addrStr,
			Desc: "Did not receive mail!",
		}
		data, _ := json.Marshal(mj)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func WriteTestJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "my=af")

	st := &struct {
		Name string
		Age  int
	}{
		Name: "Hello",
		Age:  1,
	}

	data, _ := json.Marshal(st)
	w.Write(data)
}

func RecordURL(r *http.Request) {
	util.HTTPLog.Printf("Method: %s, URL: %s\n", r.Method, r.URL)
}

func TestWebSocket(w http.ResponseWriter, r *http.Request) {
	player := &test.Player{
		ID:   123,
		Name: "321",
	}
	data, _ := proto.Marshal(player)
	upgrader := &websocket.Upgrader{
		//CheckOrigin: func(r *http.Request) bool {
		//	return true
		//},
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer wsConn.Close()

	for {
		err = wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second)
	}
}

func PushJsonMail(w http.ResponseWriter, r *http.Request) {
	pfx := util.IDGen.GetID()
	addrStr := pfx + tpldata.AddrSuf

	upgrader := &websocket.Upgrader{}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.HTTPLog.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Can not user WebSocket!"))
	}

	userInfo := &module.UserInfo{
		MailAddr: addrStr,
		Upgrader: upgrader,
		WSConn:   wsConn,
	}
	userInfo.PushMailAddr()
	module.Store.M.Store(userInfo.MailAddr, userInfo)
	// todo: 释放连接?
	module.Store.DelUser(5*time.Minute, userInfo.MailAddr)
}
