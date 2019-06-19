package web

import (
	"net/http"
	"smtp/util"
	"smtp/web/api"
)

type HTTPServer struct {
}

func (s *HTTPServer) Handle() {
	http.HandleFunc("/", Welcome)

	// todo: 路径匹配规则需要研究
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	http.HandleFunc("/favicon.ico", Favicon)
	// 已启用, 目前留作记录
	//http.HandleFunc("/mail", api.WriteJsonMail)
	http.HandleFunc("/ws", api.PushJsonMail)

	util.HTTPLog.Println("Http server started!")
	err := http.ListenAndServe(":8025", nil)
	if err != nil {
		util.HTTPLog.Fatalln(err)
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	api.RecordURL(r)
	http.Redirect(w, r, "/static/index.html", 301)
}

// Favicon 用于重定向到图标 url
func Favicon(w http.ResponseWriter, r *http.Request) {
	api.RecordURL(r)
	http.Redirect(w, r, "/static/favicon.ico", 301)
}
