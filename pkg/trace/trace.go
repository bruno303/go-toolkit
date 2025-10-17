package trace

import (
	"context"
	"errors"
	"sync"
)

type (
	EndFunc     func()
	TraceKind   int
	TraceConfig struct {
		Kind      TraceKind
		TraceName string
		SpanName  string
	}
	TraceIDs struct {
		TraceID string
		SpanID  string
		IsValid bool
	}
	Tracer interface {
		Trace(ctx context.Context, cfg *TraceConfig, cb TraceCallback) (any, error)
		ExtractTraceIds(ctx context.Context) TraceIDs
		InjectAttributes(ctx context.Context, attrs ...Attribute)
		InjectError(ctx context.Context, err error)
	}
	TraceCallback func(ctx context.Context) (any, error)
)

var (
	tracer Tracer = NewNoOpTracer()
	once   sync.Once
)

const (
	_ TraceKind = iota
	TraceKindServer
	TraceKindConsumer
	TraceKindProducer
)

func GetTracer() Tracer {
	return tracer
}

func SetTracer(t Tracer) {
	once.Do(func() {
		tracer = t
	})
}

func Trace(ctx context.Context, cfg *TraceConfig, cb TraceCallback) (any, error) {
	return tracer.Trace(ctx, cfg, cb)
}

func ExtractTraceIds(ctx context.Context) TraceIDs {
	return tracer.ExtractTraceIds(ctx)
}

func InjectAttributes(ctx context.Context, attrs ...Attribute) {
	tracer.InjectAttributes(ctx, attrs...)
}

func InjectError(ctx context.Context, err error) {
	tracer.InjectError(ctx, err)
}

func NameConfig(traceName string, spanName string) *TraceConfig {
	return &TraceConfig{TraceName: traceName, SpanName: spanName}
}

func DefaultTraceCfg() *TraceConfig {
	return &TraceConfig{
		Kind: TraceKindServer,
	}
}

func (c *TraceConfig) Validate() error {
	if c.Kind == 0 {
		c.Kind = TraceKindServer
	}
	if c.TraceName == "" {
		return errors.New("TraceName must be informed")
	}
	if c.SpanName == "" {
		return errors.New("SpanName must be informed")
	}
	return nil
}
