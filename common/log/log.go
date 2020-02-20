package log

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

var logger *logrus.Logger

func init() {
	logger = &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func Infof(f string, v ...interface{}) {
	logger.Infof(f, v...)
}

func Warnf(f string, v ...interface{}) {
	logger.Warnf(f, v...)
}

func Errorf(err error, f string, v ...interface{}) {
	logger.WithError(err).Errorf(f, v...)
}

func Fatalf(err error, f string, v ...interface{}) {
	logger.WithError(err).Errorf(f, v...)
}

func Panic(e interface{}) {
	logger.Panic(e)
}
