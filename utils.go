package colorlog

import (
	"os"
	"path"
	"runtime"
	"strings"
)

func levleColor(levle uint) string {
	switch levle {
	case debug:
		return blue
	case info:
		return green
	case warning:
		return yellow
	case err:
		return red
	case fatal:
		return magenta
	default:
		return red
	}
}

func LevleToInt(s string) uint {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return debug
	case "info":
		return info
	case "warning":
		return warning
	case "error":
		return err
	case "fatal":
		return fatal
	case "none":
		return none
	default:
		return debug
	}
}

func IntToLevle(i uint) string {
	switch i {
	case debug:
		return "DEBUG"
	case info:
		return "INFO"
	case warning:
		return "WARNING"
	case err:
		return "ERROR"
	case fatal:
		return "FATAL"
	case none:
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
