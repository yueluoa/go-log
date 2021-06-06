package hooks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/iiiang/go-log/level"

	"github.com/iiiang/go-log/format"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

var (
	ElasticIndex    = "testlog"
	ElasticHost     = "http://zhigufengyu.cn:9200"
	ElasticIndexErr = fmt.Errorf("%v", "索引不存在,es无法搜集该信息,请先创建~")
)

type IndexNameFunc func() string

type fireFunc func(entry *logrus.Entry, hook *ElasticHook) error

type ElasticHook struct {
	client    *elastic.Client
	index     IndexNameFunc
	levels    []level.Level
	ctx       context.Context
	ctxCancel context.CancelFunc
	fireFunc  fireFunc
}

type message struct {
	Timestamp string `json:"@timestamp"`
	Message   string `json:"Message,omitempty"`
	Level     string `json:"Level,omitempty"`
}

func NewElastic(index string) *ElasticHook {
	return newElastic(index)
}

func newElastic(index string) *ElasticHook {
	client, err := elastic.NewClient(elastic.SetSniff(false),
		elastic.SetURL(ElasticHost))
	if err != nil {
		fmt.Println("elastic.NewClient err: ", err)
		return nil
	}
	elasticHook, err := NewElasticHook(client, index)
	if err != nil {
		fmt.Println("elogrus.NewElasticHook err: ", err)
		return nil
	}

	return elasticHook
}

func NewElasticHook(client *elastic.Client, index string) (*ElasticHook, error) {
	return newElasticHook(client, index)
}

func newElasticHook(client *elastic.Client, index string) (*ElasticHook, error) {
	es := &ElasticHook{
		client: client,
		index:  func() string { return index },
	}
	err := es.initHookFuncAndFireFunc(syncFireFunc)

	return es, err
}

func (es *ElasticHook) initHookFuncAndFireFunc(fireFunc fireFunc) error {

	ctx, cancel := context.WithCancel(context.TODO())

	exists, err := es.client.IndexExists(es.index()).Do(ctx)
	if err != nil {
		cancel()
		return err
	}
	if !exists {
		return ElasticIndexErr
	}

	es.ctx = ctx
	es.ctxCancel = cancel
	es.fireFunc = fireFunc

	return nil
}

func (es *ElasticHook) Fire(entry *logrus.Entry) error {
	return es.fireFunc(entry, es)
}

func (es *ElasticHook) Levels() []logrus.Level {
	lvs := make([]logrus.Level, 0)
	if len(es.levels) == 0 {
		es.levels = level.AllLevels
	}
	for _, lv := range es.levels {
		lvs = append(lvs, logrus.Level(lv))
	}

	return lvs
}

func createMessage(entry *logrus.Entry) *message {
	lv := entry.Level.String()

	f := format.NewLogFormat()
	byteData, _ := f.Format(entry)

	msg := string(byteData)
	msg = strings.Replace(msg, "\n", "", -1)

	return &message{
		Timestamp: entry.Time.Format(time.RFC3339),
		Message:   msg,
		Level:     strings.ToUpper(lv),
	}
}

func syncFireFunc(entry *logrus.Entry, hook *ElasticHook) error {
	_, err := hook.client.
		Index().
		Index(hook.index()).
		Type("log").
		BodyJson(createMessage(entry)).
		Do(hook.ctx)

	return err
}
