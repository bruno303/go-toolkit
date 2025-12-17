package metric

import (
	"context"
	"testing"
)

func TestNoOpMeter(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	err := meter.AddCounter(ctx, "test.counter", "A test counter", "1", 1.0, NewAttribute("key", "value"))
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}
	err = meter.AddCounter(ctx, "test.counter", "A test counter", "1", 5.0)
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}

	err = meter.AddGauge(ctx, "test.gauge", "A test gauge", "celsius", 23.5, NewAttribute("location", "server1"))
	if err != nil {
		t.Fatalf("failed to add gauge: %v", err)
	}
}

func TestGlobalMeter(t *testing.T) {
	ctx := context.Background()

	meter := GetMeter()
	if meter == nil {
		t.Fatal("expected meter to not be nil")
	}

	err := meter.AddCounter(ctx, "test.counter", "test", "1", 1.0)
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}
}

func TestSetMeter(t *testing.T) {
	customMeter := NewNoOpMeter()
	SetMeter(customMeter)

	retrievedMeter := GetMeter()
	if retrievedMeter == nil {
		t.Fatal("expected meter to not be nil")
	}

	ctx := context.Background()
	err := retrievedMeter.AddCounter(ctx, "test.counter", "test", "1", 1.0)
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}
}

func TestNewAttribute(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value any
	}{
		{"string", "key1", "value1"},
		{"int", "key2", 42},
		{"float", "key3", 3.14},
		{"bool", "key4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := NewAttribute(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("expected key %s, got %s", tt.key, attr.Key)
			}
			if attr.Value != tt.value {
				t.Errorf("expected value %v, got %v", tt.value, attr.Value)
			}
		})
	}
}

func TestAttributeTypes(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	err := meter.AddCounter(ctx, "test.counter", "test", "1", 1.0,
		NewAttribute("string", "value"),
		NewAttribute("int", 123),
		NewAttribute("int64", int64(456)),
		NewAttribute("float64", 78.9),
		NewAttribute("bool", true),
		NewAttribute("strings", []string{"a", "b", "c"}),
		NewAttribute("ints", []int{1, 2, 3}),
		NewAttribute("int64s", []int64{4, 5, 6}),
		NewAttribute("float64s", []float64{7.8, 9.0}),
		NewAttribute("bools", []bool{true, false}),
	)
	if err != nil {
		t.Fatalf("failed to add counter with attributes: %v", err)
	}
}

func TestMultipleMetricsFromSameMeter(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	err := meter.AddCounter(ctx, "counter1", "First counter", "1", 1.0)
	if err != nil {
		t.Fatalf("failed to add counter1: %v", err)
	}

	err = meter.AddCounter(ctx, "counter2", "Second counter", "1", 2.0)
	if err != nil {
		t.Fatalf("failed to add counter2: %v", err)
	}

	err = meter.AddGauge(ctx, "gauge1", "First gauge", "celsius", 25.0)
	if err != nil {
		t.Fatalf("failed to add gauge: %v", err)
	}

	err = meter.AddGauge(ctx, "histogram1", "First histogram", "ms", 100.0)
	if err != nil {
		t.Fatalf("failed to add histogram: %v", err)
	}
}

func TestCounterWithMultipleAttributes(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	err := meter.AddCounter(ctx, "http.requests", "HTTP requests counter", "1", 1.0,
		NewAttribute("method", "GET"),
		NewAttribute("status", 200),
		NewAttribute("path", "/api/users"),
		NewAttribute("success", true),
	)
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}

	err = meter.AddCounter(ctx, "http.requests", "HTTP requests counter", "1", 1.0,
		NewAttribute("method", "POST"),
		NewAttribute("status", 201),
		NewAttribute("path", "/api/users"),
		NewAttribute("success", true),
	)
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}

	err = meter.AddCounter(ctx, "http.requests", "HTTP requests counter", "1", 1.0,
		NewAttribute("method", "GET"),
		NewAttribute("status", 404),
		NewAttribute("path", "/api/missing"),
		NewAttribute("success", false),
	)
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}
}

func TestHistogramForLatencies(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	latencies := []float64{50.5, 100.2, 150.8, 200.1, 75.3, 125.9}
	for _, latency := range latencies {
		err := meter.AddGauge(ctx, "http.request.duration", "HTTP request duration", "ms", latency,
			NewAttribute("endpoint", "/api/users"),
			NewAttribute("method", "GET"),
		)
		if err != nil {
			t.Fatalf("failed to add gauge: %v", err)
		}
	}
}

func TestGaugeForMetrics(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	err := meter.AddGauge(ctx, "system.cpu.usage", "CPU usage percentage", "%", 45.5, NewAttribute("host", "server1"))
	if err != nil {
		t.Fatalf("failed to add cpu gauge: %v", err)
	}

	err = meter.AddGauge(ctx, "system.memory.usage", "Memory usage", "bytes", 2048000000, NewAttribute("host", "server1"))
	if err != nil {
		t.Fatalf("failed to add memory gauge: %v", err)
	}

	err = meter.AddGauge(ctx, "system.cpu.usage", "CPU usage percentage", "%", 67.8, NewAttribute("host", "server2"))
	if err != nil {
		t.Fatalf("failed to add cpu gauge: %v", err)
	}

	err = meter.AddGauge(ctx, "system.memory.usage", "Memory usage", "bytes", 4096000000, NewAttribute("host", "server2"))
	if err != nil {
		t.Fatalf("failed to add memory gauge: %v", err)
	}
}

func TestUpDownCounterForConnections(t *testing.T) {
	ctx := context.Background()
	meter := NewNoOpMeter()

	err := meter.AddCounter(ctx, "active.connections", "Active connections", "connections", 10.0, NewAttribute("type", "websocket"))
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}

	err = meter.AddCounter(ctx, "active.connections", "Active connections", "connections", 5.0, NewAttribute("type", "websocket"))
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}

	err = meter.AddCounter(ctx, "active.connections", "Active connections", "connections", 8.0, NewAttribute("type", "http"))
	if err != nil {
		t.Fatalf("failed to add counter: %v", err)
	}
}
