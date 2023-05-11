package colorlog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-colorable"
)

var (
	log         = logger{levle: L_Info, logPrint: true, stdout: colorStdout, maxBackLog: 3, message: make(chan *logmsg, 100)}
	timeFormate = "[2006-01-02 15:04:05]"

	colorStdout    = colorable.NewColorableStdout()
	nonColorStdout = colorable.NewNonColorable(os.Stdout)
)

type LogLevle uint

const (
	L_Debug LogLevle = iota
	L_Info
	L_Warning
	L_Error
	L_Fatal
	L_None
)

type logger struct {
	levle       LogLevle
	logPrint    bool
	fileOBJ     *os.File
	fileSize    int64
	errFileOBJ  *os.File
	errFileSize int64
	logFile     *bufio.Writer
	stdout      io.Writer
	message     chan *logmsg
	maxBackLog  uint
}

type logmsg struct {
	levle    LogLevle
	message  string
	time     time.Time
	funcName string
	filename string
	line     int
}

func init() {
	log.fileInit()
	go log.backWriteLog()
}

func EnableLogPrint(v bool) {
	log.logPrint = v
}

func EnableColor(v bool) {
	if v {
		log.stdout = colorStdout
	} else {
		log.stdout = nonColorStdout
	}
}

func SetTimeFormate(Formate string) {
	timeFormate = Formate
}

func SetLogLevle(levle LogLevle) {
	log.levle = levle
}

func MaxLogFile(count uint) {
	log.maxBackLog = count
}

func (l *logger) fileInit() {
	if !FileExists("./logs") {
		os.Mkdir("./logs", os.ModePerm)
	}
	f, err := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		l.logprint(&logmsg{levle: L_Fatal, message: "Open log.log err: " + err.Error(), time: time.Now()})
		panic(err)
	}
	ef, err := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		l.logprint(&logmsg{levle: L_Fatal, message: "Open err.log err: " + err.Error(), time: time.Now()})
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	l.fileSize = fi.Size()
	l.fileOBJ = f
	efi, err := ef.Stat()
	if err != nil {
		panic(err)
	}
	l.errFileSize = efi.Size()
	l.logFile = bufio.NewWriter(f)
	l.errFileOBJ = ef
}

func logfd(levle LogLevle, format string, a ...interface{}) {
	if log.levle <= levle {
		var msg logmsg
		if log.levle == L_Debug {
			filename, funcName, line := getInfo()
			msg = logmsg{
				levle:    levle,
				message:  fmt.Sprintf(format, a...),
				time:     time.Now(),
				funcName: funcName,
				filename: filename,
				line:     line,
			}
		} else {
			msg = logmsg{
				levle:   levle,
				message: fmt.Sprintf(format, a...),
				time:    time.Now(),
			}
		}
		select {
		case log.message <- &msg:
		default:
		}
		log.logprint(&msg)
	}
}

func logd(levle LogLevle, a ...interface{}) {
	if log.levle <= levle {
		var msg logmsg
		if log.levle == L_Debug {
			filename, funcName, line := getInfo()
			msg = logmsg{
				levle:    levle,
				message:  fmt.Sprint(a...),
				time:     time.Now(),
				funcName: funcName,
				filename: filename,
				line:     line,
			}
		} else {
			msg = logmsg{
				levle:   levle,
				message: fmt.Sprint(a...),
				time:    time.Now(),
			}
		}
		select {
		case log.message <- &msg:
		default:
		}
		log.logprint(&msg)
	}
}

func (l *logger) backWriteLog() {
	for msgtmp := range l.message {
		l.deleteBacLog()
		l.backupLog()
		n, err := l.writeToFile(msgtmp, l.logFile)
		l.fileSize += int64(n)
		if err != nil {
			continue
		}
		if msgtmp.levle >= L_Error {
			l.backupErrLog()
			n, err := l.writeToFile(msgtmp, l.errFileOBJ)
			l.errFileSize += int64(n)
			if err != nil {
				continue
			}
		}
	}
	l.logFile.Flush()
	l.fileOBJ.Close()
	l.errFileOBJ.Close()
}
