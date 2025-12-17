package metric

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type OtelMeter struct {
	meter metric.Meter
}

var _ Meter = (*OtelMeter)(nil)

func NewOtelMeter(meter metric.Meter) *OtelMeter {
	return &OtelMeter{meter: meter}
}

func (m *OtelMeter) AddGauge(ctx context.Context, name string, description string, unit string, value float64, attrs ...Attribute) error {
	gauge, err := m.meter.Float64Gauge(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge %s: %w", name, err)
	}
	gauge.Record(ctx, value, metric.WithAttributes(toOtelAttributes(attrs)...))
	return nil
}

func (m *OtelMeter) AddCounter(ctx context.Context, name string, description string, unit string, value float64, attrs ...Attribute) error {
	upDownCounter, err := m.meter.Float64UpDownCounter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		return fmt.Errorf("failed to create up-down counter %s: %w", name, err)
	}
	upDownCounter.Add(ctx, value, metric.WithAttributes(toOtelAttributes(attrs)...))
	return nil
}

func toOtelAttributes(attrs []Attribute) []attribute.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	otelAttrs := make([]attribute.KeyValue, 0, len(attrs))
	for _, attr := range attrs {
		otelAttrs = append(otelAttrs, toOtelAttribute(attr))
	}
	return otelAttrs
}

func toOtelAttribute(attr Attribute) attribute.KeyValue {
	switch v := attr.Value.(type) {
	case string:
		return attribute.String(attr.Key, v)
	case int:
		return attribute.Int(attr.Key, v)
	case int64:
		return attribute.Int64(attr.Key, v)
	case float64:
		return attribute.Float64(attr.Key, v)
	case bool:
		return attribute.Bool(attr.Key, v)
	case []string:
		return attribute.StringSlice(attr.Key, v)
	case []int:
		return attribute.IntSlice(attr.Key, v)
	case []int64:
		return attribute.Int64Slice(attr.Key, v)
	case []float64:
		return attribute.Float64Slice(attr.Key, v)
	case []bool:
		return attribute.BoolSlice(attr.Key, v)
	default:

		return attribute.String(attr.Key, fmt.Sprintf("%v", v))
	}
}
