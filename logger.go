package go_log

type Logger interface {
	Debug(args ...interface{})
	DebugF(format string, args ...interface{})
	Info(args ...interface{})
	InfoF(format string, args ...interface{})
	Error(args ...interface{})
	ErrorF(format string, args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
	Print(args ...interface{})
}
