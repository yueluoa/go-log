package go_log

import (
	"github.com/iiiang/go-log/level"
	"github.com/sirupsen/logrus"
)

type Option interface {
	apply(*Log)
}

type logOption struct {
	f func(*Log)
}

func (flo *logOption) apply(log *Log) {
	flo.f(log)
}

func NewLogOption(f func(*Log)) *logOption {
	return &logOption{
		f: f,
	}
}

func WithLevel(lv string) Option {
	return NewLogOption(func(l *Log) {
		l.level = level.ParseLevel(lv)
	})
}

func WithFormatter(f logrus.Formatter) Option {
	return NewLogOption(func(l *Log) {
		l.formatter = f
	})

}

func WithReportElastic(isReport bool) Option {
	return NewLogOption(func(l *Log) {
		l.reportElastic = isReport
	})
}

func WithElasticIndex(index string) Option {
	return NewLogOption(func(l *Log) {
		l.elasticIndex = index
	})
}

func WithReportCaller(isReport bool) Option {
	return NewLogOption(func(l *Log) {
		l.reportCaller = isReport
	})
}

func WithOutLevel(out OutStatus) Option {
	return NewLogOption(func(l *Log) {
		l.outLevel = out
	})
}
