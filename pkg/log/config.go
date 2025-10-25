package log

type (
	Level         uint
	LogType       uint
	LoggerFactory func(string) Logger

	LogConfig struct {
		Levels             map[string]Level
		Type               LogType
		MultipleLogConfig  MultipleLogConfig
		SingletonLogConfig SingletonLogConfig
	}

	SingletonLogConfig struct {
		Logger Logger
	}

	MultipleLogConfig struct {
		Factory LoggerFactory
	}
)

const (
	LogTypeSingleton LogType = iota
	LogTypeMultiple
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)
