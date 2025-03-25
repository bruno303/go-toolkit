package log

import (
	"context"
	"sync"
)

var (
	log     Logger              = nil
	factory func(string) Logger = func(name string) Logger {
		return NewSlogAdapter(SlogAdapterOpts{
			Level:      LevelInfo,
			FormatJson: false,
			Source:     name,
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
	Shutdown(context.Context) error
}

func Log() Logger {
	mutex.TryLock()
	defer mutex.Unlock()
	if log == nil {
		log = factory("default")
	}
	return log
}

func SetLogger(lg Logger) {
	log = lg
}

func SetLoggerFactory(f func(string) Logger) {
	factory = f
}

func NewLogger(name string) Logger {
	return factory(name)
}
