package main

import (
	"./module"
	"log"
	"time"
)
import "./web"

func main() {
	s := module.SMTPServer{}
	go s.ListenAndAccept()
	go web.Handle()

	log.Println("sleep...")
	time.Sleep(10 * time.Minute)
}
