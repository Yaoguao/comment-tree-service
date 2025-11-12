package logr

import (
	"github.com/sirupsen/logrus"
	"io"
)

func InitLogrusLog(env string, w io.Writer) *logrus.Logger {

	log := &logrus.Logger{
		Out:       w,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
	}

	switch env {
	case "dev":
		log.SetLevel(logrus.TraceLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
