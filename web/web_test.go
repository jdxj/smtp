package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWeb_ListenAndAccept(t *testing.T) {
	http.HandleFunc("/", Hello)
	http.HandleFunc("/hel", Hello)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("abc"))
}
