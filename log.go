package go_log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/iiiang/go-log/format"
	"github.com/iiiang/go-log/hooks"
	"github.com/iiiang/go-log/level"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const DefaultDir = "./log/default"

const MaxLogFileCount = 8 // 文件最大保存份数

type OutStatus uint8

const (
	Unknown OutStatus = iota
	Terminal
	File
	TerminalAndFile
)

type Log struct {
	log           *logrus.Logger
	level         level.Level // 日志级别
	formatter     logrus.Formatter
	outLevel      OutStatus
	reportCaller  bool
	reportElastic bool // 输出到es
	elasticIndex  string
	path          string // 文件路径
	fileName      string
}

func New(path string) Logger {
	return newLog(path, loadConfig()...)
}

func NewLog(path string, opts ...Option) Logger {
	return newLog(path, opts...)
}

func newLog(path string, opts ...Option) *Log {
	l := &Log{}
	l.level = level.InfoLevel
	l.path = path
	l.initLog(opts...)

	return l
}

func (l *Log) initLog(opts ...Option) {
	for _, opt := range opts {
		opt.apply(l)
	}

	l.log = logrus.New()
	l.fileName = l.level.String()
	lv := logrus.Level(l.level)
	l.log.SetLevel(lv)

	if l.formatter == nil {
		l.SetFormatter(format.NewLogFormat())
	}
	l.log.SetFormatter(l.formatter)

	if l.reportCaller {
		l.log.ReportCaller = l.reportCaller
	}
	if l.reportElastic {
		if l.elasticIndex == "" {
			l.elasticIndex = hooks.ElasticIndex
		}
		if e := hooks.NewElastic(l.elasticIndex); e != nil {
			l.log.AddHook(e)
		}
	}

	l.log.SetOutput(l.loadOut())
}

func (l *Log) SetLevel(level level.Level) {
	l.level = level
}

func (l *Log) SetOutLevel(outLevel OutStatus) {
	l.outLevel = outLevel
}

func (l *Log) SetFormatter(f logrus.Formatter) {
	l.formatter = f
}

func (l *Log) SetReportElastic(isReport bool) {
	l.reportElastic = isReport
}

func (l *Log) SetReportCaller(isReport bool) {
	l.reportCaller = isReport
}

func (l *Log) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Log) DebugF(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Log) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Log) InfoF(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Log) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Log) ErrorF(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Log) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Log) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Log) Print(args ...interface{}) {
	l.log.Print(args...)
}

func (l *Log) loadOut() io.Writer {
	if l.path == "" {
		l.path = DefaultDir
	}
	var out io.Writer
	if l.outLevel == Terminal {
		out = os.Stdout
	}
	if l.outLevel == File ||
		l.outLevel == TerminalAndFile ||
		l.outLevel == Unknown {
		writer, err := initLogFile(l.path, l.fileName)
		if err != nil {
			fmt.Println(err)
			return out
		}
		out = writer
		if l.outLevel == Unknown {
			out = writer
		}
		if l.outLevel == TerminalAndFile {
			out = io.MultiWriter([]io.Writer{writer, os.Stdout}...)
		}
	}

	return out
}

func initLogFile(path, fileName string) (io.Writer, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create dir failed: %v\n", err)
		}
	}
	fileSuffix := time.Now().Format(".2006-01-02")
	filePath := fmt.Sprintf("%s/%s", path, fileName)
	writer, err := rotateLogs.New(
		filePath+fileSuffix,
		rotateLogs.WithRotationCount(MaxLogFileCount),
		rotateLogs.WithRotationTime(time.Hour*24))
	if err != nil {
		return nil, fmt.Errorf("failed to create rotatelogs: %v\n", err)
	}

	return writer, nil
}

// 默认配置
func loadConfig() []Option {
	opts := []Option{
		WithLevel(level.InfoLevel),
		WithReportElastic(true),
		WithOutLevel(File),
	}

	return opts
}
