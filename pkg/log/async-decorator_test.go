package log

import (
	"context"
	"errors"
	"testing"
)

func TestAsyncDecorator_Info(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}
	ad := NewAsyncDecorator(logger)

	ad.Info(context.Background(), "test message")
	ad.Shutdown(ctx)

	if len(logger.messages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(logger.messages))
	}
	if logger.messages[0].level != "Info" {
		t.Errorf("expected log level Info, got %s", logger.messages[0].level)
	}
	if logger.messages[0].msg != "test message" {
		t.Errorf("expected log message 'test message', got %s", logger.messages[0].msg)
	}
}

func TestAsyncDecorator_Debug(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}
	ad := NewAsyncDecorator(logger)

	ad.Debug(context.Background(), "test message")
	ad.Shutdown(ctx)

	if len(logger.messages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(logger.messages))
	}
	if logger.messages[0].level != "Debug" {
		t.Errorf("expected log level Debug, got %s", logger.messages[0].level)
	}
	if logger.messages[0].msg != "test message" {
		t.Errorf("expected log message 'test message', got %s", logger.messages[0].msg)
	}
}

func TestAsyncDecorator_Warn(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}
	ad := NewAsyncDecorator(logger)

	ad.Warn(context.Background(), "test message")
	ad.Shutdown(ctx)

	if len(logger.messages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(logger.messages))
	}
	if logger.messages[0].level != "Warn" {
		t.Errorf("expected log level Warn, got %s", logger.messages[0].level)
	}
	if logger.messages[0].msg != "test message" {
		t.Errorf("expected log message 'test message', got %s", logger.messages[0].msg)
	}
}

func TestAsyncDecorator_Error(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}
	ad := NewAsyncDecorator(logger)

	ad.Error(context.Background(), "test message", errors.New("test error"))
	ad.Shutdown(ctx)

	if len(logger.messages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(logger.messages))
	}
	if logger.messages[0].level != "Error" {
		t.Errorf("expected log level Error, got %s", logger.messages[0].level)
	}
	if logger.messages[0].msg != "test message" {
		t.Errorf("expected log message 'test message', got %s", logger.messages[0].msg)
	}
	if logger.messages[0].err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestAsyncDecorator_Shutdown(t *testing.T) {
	ctx := context.Background()
	logger := &mockLogger{}
	ad := NewAsyncDecorator(logger)

	ad.Info(context.Background(), "test message")
	ad.Shutdown(ctx)

	if len(logger.messages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(logger.messages))
	}
}

type mockLogger struct {
	messages []logMessage
}

func (m *mockLogger) SetLevel(l Level) error {
	return nil
}

func (m *mockLogger) Shutdown(context.Context) error {
	return nil
}

func (m *mockLogger) Info(ctx context.Context, msg string, args ...any) {
	m.messages = append(m.messages, logMessage{
		level: "Info",
		msg:   msg,
		args:  args,
	})
}

func (m *mockLogger) Debug(ctx context.Context, msg string, args ...any) {
	m.messages = append(m.messages, logMessage{
		level: "Debug",
		msg:   msg,
		args:  args,
	})
}

func (m *mockLogger) Warn(ctx context.Context, msg string, args ...any) {
	m.messages = append(m.messages, logMessage{
		level: "Warn",
		msg:   msg,
		args:  args,
	})
}

func (m *mockLogger) Error(ctx context.Context, msg string, err error) {
	m.messages = append(m.messages, logMessage{
		level: "Error",
		msg:   msg,
		err:   err,
	})
}

type logMessage struct {
	level string
	msg   string
	args  []any
	err   error
}
