package colorlog

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"time"
)

var backre, _ = regexp.Compile(`^(\d{4}_\d{2}_\d{2}_\d{2}_\d{2}_\d{2})_(.*?).log`)

func (l *logger) deleteBacLog() {
	files, err := os.ReadDir("./logs")
	if err != nil {
		Errorf("Index logs file err: %v", err)
		return
	}
	if len(files) <= int(l.maxBackLog) {
		return
	}
	logCount, errCount, logTime, errTime := findLogFile(&files)
	if logCount > l.maxBackLog {
		delectLog(logTime, "log")
	}
	if errCount > l.maxBackLog {
		delectLog(errTime, "err")
	}
}

func findLogFile(files *[]fs.DirEntry) (uint, uint, time.Time, time.Time) {
	var (
		logCount uint
		errCount uint
		logTime  = time.Now()
		errTime  = time.Now()
	)
	for _, file := range *files {
		s := backre.FindStringSubmatch(file.Name())
		if len(s) != 3 {
			continue
		}
		switch s[2] {
		case "log":
			logCount++
			t, _ := time.Parse("2006_01_02_15_04_05", s[1])
			if t.Before(logTime) {
				logTime = t
			}
		case "err":
			errCount++
			t, _ := time.Parse("2006_01_02_15_04_05", s[1])
			if t.Before(errTime) {
				errTime = t
			}
		default:
		}
	}
	return logCount, errCount, logTime, errTime
}

func delectLog(fileTime time.Time, types string) {
	os.Remove(fmt.Sprintf("./logs/%s_%s.log", fileTime.Format("2006_01_02_15_04_05"), types))
}
