package metric_test

import (
	"context"
	"fmt"
	"time"

	"github.com/bruno303/go-toolkit/pkg/metric"
)

// Example_basicUsage demonstrates basic metric operations
func Example_basicUsage() {
	ctx := context.Background()

	// Get the global meter (returns no-op meter if not configured)
	meter := metric.GetMeter()

	// Add a counter with attributes
	_ = meter.AddCounter(ctx,
		"http.requests.total",
		"Total number of HTTP requests",
		"1",
		1.0,
		metric.NewAttribute("method", "GET"),
		metric.NewAttribute("status", 200),
	)

	// Add a gauge for latency tracking
	_ = meter.AddGauge(ctx,
		"http.request.duration",
		"HTTP request duration",
		"ms",
		150.5,
		metric.NewAttribute("endpoint", "/api/users"),
	)

	// Add a gauge for current temperature
	_ = meter.AddGauge(ctx,
		"system.temperature",
		"System temperature",
		"celsius",
		45.2,
		metric.NewAttribute("component", "cpu"),
	)

	fmt.Println("Metrics recorded successfully")
	// Output: Metrics recorded successfully
}

// Example_setupOpenTelemetry demonstrates setting up OpenTelemetry metrics
// NOTE: This example is commented out as it requires proper OpenTelemetry setup
/*
func Example_setupOpenTelemetry() {
	ctx := context.Background()

	// Setup OpenTelemetry metrics
	shutdown, err := metric.SetupOTelMetrics(ctx, metric.Config{
		ApplicationName:    "my-service",
		ApplicationVersion: "1.0.0",
		Environment:        "production",
		Enabled:            true,
		Port:               9090,
		Path:               "/metrics",
	})
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	// Now use the global meter
	meter := metric.GetMeter()
	_ = meter.AddCounter(ctx, "app.events", "Application events", "1", 1.0, metric.NewAttribute("type", "user_login"))

	fmt.Println("OpenTelemetry metrics setup complete")
	// Output: OpenTelemetry metrics setup complete
}
*/


// Example_httpMetrics demonstrates tracking HTTP metrics
func Example_httpMetrics() {
	ctx := context.Background()
	meter := metric.GetMeter()

	// Simulate handling a request
	start := time.Now()

	// ... process request ...
	time.Sleep(10 * time.Millisecond) // Simulate processing

	duration := float64(time.Since(start).Milliseconds())

	// Record metrics
	_ = meter.AddCounter(ctx,
		"http.server.requests",
		"Total HTTP requests received",
		"1",
		1.0,
		metric.NewAttribute("method", "GET"),
		metric.NewAttribute("path", "/api/users"),
		metric.NewAttribute("status", 200),
	)

	_ = meter.AddGauge(ctx,
		"http.server.duration",
		"HTTP request duration",
		"ms",
		duration,
		metric.NewAttribute("method", "GET"),
		metric.NewAttribute("path", "/api/users"),
	)

	fmt.Println("HTTP metrics recorded")
	// Output: HTTP metrics recorded
}

// Example_databaseMetrics demonstrates tracking database metrics
func Example_databaseMetrics() {
	ctx := context.Background()
	meter := metric.GetMeter()

	// Track a query
	_ = meter.AddCounter(ctx,
		"db.queries.total",
		"Total database queries",
		"1",
		1.0,
		metric.NewAttribute("operation", "SELECT"),
		metric.NewAttribute("table", "users"),
	)

	_ = meter.AddGauge(ctx,
		"db.query.duration",
		"Database query duration",
		"ms",
		25.5,
		metric.NewAttribute("operation", "SELECT"),
		metric.NewAttribute("table", "users"),
	)

	_ = meter.AddGauge(ctx,
		"db.connections.active",
		"Active database connections",
		"1",
		10.0,
		metric.NewAttribute("pool", "primary"),
	)

	fmt.Println("Database metrics recorded")
	// Output: Database metrics recorded
}

// Example_businessMetrics demonstrates tracking business metrics
func Example_businessMetrics() {
	ctx := context.Background()
	meter := metric.GetMeter()

	// Record business events
	_ = meter.AddCounter(ctx,
		"business.orders.total",
		"Total orders processed",
		"1",
		1.0,
		metric.NewAttribute("status", "completed"),
		metric.NewAttribute("region", "us-west"),
	)

	_ = meter.AddCounter(ctx,
		"business.revenue.total",
		"Total revenue",
		"USD",
		99.99,
		metric.NewAttribute("currency", "USD"),
		metric.NewAttribute("region", "us-west"),
	)

	_ = meter.AddGauge(ctx,
		"business.users.active",
		"Currently active users",
		"1",
		1523.0,
		metric.NewAttribute("tier", "premium"),
	)

	fmt.Println("Business metrics recorded")
	// Output: Business metrics recorded
}

// Example_systemMetrics demonstrates tracking system metrics
func Example_systemMetrics() {
	ctx := context.Background()
	meter := metric.GetMeter()

	// Record system metrics
	_ = meter.AddGauge(ctx,
		"system.cpu.usage",
		"CPU usage percentage",
		"%",
		45.5,
		metric.NewAttribute("host", "server-01"),
		metric.NewAttribute("core", 0),
	)

	_ = meter.AddGauge(ctx,
		"system.memory.usage",
		"Memory usage",
		"bytes",
		4294967296, // 4GB in bytes
		metric.NewAttribute("host", "server-01"),
		metric.NewAttribute("type", "used"),
	)

	_ = meter.AddCounter(ctx,
		"system.disk.io",
		"Disk I/O operations",
		"1",
		1000.0,
		metric.NewAttribute("host", "server-01"),
		metric.NewAttribute("operation", "read"),
		metric.NewAttribute("device", "sda1"),
	)

	fmt.Println("System metrics recorded")
	// Output: System metrics recorded
}

// Example_customMeter demonstrates using a custom meter
func Example_customMeter() {
	ctx := context.Background()

	// Create a custom no-op meter
	customMeter := metric.NewNoOpMeter()
	metric.SetMeter(customMeter)

	// Use the custom meter
	meter := metric.GetMeter()
	_ = meter.AddCounter(ctx, "custom.events", "Custom events", "1", 1.0)

	fmt.Println("Custom meter configured")
	// Output: Custom meter configured
}

// Example_multipleAttributes demonstrates using multiple attributes effectively
func Example_multipleAttributes() {
	ctx := context.Background()
	meter := metric.GetMeter()

	// Record with multiple dimensions
	_ = meter.AddCounter(ctx,
		"api.calls",
		"API calls",
		"1",
		1.0,
		metric.NewAttribute("service", "user-service"),
		metric.NewAttribute("method", "GET"),
		metric.NewAttribute("endpoint", "/users/:id"),
		metric.NewAttribute("status_code", 200),
		metric.NewAttribute("region", "us-east-1"),
		metric.NewAttribute("environment", "production"),
		metric.NewAttribute("version", "v2"),
		metric.NewAttribute("authenticated", true),
	)

	fmt.Println("Multi-attribute metrics recorded")
	// Output: Multi-attribute metrics recorded
}

// Example_errorTracking demonstrates tracking errors with metrics
func Example_errorTracking() {
	ctx := context.Background()
	meter := metric.GetMeter()

	// Record different error types
	_ = meter.AddCounter(ctx,
		"app.errors.total",
		"Total application errors",
		"1",
		1.0,
		metric.NewAttribute("type", "database_error"),
		metric.NewAttribute("severity", "high"),
		metric.NewAttribute("component", "user-repository"),
	)

	_ = meter.AddCounter(ctx,
		"app.errors.total",
		"Total application errors",
		"1",
		1.0,
		metric.NewAttribute("type", "validation_error"),
		metric.NewAttribute("severity", "low"),
		metric.NewAttribute("component", "user-controller"),
	)

	_ = meter.AddCounter(ctx,
		"app.errors.total",
		"Total application errors",
		"1",
		1.0,
		metric.NewAttribute("type", "timeout_error"),
		metric.NewAttribute("severity", "medium"),
		metric.NewAttribute("component", "external-api-client"),
	)

	fmt.Println("Error metrics recorded")
	// Output: Error metrics recorded
}
