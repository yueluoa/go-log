package go_log

import (
	"github.com/iiiang/go-log/level"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	LogConfig
	LogField
}

type LogConfig interface {
	SetFormatter(f logrus.Formatter)
	SetReportElastic(isReport bool)
	SetReportCaller(isReport bool)
	SetOutLevel(outLevel OutStatus) // 日志输出级别, 1: 终端, 2: 文件输出, 3: 同时输出
	SetLevel(level level.Level)
}

type LogField interface {
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
