package trace

import (
	"context"
)

type NoOpTracer struct{}

func NewNoOpTracer() NoOpTracer {
	return NoOpTracer{}
}

func (t NoOpTracer) Trace(ctx context.Context, cfg *TraceConfig, cb TraceCallback) (any, error) {
	return cb(ctx)
}

func (t NoOpTracer) ExtractTraceIds(ctx context.Context) TraceIDs {
	return TraceIDs{
		TraceID: "",
		SpanID:  "",
		IsValid: false,
	}
}

func (t NoOpTracer) InjectAttributes(ctx context.Context, attrs ...Attribute) {}

func (t NoOpTracer) InjectError(ctx context.Context, err error) {}
