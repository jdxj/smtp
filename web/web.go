package web

import (
	"log"
	"net/http"
	"smtp/module"
)

type Server struct {
}

func Handle() {
	http.HandleFunc("/", testHello)
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
