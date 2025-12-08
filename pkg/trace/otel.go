package trace

import (
	"context"
	"fmt"
	"time"

	"github.com/bruno303/go-toolkit/pkg/utils/array"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	tracelib "go.opentelemetry.io/otel/trace"
)

type OtelTracerAdapter struct{}

func NewOtelTracerAdapter() OtelTracerAdapter {
	return OtelTracerAdapter{}
}

func (t OtelTracerAdapter) Trace(ctx context.Context, cfg *TraceConfig, cb TraceCallback) (any, error) {
	if cfg == nil {
		cfg = DefaultTraceCfg()
	}
	cfg.Validate()

	ctx, span := startSpan(ctx, cfg.TraceName, cfg.SpanName)
	defer span.End()
	res, err := cb(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return res, err
}

func (t OtelTracerAdapter) ExtractTraceIds(ctx context.Context) TraceIDs {
	span := tracelib.SpanFromContext(ctx)
	return TraceIDs{
		TraceID: span.SpanContext().TraceID().String(),
		SpanID:  span.SpanContext().SpanID().String(),
		IsValid: span.SpanContext().IsValid(),
	}
}

func (t OtelTracerAdapter) InjectAttributes(ctx context.Context, attrs ...Attribute) {
	span := tracelib.SpanFromContext(ctx)
	if span == nil {
		return
	}
	otelAttrs := array.Map(attrs, func(a Attribute) attribute.KeyValue {
		return attribute.String(a.Key, a.Value)
	})
	span.SetAttributes(otelAttrs...)
}

func (t OtelTracerAdapter) InjectError(ctx context.Context, err error) {
	span := tracelib.SpanFromContext(ctx)
	if span == nil {
		return
	}
	span.RecordError(err)
}

func startSpan(ctx context.Context, tracerName string, spanName string) (context.Context, tracelib.Span) {
	return otel.Tracer(tracerName).Start(
		ctx,
		fmt.Sprintf("%s.%s", tracerName, spanName),
		tracelib.WithTimestamp(time.Now()),
	)
}

func EndTrace(ctx context.Context) {
	span := tracelib.SpanFromContext(ctx)
	if span == nil {
		return
	}
	span.End()
}
