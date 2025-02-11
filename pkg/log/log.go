package log

import (
	"context"
	"sync"
)

var (
	log     Logger        = nil
	factory func() Logger = func() Logger {
		return NewSlogAdapter(SlogAdapterOpts{
			Level:      LevelInfo,
			FormatJson: false,
		})
	}
	mutex sync.Mutex
)

type Level uint

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, err error)
	SetLevel(l Level) error
}

func Log() Logger {
	mutex.TryLock()
	defer mutex.Unlock()
	if log == nil {
		log = factory()
	}
	return log
}

func SetLogger(lg Logger) {
	log = lg
}

func SetLoggerFactory(f func() Logger) {
	factory = f
}

func NewLogger(name string) Logger {
	return factory()
}
