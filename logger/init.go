package logger

import (
	"github.com/sirupsen/logrus"
	"strings"
)

func GetLogger(opts Opts) Logger {
	l := logrus.StandardLogger()

	switch level := strings.ToLower(opts.Level); level {
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warning":
		l.SetLevel(logrus.WarnLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	}

	switch format := strings.ToLower(opts.Format); format {
	case "logfmt":
		l.SetFormatter(new(logrus.TextFormatter))
	case "json":
		l.SetFormatter(new(logrus.JSONFormatter))
	}

	return WrapLogrus(l)
}
