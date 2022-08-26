package log

import (
	"fmt"
	"os"
	"time"
)

func (l *logger) backupLog() {
	file, _ := l.fileOBJ.Stat()
	// 2M
	if file.Size() >= 2097152 {
		l.fileOBJ.Close()
		os.Rename(`./logs/log.log`, fmt.Sprint(`./logs/`, time.Now().Format("2006_01_02_15_04_05_log.log")))
		f, _ := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		l.fileOBJ = f
	}
}

func (l *logger) backupErrLog() {
	file, _ := l.errFileOBJ.Stat()
	if file.Size() >= 2097152 {
		l.errFileOBJ.Close()
		os.Rename(`./logs/err.log`, fmt.Sprint(`./logs/`, time.Now().Format("2006_01_02_15_04_05_err.log")))
		f, _ := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		l.errFileOBJ = f
	}
}
