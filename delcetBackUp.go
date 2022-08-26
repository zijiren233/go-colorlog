package log

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

var backre, _ = regexp.Compile(`^(\d{4}_\d{2}_\d{2}_\d{2}_\d{2}_\d{2})_(.*?).log`)

func (l *logger) deleteBacLog() {
	files, err := os.ReadDir("./logs")
	if err != nil {
		Errorf("索引 logs 文件失败!")
		return
	}
	if len(files) <= int(l.maxBackLog) {
		return
	}
	var (
		logCount uint
		errCount uint
		logTime  = time.Now()
		errTime  = time.Now()
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
	if logCount > l.maxBackLog {
		delectLog(errTime, "log")
	}
	if errCount > l.maxBackLog {
		delectLog(errTime, "err")
	}
}

func delectLog(fileTime time.Time, types string) {
	os.Remove(fmt.Sprintf("./logs/%s_%s.log", fileTime.Format("2006_01_02_15_04_05"), types))
}
