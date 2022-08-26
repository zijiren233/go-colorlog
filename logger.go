package logger

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
)

var (
	log       logger
	backre, _ = regexp.Compile(`^(\d{4}_\d{2}_\d{2}_\d{2}_\d{2}_\d{2})_(.*?).log`)
)

const (
	Debug uint = iota
	Info
	Warning
	Error
	Fatal
	None
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type logger struct {
	levle        uint
	disableprint bool
	fileOBJ      *os.File
	errFileOBJ   *os.File
	message      chan *logmsg
	maxBackLog   uint
}

type logmsg struct {
	levle    uint
	message  string
	now      string
	funcName string
	filename string
	line     int
}

func levleColor(levle uint) string {
	switch levle {
	case Debug:
		return blue
	case Info:
		return green
	case Warning:
		return yellow
	case Error:
		return red
	case Fatal:
		return magenta
	default:
		return red
	}
}

func LevleToInt(s string) uint {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warning":
		return Warning
	case "error":
		return Error
	case "fatal":
		return Fatal
	case "none":
		return None
	default:
		return Debug
	}
}

func IntToLevle(i uint) string {
	switch i {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	case None:
		return "NONE"
	default:
		return "DEBUG"
	}
}

func init() {
	log = logger{levle: Info, disableprint: false, maxBackLog: 3, message: make(chan *logmsg, 50)}
	log.fileInit()
	go log.backWriteLog()
}

func DisableLogPrint() {
	log.disableprint = true
}

func EnableLogPrint() {
	log.disableprint = false
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

func (l *logger) log(levle uint, format string, a ...interface{}) {
	if l.levle <= levle {
		var log logmsg
		if l.levle == Debug {
			filename, funcName, line := getInfo()
			log = logmsg{
				levle:    levle,
				message:  fmt.Sprintf(format, a...),
				now:      time.Now().Format("[2006-01-02 15:04:05]"),
				funcName: funcName,
				filename: filename,
				line:     line,
			}
		} else {
			log = logmsg{
				levle:   levle,
				message: fmt.Sprintf(format, a...),
				now:     time.Now().Format("[2006-01-02 15:04:05]"),
			}
		}
		select {
		case l.message <- &log:
		default:
		}
	}
}

func SetLogLevle(levle uint) {
	log.levle = levle
}

func (l *logger) backWriteLog() {
	var msgtmp *logmsg
	stdout := colorable.NewColorableStdout()
	for {
		l.deleteBacLog()
		msgtmp = <-l.message
		l.backupLog()
		if l.levle == Debug {
			fmt.Fprintf(l.fileOBJ, "%s [%s] [%s|%s|%d] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.filename, msgtmp.funcName, msgtmp.line, msgtmp.message)
		} else {
			fmt.Fprintf(l.fileOBJ, "%s [%s] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.message)
		}
		if !l.disableprint {
			fmt.Fprintf(stdout, "%s |%s %s %s| %s\n", msgtmp.now, levleColor(msgtmp.levle), IntToLevle(msgtmp.levle), reset, msgtmp.message)
		}
		if msgtmp.levle >= Error {
			l.backupErrLog()
			if l.levle == Debug {
				fmt.Fprintf(l.errFileOBJ, "%s [%s] [%s|%s|%d] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.filename, msgtmp.funcName, msgtmp.line, msgtmp.message)
			} else {
				fmt.Fprintf(l.errFileOBJ, "%s [%s] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.message)
			}
		}
	}
}

func (l *logger) deleteBacLog() {
	files, err := os.ReadDir("./logs")
	if err != nil {
		Errorf("索引 logs 文件失败!")
		return
	}
	var (
		logCount uint
		errCount uint
		logTime  = make(map[time.Time]string)
		errTime  = make(map[time.Time]string)
	)
	for _, file := range files {
		s := backre.FindStringSubmatch(file.Name())
		if len(s) != 3 {
			continue
		}
		switch s[2] {
		case "log":
			logCount++
			t, _ := time.Parse("2006_01_02_15_04_05", s[1])
			logTime[t] = file.Name()
		case "err":
			errCount++
			t, _ := time.Parse("2006_01_02_15_04_05", s[1])
			errTime[t] = file.Name()
		default:
		}
	}
	var (
		earliestLog = time.Now()
		earliestErr = time.Now()
	)
	if logCount > l.maxBackLog {
		for k := range logTime {
			if k.After(earliestLog) {
				earliestLog = k
			}
		}
		os.Remove("./logs/" + logTime[earliestLog])
	}
	if errCount > l.maxBackLog {
		for k := range logTime {
			if k.After(earliestErr) {
				earliestErr = k
			}
		}
		os.Remove("./logs/" + logTime[earliestErr])
	}
}

func getInfo() (string, string, int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "", "", 0
	}
	return path.Base(file), path.Base(runtime.FuncForPC(pc).Name()), line
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
