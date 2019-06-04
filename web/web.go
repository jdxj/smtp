package web

import (
	"html/template"
	"net/http"
	"smtp/util"
	"smtp/web/api"
	"smtp/web/tpl"
	"smtp/web/tpldata"
)

type HTTPServer struct {
}

func (s *HTTPServer) Handle() {
	// todo: 路径匹配规则需要研究
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	http.HandleFunc("/favicon.ico", Favicon)
	http.HandleFunc("/mail", api.WriteJsonMail)
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

// Favicon 用于重定向到图标 url
func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/favicon.ico", 301)
}
