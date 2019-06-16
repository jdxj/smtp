package web

import (
	"fmt"
	"github.com/gorilla/websocket"
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

func TestStr2Bit(t *testing.T) {
	str := "abc"
	fmt.Printf("%b\n", []byte(str))
	str = "1011100111011100011111111011110"
	fmt.Println(len(str))
}

func TestHttpFileServer(t *testing.T) {
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir("static")))
	if err != nil {
		fmt.Println(err)
	}
}

func TestWebSocket(t *testing.T) {
	http.HandleFunc("/ws", handWS)
	http.ListenAndServe(":8080", nil)
}

func handWS(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer wsConn.Close()

	for {
		mt, p, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("message type: ", mt)
		fmt.Println("data: ", p)
	}
}

func TestSlice(t *testing.T) {
	s1 := make([]int, 2)
	s1[0] = 11
	s1[1] = 22

	s2 := s1[:0]
	fmt.Println(s2)
	fmt.Println(s2 == nil)
	fmt.Println(len(s2))
	fmt.Println("--------------")

	var s3 []int
	fmt.Println(s3 == nil)
}