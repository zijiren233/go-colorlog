package colorlog

import (
	"bufio"
	"fmt"
)

func (l *logger) writeToFile(msgtmp *logmsg, file *bufio.Writer) {
	if l.levle == L_Debug {
		fmt.Fprintf(file, "%s [%s] [%s|%s|%d] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.filename, msgtmp.funcName, msgtmp.line, msgtmp.message)
	} else {
		fmt.Fprintf(file, "%s [%s] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.message)
	}
}

func (l *logger) logprint(msgtmp *logmsg) {
	if l.logPrint {
		fmt.Fprintf(l.stdout, "%s |%s %s %s| %s\n", msgtmp.now, levleColor(msgtmp.levle), IntToLevle(msgtmp.levle), reset, msgtmp.message)
	}
}
