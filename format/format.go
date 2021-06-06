package format

import (
	"github.com/sirupsen/logrus"
)

type LogFormat interface {
	Format(entry *logrus.Entry) ([]byte, error)
}
