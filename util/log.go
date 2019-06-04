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

	//HTTPLog = log.New(os.Stderr, httpPfix, flags)
	//SMTPLog = log.New(os.Stderr, smtpPfix, flags)
	foH := &fileOut{
		fileName: "httpLog.log",
	}
	foS := &fileOut{
		fileName: "smtpLog.log",
	}

	HTTPLog = log.New(foH, httpPfix, flags)
	SMTPLog = log.New(foS, smtpPfix, flags)
}

type fileOut struct {
	fileName string
}

func (fo *fileOut) Write(p []byte) (n int, err error) {
	f, err := os.OpenFile(fo.fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	return f.Write(p)
}
