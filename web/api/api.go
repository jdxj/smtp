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

	//i := 0
	//for {
	//	iStr := strconv.Itoa(i)
	//	err = wsConn.WriteMessage(websocket.TextMessage, []byte(iStr))
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	i++
	//	time.Sleep(time.Second)
	//}
	for {
		err = wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second)
	}
}
