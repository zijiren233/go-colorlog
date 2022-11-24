package colorlog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func (l *logger) backupLog() {
	info, err := l.fileOBJ.Stat()
	if err != nil {
		l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
		return
	}
	// 2M
	if info.Size()+int64(l.logFile.Buffered()) >= 2097152 {
		l.logFile.Flush()
		l.fileOBJ.Close()
		os.Rename(`./logs/log.log`, fmt.Sprint(`./logs/`, time.Now().Format("2006_01_02_15_04_05_log.log")))
		f, err := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
			return
		}
		l.fileOBJ = f
		l.logFile = bufio.NewWriter(f)
	}
}

func (l *logger) backupErrLog() {
	info, err := l.errFileOBJ.Stat()
	if err != nil {
		l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
		return
	}
	if info.Size()+int64(l.errFile.Buffered()) >= 2097152 {
		l.errFile.Flush()
		l.errFileOBJ.Close()
		os.Rename(`./logs/err.log`, fmt.Sprint(`./logs/`, time.Now().Format("2006_01_02_15_04_05_err.log")))
		f, err := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
			return
		}
		l.errFileOBJ = f
		l.errFile = bufio.NewWriter(f)
	}
}
