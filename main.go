package main

import (
	"log"
	"smtp/module"
	"smtp/web"
	"time"
)

func main() {
	s := module.SMTPServer{}
	go s.ListenAndAccept()
	go web.Handle()

	log.Println("sleep...")
	time.Sleep(20 * time.Minute)
	log.Println("Server quits actively!")
}
