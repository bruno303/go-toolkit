package metric

import (
	"context"
)

type NoOpMeter struct{}

var _ Meter = (*NoOpMeter)(nil)

func NewNoOpMeter() *NoOpMeter {
	return &NoOpMeter{}
}

func (m *NoOpMeter) AddCounter(ctx context.Context, name string, description string, unit string, value float64, attrs ...Attribute) error {
	return nil
}

func (m *NoOpMeter) AddGauge(ctx context.Context, name string, description string, unit string, value float64, attrs ...Attribute) error {
	return nil
}
