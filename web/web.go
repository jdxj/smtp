package web

import (
	"html/template"
	"log"
	"net/http"
	"smtp/module"
	"smtp/util"
	"smtp/web/tpl"
	"smtp/web/tpldata"
)

type Server struct {
}

func Handle() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/favicon.ico", Favicon)
	http.HandleFunc("/hel", testHello)
	err := http.ListenAndServe(":8025", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func testHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request URL: ", r.RequestURI)
	cookie, err := r.Cookie("kkk")
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "kkk", Value: "vvv"})
	} else {
		log.Printf("a cookie: %s\n", cookie.Value)
	}

	mailMsg, ok := module.Store.Get().(*module.MailMsg)
	if !ok {
		w.Write([]byte("no data!"))
		return
	}
	w.Write([]byte(mailMsg.String()))
}

func Index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.New("index").Parse(tpl.Index)
	if err != nil {
		w.Write([]byte("no data1!"))
		return
	}

	data := tpldata.MailMod{
		Addr: util.IDGen.GetID() + tpldata.AddrSuf,
	}
	err = temp.Execute(w, data)
	if err != nil {
		log.Println(err)
		w.Write([]byte("no data2!"))
		return
	}
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("no favicon!"))
}
