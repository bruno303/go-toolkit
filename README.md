# go-toolkit

Go toolkit

## Generating a new tag

```shell
git tag -a v1.0.0
git push origin v1.0.0
```

## Logging

```go
log.SetLogger(
  slog.NewSlogAdapter(
    slog.SlogAdapterOpts{
      Level:                 log.LevelDebug,
      FormatJson:            true,
      ExtractAdditionalInfo: func(context.Context) []any {},
    },
  )
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
