package util

import (
	"log"
	"os"
)

var HTTPLog *log.Logger
var SMTPLog *log.Logger

func init() {
	httpPfix := "[HTTPServer]"
	smtpPfix := "[SMTPServer]"

	flags := log.Ldate | log.Ltime | log.Lshortfile

	HTTPLog = log.New(os.Stderr, httpPfix, flags)
	SMTPLog = log.New(os.Stderr, smtpPfix, flags)
}
