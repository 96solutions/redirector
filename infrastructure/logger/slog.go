package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/lroman242/redirector/config"
)

// NewLogger creates a new slog.Logger instance.
func NewLogger(conf *config.LoggerConf) *slog.Logger {
	// write to stdout by default
	//TODO: add ability to provide writer
	var w io.Writer = os.Stdout

	var l slog.Level
	if err := l.UnmarshalText([]byte(conf.Level)); err != nil {
		l = slog.LevelInfo
	}

	options := &slog.HandlerOptions{
		AddSource: conf.AddSource,
		Level:     l,
	}

	var h slog.Handler = slog.NewTextHandler(w, options)

	if conf.IsJSON {
		h = slog.NewJSONHandler(w, options)
	}

	logger := slog.New(h)

	if conf.ReplaceDefault {
		slog.SetDefault(logger)
	}

	return logger
}

// WithDefaultAttrs returns logger with default attributes.
func WithDefaultAttrs(logger *slog.Logger, attrs ...slog.Attr) *slog.Logger {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

// ErrAttr prepares slog.Attr which contains error data.
func ErrAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// "unique" key for the context to store logger instance.
type ctxLoggerKey struct{}

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

// loggerFromContext returns logger from context.
func loggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}
