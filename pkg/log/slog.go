package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type (
	ExtractFunc func(context.Context) []any
	SlogAdapter struct {
		logger                *slog.Logger
		level                 *slog.LevelVar
		extractAdditionalInfo func(context.Context) []any
	}
	SlogAdapterOpts struct {
		Level                 Level
		FormatJson            bool
		Source                string
		ExtractAdditionalInfo func(context.Context) []any
	}
)

func NewSlogAdapter(opts SlogAdapterOpts) SlogAdapter {
	level := toSlogLevel(opts.Level)
	levelVar := &slog.LevelVar{}
	levelVar.Set(level)
	var extractInfo func(context.Context) []any = nil

	if opts.ExtractAdditionalInfo == nil {
		opts.ExtractAdditionalInfo = func(context.Context) []any { return nil }
	}
	extractInfo = func(ctx context.Context) []any {
		ai := make([]any, 0)
		ai = append(ai, "source", opts.Source)
		ai = append(ai, opts.ExtractAdditionalInfo(ctx)...)
		return ai
	}

	handlerOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     levelVar,
	}

	var handler slog.Handler
	if opts.FormatJson {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}
	return SlogAdapter{
		logger:                slog.New(handler),
		level:                 levelVar,
		extractAdditionalInfo: extractInfo,
	}
}

func (l SlogAdapter) Info(ctx context.Context, msg string, args ...any) {
	if !l.logger.Enabled(ctx, slog.LevelInfo) {
		return
	}

	l.logger.InfoContext(ctx, fmt.Sprintf(msg, args...), l.extractAdditionalInfo(ctx)...)
}

func (l SlogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	if !l.logger.Enabled(ctx, slog.LevelDebug) {
		return
	}
	l.logger.DebugContext(ctx, fmt.Sprintf(msg, args...), l.extractAdditionalInfo(ctx)...)
}

func (l SlogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	if !l.logger.Enabled(ctx, slog.LevelWarn) {
		return
	}
	l.logger.WarnContext(ctx, fmt.Sprintf(msg, args...), l.extractAdditionalInfo(ctx)...)
}

func (l SlogAdapter) Error(ctx context.Context, msg string, err error) {
	additionalData := l.extractAdditionalInfo(ctx)
	additionalData = append(additionalData, "err", err.Error())
	l.logger.ErrorContext(ctx, msg, additionalData...)
}

func (la SlogAdapter) SetLevel(l Level) error {
	la.level.Set(toSlogLevel(l))
	return nil
}

func toSlogLevel(l Level) slog.Level {
	switch l {
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelDebug:
		return slog.LevelDebug
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
