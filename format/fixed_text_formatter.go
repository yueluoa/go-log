package format

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	fieldMsgKey   = "_msg"
	fieldLevelKey = "_level"
	fieldTimeKey  = "_time"
	fieldFuncKey  = "_func"
	fieldFileKey  = "_file"
)

type Field struct {
	Key   string
	Value interface{}
}

type fixedTextFormatter struct {
	D []*Field
}

func NewLogFormat() LogFormat {
	return &fixedTextFormatter{
		D: make([]*Field, 0),
	}
}

func (ftf fixedTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	keys := make([]string, 0, len(entry.Data))
	for k, _ := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(ftf.D) == 0 {
		ftf.D = make([]*Field, 0)
	}
	if entry.HasCaller() {
		fieldFunc := entry.Caller.Function
		fieldFile := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		ftf.D = append(ftf.D, &Field{Key: fieldFuncKey, Value: fieldFunc})
		ftf.D = append(ftf.D, &Field{Key: fieldFileKey, Value: fieldFile})
	}

	for _, v := range keys {
		ftf.D = append(ftf.D, &Field{Key: v, Value: entry.Data[v]})
	}
	if entry.Message != "" {
		ftf.D = append(ftf.D, &Field{Key: fieldMsgKey, Value: entry.Message})
	}

	b := &bytes.Buffer{}
	b.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	b.WriteString(" [" + strings.ToUpper(entry.Level.String()) + "] ")
	for _, v := range ftf.D {
		b.WriteString(v.Key)
		b.WriteString("= ")
		b.WriteString(fmt.Sprint(v.Value, " "))
	}
	b.WriteString("\n")

	return b.Bytes(), nil
}

type fieldKey string

type FieldMap map[fieldKey]string

func (f FieldMap) getFieldKey(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}
