package colorlog

func Debugf(format string, a ...any) {
	logfd(L_Debug, format, a...)
}

func Infof(format string, a ...any) {
	logfd(L_Info, format, a...)
}

func Warringf(format string, a ...any) {
	logfd(L_Warning, format, a...)
}

func Errorf(format string, a ...any) {
	logfd(L_Error, format, a...)
}

func Fatalf(format string, a ...any) {
	logfd(L_Fatal, format, a...)
}

func Debug(a ...any) {
	logd(L_Debug, a...)
}

func Info(a ...any) {
	logd(L_Info, a...)
}

func Warring(a ...any) {
	logd(L_Warning, a...)
}

func Error(a ...any) {
	logd(L_Error, a...)
}

func Fatal(a ...any) {
	logd(L_Fatal, a...)
}
