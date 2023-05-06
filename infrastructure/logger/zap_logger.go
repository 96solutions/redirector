package logger

import (
	"os"

	"github.com/lroman242/redirector/domain/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapSugarLoggerWrap struct {
	*zap.SugaredLogger
}

func (l *zapSugarLoggerWrap) With(args ...interface{}) service.Logger {
	l.SugaredLogger = l.SugaredLogger.With(args...)

	return l
}

func NewZapLogger(logLevel string) service.Logger {
	lvl, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		lvl,
	)
	//fileCore := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	//	zapcore.AddSync(),
	//	zap.ErrorLevel,
	//)

	core := zapcore.NewTee(consoleCore)

	logger := zap.New(core, zapOptions()...)

	return &zapSugarLoggerWrap{
		logger.Sugar(),
	}
}

func zapOptions() []zap.Option {
	return []zap.Option{zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller()}
}
