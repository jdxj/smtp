package main

import "./module"

func main() {
	s := module.SMTPServer{}
	s.ListenAndAccept()
}
