package logger

import (
	"os"

	"github.com/lroman242/redirector/domain/service"
	log "github.com/sirupsen/logrus"
)

type logrusLogger struct {
	*log.Entry
}

func NewLogrusLogger(logLevel string) service.Logger {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	// Only log the warning severity or above.
	logger.SetLevel(lvl)

	return &logrusLogger{
		log.NewEntry(logger),
	}
}

func (l *logrusLogger) With(args ...interface{}) service.Logger {
	for len(args) >= 2 {
		l.Entry = l.WithField(args[0].(string), args[1])
		args = args[2:]
	}

	return l
}

func (l *logrusLogger) Debug(args ...interface{}) {
	log.Debug(args...)
}
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
func (l *logrusLogger) Debugln(args ...interface{}) {
	log.Debugln(args...)
}
func (l *logrusLogger) Info(args ...interface{}) {
	log.Info(args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
func (l *logrusLogger) Infoln(args ...interface{}) {
	log.Infoln(args...)
}
func (l *logrusLogger) Warn(args ...interface{}) {
	log.Warn(args...)
}
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}
func (l *logrusLogger) Warnln(args ...interface{}) {
	log.Warnln(args...)
}
func (l *logrusLogger) Error(args ...interface{}) {
	log.Error(args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
func (l *logrusLogger) Errorln(args ...interface{}) {
	log.Errorln(args...)
}
func (l *logrusLogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
func (l *logrusLogger) Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}
