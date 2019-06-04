package api

import (
	"encoding/json"
	"net/http"
	"smtp/module"
	"smtp/util"
	"smtp/web/tpldata"
)

func WriteJsonMail(w http.ResponseWriter, r *http.Request) {
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
		addrStr := pfx+tpldata.AddrSuf

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
		Age int
	}{
		Name: "Hello",
		Age: 1,
	}

	data, _ := json.Marshal(st)
	w.Write(data)
}
