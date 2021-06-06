package hooks

import (
	"runtime"
	"strings"
	"sync"

	"github.com/iiiang/go-log/level"

	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
)

type redirectHook struct {
	logLevels     []level.Level
	elasticHook   *elogrus.ElasticHook
	reportCaller  bool
	reportElastic bool // 是否发送es
}

func NewRedirectHook() *redirectHook {
	return &redirectHook{}
}

func (h *redirectHook) Levels() []logrus.Level {
	lvs := make([]logrus.Level, 0)
	if len(h.logLevels) == 0 {
		h.logLevels = level.AllLevels
	}
	for _, lv := range h.logLevels {
		lvs = append(lvs, logrus.Level(lv))
	}

	return lvs
}

func (h *redirectHook) Fire(entry *logrus.Entry) error {
	if h.reportCaller {
		entry.Caller = getCaller()
	}

	return nil
}

func (h *redirectHook) SetReportCaller(isReport bool) {
	h.reportCaller = isReport
}

const (
	maximumCallerDepth int = 25
	knownLogFrames     int = 7
)

var (
	logPackage         string
	callerInitOnce     sync.Once
	minimumCallerDepth int
)

func getCaller() *runtime.Frame {
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logPackage = getPackageName(funcName)
				break
			}
		}
		minimumCallerDepth = knownLogFrames
	})

	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		if pkg != logPackage {
			return &f
		}
	}

	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
