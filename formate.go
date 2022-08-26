package colorlog

import (
	"fmt"
	"os"
)

func (l *logger) writeToFile(msgtmp *logmsg, FileOBJ *os.File) {
	if l.levle == Debug {
		fmt.Fprintf(FileOBJ, "%s [%s] [%s|%s|%d] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.filename, msgtmp.funcName, msgtmp.line, msgtmp.message)
	} else {
		fmt.Fprintf(FileOBJ, "%s [%s] %s\n", msgtmp.now, IntToLevle(msgtmp.levle), msgtmp.message)
	}
}

func (l *logger) logprint(msgtmp *logmsg) {
	if l.logPrint {
		if l.logColor {
			fmt.Fprintf(colorStdout, "%s |%s %s %s| %s\n", msgtmp.now, levleColor(msgtmp.levle), IntToLevle(msgtmp.levle), reset, msgtmp.message)
		} else {
			fmt.Fprintf(NonColorStdout, "%s |%s %s %s| %s\n", msgtmp.now, levleColor(msgtmp.levle), IntToLevle(msgtmp.levle), reset, msgtmp.message)
		}
	}
}
