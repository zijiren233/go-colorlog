package colorlog

import (
	"os"
	"path"
	"runtime"
	"strings"
)

func levleColor(levle uint) string {
	switch levle {
	case L_Debug:
		return blue
	case L_Info:
		return green
	case L_Warning:
		return yellow
	case L_Error:
		return red
	case L_Fatal:
		return magenta
	default:
		return red
	}
}

func LevleToInt(s string) uint {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return L_Debug
	case "info":
		return L_Info
	case "warning":
		return L_Warning
	case "error":
		return L_Error
	case "fatal":
		return L_Fatal
	case "none":
		return L_None
	default:
		return L_Debug
	}
}

func IntToLevle(i uint) string {
	switch i {
	case L_Debug:
		return "DEBUG"
	case L_Info:
		return "INFO"
	case L_Warning:
		return "WARNING"
	case L_Error:
		return "ERROR"
	case L_Fatal:
		return "FATAL"
	case L_None:
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
		return "???", "???", 0
	}
	return path.Base(file), path.Base(runtime.FuncForPC(pc).Name()), line
}
