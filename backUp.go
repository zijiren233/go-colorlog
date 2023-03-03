package colorlog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func (l *logger) backupLog() error {
	// info, err := l.fileOBJ.Stat()
	// if err != nil {
	// 	l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
	// 	return
	// }
	// 2M
	if l.fileSize+int64(l.logFile.Buffered()) >= 2097152 {
		l.logFile.Flush()
		l.fileOBJ.Close()
		os.Rename(`./logs/log.log`, fmt.Sprint(`./logs/`, time.Now().Format("2006_01_02_15_04_05_log.log")))
		f, err := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
			return err
		}
		l.fileOBJ = f
		l.logFile = bufio.NewWriter(f)
		l.fileSize = 0
	}
	return nil
}

func (l *logger) backupErrLog() error {
	// info, err := l.errFileOBJ.Stat()
	// if err != nil {
	// 	l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
	// 	return
	// }
	if l.errFileSize >= 2097152 {
		l.errFileOBJ.Close()
		os.Rename(`./logs/err.log`, fmt.Sprint(`./logs/`, time.Now().Format("2006_01_02_15_04_05_err.log")))
		f, err := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			l.logprint(&logmsg{levle: L_Fatal, message: "Backup ErrLog Err: " + err.Error(), now: time.Now().Format(timeFormate)})
			return err
		}
		l.errFileOBJ = f
		l.errFileSize = 0
	}
	return nil
}
