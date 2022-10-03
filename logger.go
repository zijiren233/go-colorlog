package colorlog

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-colorable"
)

var (
	log         logger
	timeFormate = "[2006-01-02 15:04:05]"

	colorStdout    = colorable.NewColorableStdout()
	nonColorStdout = colorable.NewNonColorable(os.Stdout)
)

const (
	Debug uint = iota
	Info
	Warning
	Error
	Fatal
	None
)

type logger struct {
	levle      uint
	logPrint   bool
	fileOBJ    *os.File
	errFileOBJ *os.File
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
	log = logger{levle: Info, logPrint: true, stdout: colorStdout, maxBackLog: 3, message: make(chan *logmsg, 50)}
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
		fmt.Println("打开日志错误!")
		panic(err)
	}
	ef, err2 := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("打开Err日志错误!")
		panic(err2)
	}
	l.fileOBJ = f
	l.errFileOBJ = ef
}

func (l *logger) log(levle uint, format string, a ...interface{}) {
	if l.levle <= levle {
		var log logmsg
		if l.levle == Debug {
			filename, funcName, line := getInfo()
			log = logmsg{
				levle:    levle,
				message:  fmt.Sprintf(format, a...),
				now:      time.Now().Format(timeFormate),
				funcName: funcName,
				filename: filename,
				line:     line,
			}
		} else {
			log = logmsg{
				levle:   levle,
				message: fmt.Sprintf(format, a...),
				now:     time.Now().Format(timeFormate),
			}
		}
		select {
		case l.message <- &log:
		default:
		}
	}
}

func (l *logger) backWriteLog() {
	var msgtmp *logmsg
	for msgtmp = range l.message {
		l.deleteBacLog()
		l.backupLog()
		l.writeToFile(msgtmp, l.fileOBJ)
		l.logprint(msgtmp)
		if msgtmp.levle >= Error {
			l.backupErrLog()
			l.writeToFile(msgtmp, l.errFileOBJ)
		}
	}
	l.fileOBJ.Close()
	l.errFileOBJ.Close()
}

func Debugf(format string, a ...interface{}) {
	log.log(Debug, format, a...)
}

func Infof(format string, a ...interface{}) {
	log.log(Info, format, a...)
}

func Warringf(format string, a ...interface{}) {
	log.log(Warning, format, a...)
}

func Errorf(format string, a ...interface{}) {
	log.log(Error, format, a...)
}

func Fatalf(format string, a ...interface{}) {
	log.log(Fatal, format, a...)
}
