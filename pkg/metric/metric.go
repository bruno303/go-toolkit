package metric

import (
	"context"
	"sync"
)

type Attribute struct {
	Key   string
	Value any
}

type Meter interface {
	AddGauge(ctx context.Context, name string, description string, unit string, value float64, attrs ...Attribute) error
	AddCounter(ctx context.Context, name string, description string, unit string, value float64, attrs ...Attribute) error
}

type MeterProvider interface {
	Meter(name string) Meter
	Shutdown(ctx context.Context) error
}

var (
	globalMeter Meter = NewNoOpMeter()
	once        sync.Once
)

func GetMeter() Meter {
	return globalMeter
}

func SetMeter(m Meter) {
	once.Do(func() {
		globalMeter = m
	})
}

func NewAttribute(key string, value any) Attribute {
	return Attribute{Key: key, Value: value}
}
