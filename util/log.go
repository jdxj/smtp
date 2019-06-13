package util

import (
	"log"
	"os"
	"time"
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
	flags    = log.Ldate | log.Ltime | log.Lshortfile

	TimeFormat = "2006-01-02_15:04:05"
	logSize    = 50 * 1024 * 1024 // 50M
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

	defer fo.rename()
	defer f.Close()

	return f.Write(p)
}

func (fo *fileOut) rename() {
	info, err := os.Stat(fo.fileName)
	if err != nil {
		log.Panicln(err)
	}
	if info.Size() <= logSize {
		return
	}

	err = os.Rename(fo.fileName, fo.fileName+"."+NowDateTime())
	if err != nil {
		log.Panicln(err)
	}
}

func NowDateTime() string {
	return time.Now().Format(TimeFormat)
}
