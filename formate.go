package colorlog

import (
	"errors"
	"fmt"
	"io"
)

var errFileNil = errors.New("file is nil")

func (l *logger) writeToFile(msgtmp *logmsg, file io.Writer) (n int, err error) {
	if file == nil {
		return 0, errFileNil
	}
	if l.levle == L_Debug || msgtmp.levle >= L_Error {
		return fmt.Fprintf(file, "%s [%s] [%s|%s|%d] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.filename, msgtmp.funcName, msgtmp.line, msgtmp.message)
	}
	return fmt.Fprintf(file, "%s [%s] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.message)
}

func (l *logger) logprint(msgtmp *logmsg) {
	if l.logPrint {
		fmt.Fprintf(l.stdout, "%s |%s %s %s| %s\n", msgtmp.now, levleColor(msgtmp.levle), IntToLevle(msgtmp.levle), reset, msgtmp.message)
	}
}
