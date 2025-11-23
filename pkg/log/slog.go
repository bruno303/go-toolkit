package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/bruno303/go-toolkit/pkg/trace"
)

type (
	ExtractFunc func(context.Context) []any
	SlogAdapter struct {
		logger                *slog.Logger
		slogLevel             *slog.LevelVar
		extractAdditionalInfo func(context.Context) []any
		name                  string
		level                 Level
	}
	SlogAdapterOpts struct {
		Level                 Level
		FormatJson            bool
		Name                  string
		ExtractAdditionalInfo func(context.Context) []any
		AddSource             bool
		Environment           string
	}
)

var _ Logger = (*SlogAdapter)(nil)

func NewSlogAdapter(opts SlogAdapterOpts) SlogAdapter {
	level := toSlogLevel(opts.Level)
	levelVar := &slog.LevelVar{}
	levelVar.Set(level)

	if opts.ExtractAdditionalInfo == nil {
		opts.ExtractAdditionalInfo = func(context.Context) []any { return nil }
	}
	extractInfo := func(ctx context.Context) []any {
		ai := make([]any, 0)
		ai = append(ai, "source", opts.Name)
		ai = append(ai, opts.ExtractAdditionalInfo(ctx)...)
		ai = append(ai, extractTraceInfo(ctx)...)
		ai = append(ai, "env", opts.Environment)
		return ai
	}

	handlerOpts := &slog.HandlerOptions{
		AddSource: opts.AddSource,
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
		slogLevel:             levelVar,
		extractAdditionalInfo: extractInfo,
		name:                  opts.Name,
		level:                 opts.Level,
	}
}

func extractTraceInfo(ctx context.Context) []any {
	traceIDs := trace.ExtractTraceIds(ctx)
	if traceIDs.IsValid {
		return []any{"trace_id", traceIDs.TraceID, "span_id", traceIDs.SpanID}
	}
	return nil
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
	la.slogLevel.Set(toSlogLevel(l))
	return nil
}

func (l SlogAdapter) Shutdown(context.Context) error {
	return nil
}

func (l SlogAdapter) Name() string {
	return l.name
}

func (l SlogAdapter) Level() Level {
	return l.level
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
