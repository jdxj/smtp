package web

import (
	"html/template"
	"net/http"
	"smtp/module"
	"smtp/util"
	"smtp/web/tpl"
	"smtp/web/tpldata"
)

type HTTPServer struct {
}

func (s *HTTPServer) Handle() {
	// todo: 路径匹配规则需要研究
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	http.HandleFunc("/favicon.ico", Favicon)
	http.HandleFunc("/helo", testHello)
	http.HandleFunc("/mail", GetMail)
	http.HandleFunc("/", testHello)

	util.HTTPLog.Println("Http server started!")
	err := http.ListenAndServe(":8025", nil)
	if err != nil {
		util.HTTPLog.Fatalln(err)
	}
}

func testHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("prefix: /, hello world!"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.New("index").Parse(tpl.IndexTpl)
	if err != nil {
		w.Write([]byte("no data1!"))
		return
	}

	data := tpldata.MailMod{
		Addr: util.IDGen.GetID() + tpldata.AddrSuf,
	}
	err = temp.Execute(w, data)
	if err != nil {
		util.HTTPLog.Println(err)
		w.Write([]byte("no data2!"))
		return
	}
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/favicon.ico", 301)
}

func GetMail(w http.ResponseWriter, r *http.Request) {
	addrpfix, ok := module.Store.M.Load(r.RemoteAddr)
	if ok { // 找到 user 标识
		addrStr := addrpfix.(string)
		mailMsgI, ok := module.Store.M.Load(addrStr)
		if ok {
			mailMsg := mailMsgI.(*module.MailMsg)
			mailMod := tpldata.MailMod{
				Addr:    mailMsg.ToAddr(),
				Content: mailMsg.String(),
			}

			temp, err := template.New("mailtpl").Parse(tpl.MailTpl)
			if err != nil {
				util.HTTPLog.Println(err)
				w.Write([]byte("internal error"))
				return
			}
			temp.Execute(w, mailMod)
		} else {
			w.Write([]byte("not receive mail!\n"))
			w.Write([]byte("mail addr: " + addrStr))
		}
	} else {
		pfx := util.IDGen.GetID()
		module.Store.M.Store(r.RemoteAddr, pfx+tpldata.AddrSuf)
		module.Store.DelUser(util.Dur, r.RemoteAddr)
		w.Write([]byte("mail addr: " + pfx + tpldata.AddrSuf))
	}
}
