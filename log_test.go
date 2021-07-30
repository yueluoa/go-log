package go_log

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

const path = "test/ddd"

func Test_LLog(t *testing.T) {
	opts := []Option{
		WithLevel("info"),
		WithOutLevel(Terminal),
	}

	var l Logger
	l = NewLog(path, opts...)
	l.Info("今天", "星期几？")
}
