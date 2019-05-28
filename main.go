package main

import (
	"log"
	"smtp/module"
	"smtp/web"
	"time"
)

// todo: 生成唯一邮件地址.

func main() {
	s := module.SMTPServer{}
	go s.ListenAndAccept()
	go web.Handle()

	log.Println("sleep...")
	time.Sleep(10 * time.Minute)
}
