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

const (
	L_Debug uint = iota
	L_Info
	L_Warning
	L_Error
	L_Fatal
	L_None
)

type logger struct {
	levle      uint
	logPrint   bool
	fileOBJ    *os.File
	errFileOBJ *os.File
	logFile    *bufio.Writer
	errFile    *bufio.Writer
	stdout     io.Writer
	message    chan *logmsg
	maxBackLog uint
}

type logmsg struct {
	levle    uint
	message  string
	now      string
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

func SetLogLevle(levle uint) {
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
		l.logprint(&logmsg{levle: L_Fatal, message: "Open log.log err: " + err.Error(), now: time.Now().Format(timeFormate)})
		panic(err)
	}
	ef, err := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		l.logprint(&logmsg{levle: L_Fatal, message: "Open err.log err: " + err.Error(), now: time.Now().Format(timeFormate)})
		panic(err)
	}
	l.fileOBJ = f
	l.logFile = bufio.NewWriter(f)
	l.errFileOBJ = ef
	l.errFile = bufio.NewWriter(ef)
}

func logfd(levle uint, format string, a ...interface{}) {
	if log.levle <= levle {
		var msg logmsg
		if log.levle == L_Debug {
			filename, funcName, line := getInfo()
			msg = logmsg{
				levle:    levle,
				message:  fmt.Sprintf(format, a...),
				now:      time.Now().Format(timeFormate),
				funcName: funcName,
				filename: filename,
				line:     line,
			}
		} else {
			msg = logmsg{
				levle:   levle,
				message: fmt.Sprintf(format, a...),
				now:     time.Now().Format(timeFormate),
			}
		}
		select {
		case log.message <- &msg:
		default:
		}
		log.logprint(&msg)
	}
}

func logd(levle uint, a ...interface{}) {
	if log.levle <= levle {
		var msg logmsg
		if log.levle == L_Debug {
			filename, funcName, line := getInfo()
			msg = logmsg{
				levle:    levle,
				message:  fmt.Sprint(a...),
				now:      time.Now().Format(timeFormate),
				funcName: funcName,
				filename: filename,
				line:     line,
			}
		} else {
			msg = logmsg{
				levle:   levle,
				message: fmt.Sprint(a...),
				now:     time.Now().Format(timeFormate),
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
		l.writeToFile(msgtmp, l.logFile)
		if msgtmp.levle >= L_Error {
			l.backupErrLog()
			l.writeToFile(msgtmp, l.errFile)
		}
	}
	l.logFile.Flush()
	l.fileOBJ.Close()
	l.errFile.Flush()
	l.errFileOBJ.Close()
}
