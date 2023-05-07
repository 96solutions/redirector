package logger

import (
	"os"

	"github.com/lroman242/redirector/domain/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type zapSugarLoggerWrap struct {
	*zap.SugaredLogger
}

func (l *zapSugarLoggerWrap) With(args ...interface{}) service.Logger {
	l.SugaredLogger = l.SugaredLogger.With(args...)

	return l
}

func NewZapLogger(logLevel string, logDir string, logFile string) service.Logger {
	lvl, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	// output logs to stdout
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		lvl,
	)
	// output logs to log file
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   logDir + "/" + logFile,
			MaxSize:    500,
			MaxAge:     30,
			MaxBackups: 5,
			LocalTime:  true,
			Compress:   true,
		}),
		zap.ErrorLevel,
	)
	//TODO: output logs to Sentry
	//TODO: output logs to Kafka
	//TODO: output logs to NewRelic

	// merge cores
	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(core, zapOptions()...)

	return &zapSugarLoggerWrap{
		logger.Sugar(),
	}
}

func zapOptions() []zap.Option {
	return []zap.Option{zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller()}
}
