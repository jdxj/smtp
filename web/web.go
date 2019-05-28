package web

import (
	"log"
	"net/http"
)

type Server struct {
}

func Handle() {
	http.HandleFunc("/", testHello)
	http.HandleFunc("/hel", testHello)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func testHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request URL: ", r.RequestURI)
	w.Write([]byte("hello world!"))
}
