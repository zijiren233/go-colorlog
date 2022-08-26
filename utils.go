package colorlog

import (
	"os"
	"path"
	"runtime"
	"strings"
)

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

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func getInfo() (string, string, int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "", "", 0
	}
	return path.Base(file), path.Base(runtime.FuncForPC(pc).Name()), line
}
