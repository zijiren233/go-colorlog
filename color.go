package colorlog

import "fmt"

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

func SGreenf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", green, fmt.Sprintf(format, a...), reset)
}

func SWhitef(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", white, fmt.Sprintf(format, a...), reset)
}

func SYellowf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", yellow, fmt.Sprintf(format, a...), reset)
}
func SRedf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", red, fmt.Sprintf(format, a...), reset)
}
func SBluef(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", blue, fmt.Sprintf(format, a...), reset)
}
func SMagentaf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", magenta, fmt.Sprintf(format, a...), reset)
}
func SCyanf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", cyan, fmt.Sprintf(format, a...), reset)
}
