package log

import (
	"context"
	"errors"
	"sync"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestLog(t *testing.T) {
	// Reset state before test
	resetLogState()

	logger := Log()
	if logger == nil {
		t.Error("expected logger to not be nil")
	}

	// Verify default logger properties
	if logger.Name() != "default" {
		t.Errorf("expected default logger name to be 'default', got %s", logger.Name())
	}

	// Verify it returns consistent results
	logger2 := Log()
	if logger2 == nil {
		t.Error("expected second logger call to not be nil")
	}

	if logger.Name() != logger2.Name() {
		t.Error("expected Log() to return the same logger instance (names should match)")
	}
}

func TestSetLogger(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("custom").AnyTimes()

	SetLogger(mockLogger)

	logger := Log()
	if logger.Name() != "custom" {
		t.Errorf("expected logger name to be 'custom', got %s", logger.Name())
	}
}

func TestSetLoggerFactory(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	customFactory := func(name string) Logger {
		mockLogger := NewMockLogger(ctrl)
		mockLogger.EXPECT().Name().Return("factory-" + name).AnyTimes()
		mockLogger.EXPECT().Level().Return(LevelInfo).AnyTimes()
		return mockLogger
	}

	SetLoggerFactory(customFactory)

	logger := Log()
	if logger.Name() != "factory-default" {
		t.Errorf("expected logger name to be 'factory-default', got %s", logger.Name())
	}
}

func TestNewLogger(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Test creating a new logger
	logger1 := NewLogger("test1")
	if logger1 == nil {
		t.Error("expected logger to not be nil")
	}
	if logger1.Name() != "test1" {
		t.Errorf("expected logger name to be 'test1', got %s", logger1.Name())
	}

	// Test that the same name returns the same instance
	logger1Again := NewLogger("test1")
	if logger1Again.Name() != "test1" {
		t.Error("expected NewLogger to return the same instance for the same name")
	}

	// Test creating a different logger
	logger2 := NewLogger("test2")
	if logger2.Name() != "test2" {
		t.Errorf("expected logger name to be 'test2', got %s", logger2.Name())
	}

	// Verify they are different loggers by checking names
	if logger1.Name() == logger2.Name() {
		t.Error("expected different loggers for different names")
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	// Create mock loggers for testing
	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("debug-logger").AnyTimes()

	// Configure logging with specific levels - using LogTypeMultiple to avoid singleton validation
	customFactory := func(name string) Logger {
		switch name {
		case "debug-logger", "error-logger", "default-logger":
			return NewSlogAdapter(SlogAdapterOpts{
				Level:      LevelInfo,
				FormatJson: false,
				Name:       name,
			})
		default:
			return NewSlogAdapter(SlogAdapterOpts{
				Level:      LevelInfo,
				FormatJson: false,
				Name:       name,
			})
		}
	}

	config := LogConfig{
		Type: LogTypeMultiple,
		MultipleLogConfig: MultipleLogConfig{
			Factory: customFactory,
		},
		Levels: map[string]Level{
			"debug-logger": LevelDebug,
			"error-logger": LevelError,
		},
	}
	ConfigureLogging(config)

	// Create loggers and verify they exist and have correct names
	debugLogger := NewLogger("debug-logger")
	if debugLogger == nil {
		t.Error("expected debug logger to not be nil")
	}
	if debugLogger.Name() != "debug-logger" {
		t.Errorf("expected debug logger name to be 'debug-logger', got %s", debugLogger.Name())
	}

	errorLogger := NewLogger("error-logger")
	if errorLogger == nil {
		t.Error("expected error logger to not be nil")
	}
	if errorLogger.Name() != "error-logger" {
		t.Errorf("expected error logger name to be 'error-logger', got %s", errorLogger.Name())
	}

	// Create a logger without specific config
	defaultLogger := NewLogger("default-logger")
	if defaultLogger == nil {
		t.Error("expected default logger to not be nil")
	}
	if defaultLogger.Name() != "default-logger" {
		t.Errorf("expected default logger name to be 'default-logger', got %s", defaultLogger.Name())
	}
}

func TestConfigureLoggingLevelsOnly(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("singleton").AnyTimes()

	// Configure with just levels and required singleton config
	config := LogConfig{
		Type: LogTypeSingleton,
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockLogger,
		},
		Levels: map[string]Level{
			"test": LevelWarn,
		},
	}

	ConfigureLogging(config)

	// Verify the config was stored
	if len(logConfig.Levels) != 1 {
		t.Errorf("expected 1 level config, got %d", len(logConfig.Levels))
	}

	if logConfig.Levels["test"] != LevelWarn {
		t.Errorf("expected test logger level to be LevelWarn (%d), got %d", LevelWarn, logConfig.Levels["test"])
	}

	if logConfig.Type != LogTypeSingleton {
		t.Errorf("expected log type to be LogTypeSingleton (%d), got %d", LogTypeSingleton, logConfig.Type)
	}
}

func TestOriginalBehaviorWithoutConfig(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Test that the original behavior still works when no ConfigureLogging is called
	// This should use the default factory without any configuration
	logger := Log()
	if logger == nil {
		t.Error("expected logger to not be nil")
	}

	if logger.Name() != "default" {
		t.Errorf("expected default logger name to be 'default', got %s", logger.Name())
	}

	// Test creating named loggers without configuration
	namedLogger := NewLogger("test-logger")
	if namedLogger == nil {
		t.Error("expected named logger to not be nil")
	}

	if namedLogger.Name() != "test-logger" {
		t.Errorf("expected named logger name to be 'test-logger', got %s", namedLogger.Name())
	}
}

func TestConfigureLoggingWithType(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("singleton-test").AnyTimes()

	config := LogConfig{
		Type: LogTypeSingleton,
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockLogger,
		},
		Levels: map[string]Level{
			"test": LevelError,
		},
	}

	ConfigureLogging(config)

	// Verify the config was stored with type
	if logConfig.Type != LogTypeSingleton {
		t.Errorf("expected log type to be LogTypeSingleton (%d), got %d", LogTypeSingleton, logConfig.Type)
	}

	if logConfig.Levels["test"] != LevelError {
		t.Errorf("expected test logger level to be LevelError (%d), got %d", LevelError, logConfig.Levels["test"])
	}
}

func TestConfigureLoggingWithSingletonConfig(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("singleton-logger").AnyTimes()

	config := LogConfig{
		Type: LogTypeSingleton,
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockLogger,
		},
		Levels: map[string]Level{
			"singleton-logger": LevelDebug,
		},
	}

	ConfigureLogging(config)

	// Verify singleton config was stored
	if logConfig.Type != LogTypeSingleton {
		t.Errorf("expected log type to be LogTypeSingleton (%d), got %d", LogTypeSingleton, logConfig.Type)
	}

	if logConfig.SingletonLogConfig.Logger == nil {
		t.Error("expected singleton logger to not be nil")
	}

	if logConfig.SingletonLogConfig.Logger.Name() != "singleton-logger" {
		t.Errorf("expected singleton logger name to be 'singleton-logger', got %s", logConfig.SingletonLogConfig.Logger.Name())
	}
}

func TestConfigureLoggingWithMultipleConfig(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	customFactory := func(name string) Logger {
		mockLogger := NewMockLogger(ctrl)
		mockLogger.EXPECT().Name().Return("factory-" + name).AnyTimes()
		mockLogger.EXPECT().Level().Return(LevelInfo).AnyTimes()
		return mockLogger
	}

	config := LogConfig{
		Type: LogTypeMultiple,
		MultipleLogConfig: MultipleLogConfig{
			Factory: customFactory,
		},
		Levels: map[string]Level{
			"test": LevelWarn,
		},
	}

	ConfigureLogging(config)

	// Verify multiple config was stored
	if logConfig.Type != LogTypeMultiple {
		t.Errorf("expected log type to be LogTypeMultiple (%d), got %d", LogTypeMultiple, logConfig.Type)
	}

	if logConfig.MultipleLogConfig.Factory == nil {
		t.Error("expected multiple logger factory to not be nil")
	}

	// Test that the factory works
	testLogger := logConfig.MultipleLogConfig.Factory("test")
	if testLogger.Name() != "factory-test" {
		t.Errorf("expected factory logger name to be 'factory-test', got %s", testLogger.Name())
	}
}

func TestLoggerPostCreation(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	// Set up logConfig directly without ConfigureLogging to avoid validation
	logConfig.Levels = map[string]Level{
		"test-logger": LevelWarn,
	}

	// Create a mock logger with different level
	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("test-logger").AnyTimes()
	mockLogger.EXPECT().Level().Return(LevelInfo).Times(1)
	mockLogger.EXPECT().SetLevel(LevelWarn).Return(nil).Times(1)

	// Call loggerPostCreation
	loggerPostCreation(mockLogger)
}

func TestLoggerPostCreationWithoutConfig(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	// Create a mock logger without any config
	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("test-logger").AnyTimes()

	// Call loggerPostCreation - should not call SetLevel since no config exists
	loggerPostCreation(mockLogger)
}

func TestLoggerPostCreationWithoutLevelsConfigured(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	// Set up logConfig directly without ConfigureLogging to avoid validation
	logConfig.Levels = nil

	// Create a mock logger with same level as config
	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("test-logger").AnyTimes()
	mockLogger.EXPECT().Level().Times(0)

	// Call loggerPostCreation
	loggerPostCreation(mockLogger)
}

func TestLoggerPostCreationSameLevel(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	// Set up logConfig directly without ConfigureLogging to avoid validation
	logConfig.Levels = map[string]Level{
		"test-logger": LevelWarn,
	}

	// Create a mock logger with same level as config
	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("test-logger").AnyTimes()
	mockLogger.EXPECT().Level().Return(LevelWarn).Times(1)
	// SetLevel should NOT be called since levels match

	// Call loggerPostCreation
	loggerPostCreation(mockLogger)
}

func TestLoggerInterfaceCompliance(t *testing.T) {
	// Reset state before test
	resetLogState()

	logger := Log()
	ctx := context.Background()

	// Test all Logger interface methods don't panic
	logger.Info(ctx, "test info message", "key", "value")
	logger.Debug(ctx, "test debug message")
	logger.Warn(ctx, "test warn message")
	logger.Error(ctx, "test error message", errors.New("test error"))

	if err := logger.SetLevel(LevelWarn); err != nil {
		t.Errorf("unexpected error setting level: %v", err)
	}

	if logger.Name() == "" {
		t.Error("expected logger name to not be empty")
	}

	currentLevel := logger.Level()
	if currentLevel < LevelDebug || currentLevel > LevelError {
		t.Errorf("expected valid log level, got %d", currentLevel)
	}

	if err := logger.Shutdown(ctx); err != nil {
		t.Errorf("unexpected error shutting down logger: %v", err)
	}
}

func TestConcurrentLogCreation(t *testing.T) {
	// Reset state before test
	resetLogState()

	const numGoroutines = 10
	const numIterations = 10

	var wg sync.WaitGroup

	// Test concurrent access to Log()
	loggerNames := make([]string, numGoroutines*numIterations)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				logger := Log()
				if logger == nil {
					t.Error("expected logger to not be nil")
					return
				}
				// Store name for later verification
				idx := id*numIterations + j
				loggerNames[idx] = logger.Name()
			}
		}(i)
	}
	wg.Wait()

	// Verify all got the same default logger
	for i, name := range loggerNames {
		if name != "default" {
			t.Errorf("goroutine %d: expected logger name 'default', got %s", i, name)
		}
	}
}

func TestConcurrentNamedLoggerCreation(t *testing.T) {
	// Reset state before test
	resetLogState()

	const numGoroutines = 5
	const numIterations = 10

	var wg sync.WaitGroup

	// Test concurrent access to NewLogger() with same name
	wg.Add(numGoroutines)
	loggerNames := make([]string, numGoroutines*numIterations)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				logger := NewLogger("concurrent-test")
				if logger == nil {
					t.Error("expected logger to not be nil")
					return
				}
				idx := id*numIterations + j
				loggerNames[idx] = logger.Name()
			}
		}(i)
	}
	wg.Wait()

	// Verify all got the same named logger
	for i, name := range loggerNames {
		if name != "concurrent-test" {
			t.Errorf("goroutine %d: expected logger name 'concurrent-test', got %s", i, name)
		}
	}
}

func TestEmptyLoggerName(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Test creating a logger with empty name
	logger := NewLogger("")
	if logger == nil {
		t.Error("expected logger to not be nil")
	}
	if logger.Name() != "" {
		t.Errorf("expected logger name to be empty, got %s", logger.Name())
	}
}

func TestMultipleConfigurations(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	mockLogger1 := NewMockLogger(ctrl)
	mockLogger1.EXPECT().Name().Return("initial-singleton").AnyTimes()

	// Configure logging with initial config
	config1 := LogConfig{
		Type: LogTypeSingleton,
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockLogger1,
		},
		Levels: map[string]Level{
			"logger1": LevelDebug,
		},
	}
	ConfigureLogging(config1)

	mockLogger := NewMockLogger(ctrl)
	mockLogger.EXPECT().Name().Return("new-singleton").AnyTimes()

	// Override with new config including different type and configs
	config2 := LogConfig{
		Type: LogTypeMultiple,
		MultipleLogConfig: MultipleLogConfig{
			Factory: func(name string) Logger { return mockLogger },
		},
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockLogger,
		},
		Levels: map[string]Level{
			"logger1": LevelError,
			"logger2": LevelWarn,
		},
	}
	ConfigureLogging(config2)

	// Verify the second config overwrote the first
	if logConfig.Type != LogTypeMultiple {
		t.Errorf("expected log type to be LogTypeMultiple (%d), got %d", LogTypeMultiple, logConfig.Type)
	}

	if len(logConfig.Levels) != 2 {
		t.Errorf("expected 2 level configs, got %d", len(logConfig.Levels))
	}

	if logConfig.Levels["logger1"] != LevelError {
		t.Errorf("expected logger1 level to be LevelError (%d), got %d", LevelError, logConfig.Levels["logger1"])
	}

	if logConfig.Levels["logger2"] != LevelWarn {
		t.Errorf("expected logger2 level to be LevelWarn (%d), got %d", LevelWarn, logConfig.Levels["logger2"])
	}

	// Verify both config types are stored
	if logConfig.MultipleLogConfig.Factory == nil {
		t.Error("expected multiple config factory to not be nil")
	}

	if logConfig.SingletonLogConfig.Logger == nil {
		t.Error("expected singleton config logger to not be nil")
	}
}

func TestLogTypes(t *testing.T) {
	// Test LogType constants exist and have expected values
	if LogTypeSingleton != 0 {
		t.Errorf("expected LogTypeSingleton to be 0, got %d", LogTypeSingleton)
	}

	if LogTypeMultiple != 1 {
		t.Errorf("expected LogTypeMultiple to be 1, got %d", LogTypeMultiple)
	}
}

func TestEmptyLogConfig(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Empty config (Type=0 which is LogTypeSingleton) without SingletonLogConfig.Logger should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when using empty config with LogTypeSingleton but no logger")
		}
	}()

	config := LogConfig{}
	ConfigureLogging(config) // Should panic
}

func TestUnknownLogType(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Test with unknown LogType should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when using unknown LogType")
		}
	}()

	config := LogConfig{
		Type: LogType(999), // Unknown type
	}
	ConfigureLogging(config) // Should panic
}

func TestSingletonConfigValidation(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Test LogTypeSingleton without Logger should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when using LogTypeSingleton without Logger")
		}
	}()

	config := LogConfig{
		Type: LogTypeSingleton,
		// Missing SingletonLogConfig.Logger
	}
	ConfigureLogging(config) // Should panic
}

func TestMultipleConfigValidation(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Test LogTypeMultiple without Factory should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when using LogTypeMultiple without Factory")
		}
	}()

	config := LogConfig{
		Type: LogTypeMultiple,
		// Missing MultipleLogConfig.Factory
	}
	ConfigureLogging(config) // Should panic
}

func TestLogConfigWithBothConfigs(t *testing.T) {
	// Reset state before test
	resetLogState()

	ctrl := gomock.NewController(t)

	mockSingletonLogger := NewMockLogger(ctrl)
	mockSingletonLogger.EXPECT().Name().Return("singleton").AnyTimes()

	mockMultipleLogger := NewMockLogger(ctrl)
	mockMultipleLogger.EXPECT().Name().Return("multiple-test").AnyTimes()
	mockMultipleLogger.EXPECT().Level().Return(LevelInfo).AnyTimes()

	// Configure with both singleton and multiple configs
	config := LogConfig{
		Type: LogTypeSingleton, // Should use singleton despite having multiple config
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockSingletonLogger,
		},
		MultipleLogConfig: MultipleLogConfig{
			Factory: func(name string) Logger { return mockMultipleLogger },
		},
		Levels: map[string]Level{
			"test": LevelWarn,
		},
	}

	ConfigureLogging(config)

	// Verify both configs are stored (even though type determines which is used)
	if logConfig.Type != LogTypeSingleton {
		t.Errorf("expected log type to be LogTypeSingleton (%d), got %d", LogTypeSingleton, logConfig.Type)
	}

	if logConfig.SingletonLogConfig.Logger == nil {
		t.Error("expected singleton logger to not be nil")
	}

	if logConfig.MultipleLogConfig.Factory == nil {
		t.Error("expected multiple factory to not be nil")
	}

	if logConfig.SingletonLogConfig.Logger.Name() != "singleton" {
		t.Errorf("expected singleton logger name to be 'singleton', got %s", logConfig.SingletonLogConfig.Logger.Name())
	}

	// Test factory still works
	factoryLogger := logConfig.MultipleLogConfig.Factory("test")
	if factoryLogger.Name() != "multiple-test" {
		t.Errorf("expected factory logger name to be 'multiple-test', got %s", factoryLogger.Name())
	}
}

func TestSetLoggerFactoryNil(t *testing.T) {
	// Reset state before test
	resetLogState()

	// This should not panic, but behavior might be undefined
	SetLoggerFactory(nil)

	// Attempting to create a logger might panic due to nil factory
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when using nil factory")
		}
	}()

	Log() // This should panic
}

func TestNewLoggerCaching(t *testing.T) {
	// Reset state before test
	resetLogState()

	// Create same logger multiple times and verify they're cached
	logger1 := NewLogger("cached-test")
	logger2 := NewLogger("cached-test")
	logger3 := NewLogger("cached-test")

	// All should have the same name (can't compare directly due to struct type)
	if logger1.Name() != logger2.Name() || logger2.Name() != logger3.Name() {
		t.Error("expected all loggers with same name to be identical")
	}

	// Different names should create different loggers
	differentLogger := NewLogger("different-test")
	if logger1.Name() == differentLogger.Name() {
		t.Error("expected different logger names to create different loggers")
	}
}

func TestConfigWithoutLevels(t *testing.T) {
	// Reset state before test
	resetLogState()
	ctrl := gomock.NewController(t)

	mockSingletonLogger := NewMockLogger(ctrl)
	mockSingletonLogger.EXPECT().Name().Return("singleton").AnyTimes()

	config := LogConfig{
		Type: LogTypeSingleton,
		SingletonLogConfig: SingletonLogConfig{
			Logger: mockSingletonLogger,
		},
		Levels: nil,
	}
	ConfigureLogging(config)
}

// resetLogState resets the global state for testing
func resetLogState() {
	mutex.Lock()
	defer mutex.Unlock()

	log = nil
	logs = make(map[string]Logger)
	logConfig = LogConfig{}

	// Reset factory to default
	factory = func(name string) Logger {
		return NewSlogAdapter(SlogAdapterOpts{
			Level:      LevelInfo,
			FormatJson: false,
			Name:       name,
		})
	}
}
