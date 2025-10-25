# go-toolkit

Go toolkit

## Generating a new tag

```shell
git tag -a v1.0.0
git push origin v1.0.0
```

## Logging

### Singleton Logger Configuration

Configure a single logger instance that will be used for all logging operations:

```go
// Create a custom logger
customLogger := log.NewSlogAdapter(log.SlogAdapterOpts{
  Level:      log.LevelDebug,
  FormatJson: true,
  Name:       "singleton-logger",
})

// Configure logging with singleton mode
log.ConfigureLogging(log.LogConfig{
  Type: log.LogTypeSingleton,
  SingletonLogConfig: log.SingletonLogConfig{
    Logger: customLogger,
  }
})

// Use the configured logger
logger := log.Log() // Returns the singleton logger
logger.Info(ctx, "This uses the singleton logger")
```

### Multiple Logger Configuration

Configure a factory function to create different loggers on demand:

```go
// Create a factory function
loggerFactory := func(name string) log.Logger {
  return log.NewSlogAdapter(log.SlogAdapterOpts{
    Level:      log.LevelInfo,
    FormatJson: false,
    Name:       name,
  })
}

// Configure logging with multiple mode
log.ConfigureLogging(log.LogConfig{
  Type: log.LogTypeMultiple,
  MultipleLogConfig: log.MultipleLogConfig{
    Factory: loggerFactory,
  },
  Levels: map[string]log.Level{
    "service.user":    log.LevelDebug,
    "service.payment": log.LevelWarn,
    "database":        log.LevelError,
  },
})

// Create different loggers
userLogger := log.NewLogger("service.user")       // Debug level
paymentLogger := log.NewLogger("service.payment") // Warn level
dbLogger := log.NewLogger("database")             // Error level

userLogger.Debug(ctx, "User operation")
paymentLogger.Warn(ctx, "Payment warning")
dbLogger.Error(ctx, "Database error", err)
```

### Manual Logger Setup (Alternative)

You can also set loggers manually without using ConfigureLogging:

```go
log.SetLogger(
  log.NewSlogAdapter(
    log.SlogAdapterOpts{
      Level:                 log.LevelDebug,
      FormatJson:            true,
      ExtractAdditionalInfo: func(context.Context) []any { return nil },
      Name:                  "default",
    },
  ),
)
```

## Trace

```go
trace.SetTracer(otel.NewOtelTracerAdapter())
shutdown, err := trace.SetupOTelSDK(ctx, trace.Config{
  ApplicationName:    "app",
  ApplicationVersion: "0.0.1",
  Endpoint:           "localhost:4317",
})
if err != nil {
  ....
}
defer shutdown()
```

## Http Middleware

```go
var h http.Handler = ......

chain, err := http.NewChain(
  http.LogMiddleware(),
  ...
  http.NewMiddleware(h),
)

server.Handle("GET /hello", chain)
```

## Shutdown

```go
shutdown.ConfigureGracefulShutdown()
shutdown.CreateListener(func() {
  // cleanup here
})
shutdown.CreateListener(func() {
  // cleanup here
})
shutdown.AwaitAll()
```

## Arrays

```go
mySlice := []string{"a","b","c"}

array.Contains(mySlice, "b") // true
array.FirstOrNil(mySlice, func(e string) bool { return e == "a" }) // a
array.Map(mySlice, func(e string) int { return len(e) }) // []int{1,1,1}
array.Join(mySlice, "-") // "a-b-c"
array.Filter(mySlice, func(e string) string { return e != "b" }) // []string{"a","c"}
array.Remove(mySlice, func(e) bool { return e == "c" }) // []string{"a","b"}
```
