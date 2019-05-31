package main

import (
	"log"
	"smtp/module"
	"smtp/web"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go StartSmtp(wg)
	go StartHttp(wg)

	wg.Wait()
	log.Println("Server quits actively!")
}

func StartSmtp(wg *sync.WaitGroup) {
	smtpServer := &module.SMTPServer{}
	smtpServer.ListenAndAccept()

	wg.Done()
}

func StartHttp(wg *sync.WaitGroup) {
	httpServer := &web.HTTPServer{}
	httpServer.Handle()

	wg.Done()
}
