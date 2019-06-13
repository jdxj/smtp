package util

import (
	"log"
	"os"
)

var (
	HTTPLog *log.Logger
	SMTPLog *log.Logger
	IPLog   *log.Logger
)

const (
	httpPfix = "[HTTPServer]"
	smtpPfix = "[SMTPServer]"
	ipPfix   = "[IPLog]"

	flags = log.Ldate | log.Ltime | log.Lshortfile
)

func init() {
	HTTPLog = log.New(newFileOUt("httpLog.log"), httpPfix, flags)
	SMTPLog = log.New(newFileOUt("smtpLog.log"), smtpPfix, flags)
	IPLog = log.New(newFileOUt("ipLog.log"), ipPfix, flags)
}

func newFileOUt(path string) *fileOut {
	fo := &fileOut{
		fileName: path,
	}
	return fo
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
