package log

import (
	"context"
	"sync"
)

//go:generate go tool mockgen -destination mocks.go -package log . Logger

var (
	log     Logger        = nil
	factory LoggerFactory = func(name string) Logger {
		return NewSlogAdapter(SlogAdapterOpts{
			Level:      LevelInfo,
			FormatJson: false,
			Name:       name,
		})
	}
	mutex sync.Mutex

	logs                = make(map[string]Logger)
	logConfig LogConfig = LogConfig{}
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, err error)
	SetLevel(l Level) error
	Shutdown(context.Context) error
	Name() string
	Level() Level
}

func Log() Logger {
	mutex.Lock()
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
	mutex.Lock()
	defer mutex.Unlock()

	if existingLog, ok := logs[name]; ok {
		loggerPostCreation(existingLog)
		return existingLog
	}

	log := factory(name)
	loggerPostCreation(log)
	logs[name] = log
	return log
}

func loggerPostCreation(logger Logger) {
	if logConfig.Levels == nil {
		return
	}

	if existingLevel, ok := logConfig.Levels[logger.Name()]; ok {
		if existingLevel != logger.Level() {
			logger.SetLevel(existingLevel)
		}
	}
}

func ConfigureLogging(config LogConfig) {
	logConfig = config

	if logConfig.Type == LogTypeMultiple {
		configureMultipleLoggers()
	} else if logConfig.Type == LogTypeSingleton {
		configureSingletonLogger()
	} else {
		panic("unknown LogType in LogConfig")
	}
}

func configureSingletonLogger() {
	if logConfig.SingletonLogConfig.Logger == nil {
		panic("SingletonLogConfig.Logger must be set for LogTypeSingleton")
	}
	log = logConfig.SingletonLogConfig.Logger
}

func configureMultipleLoggers() {
	if logConfig.MultipleLogConfig.Factory == nil {
		panic("MultipleLogConfig.Factory must be set for LogTypeMultiple")
	}
	factory = logConfig.MultipleLogConfig.Factory
}

func SetLevel(name string, level Level) error {
	mutex.Lock()
	defer mutex.Unlock()

	logConfig.Levels[name] = level

	if logger, ok := logs[name]; ok {
		return logger.SetLevel(level)
	}

	return nil
}
