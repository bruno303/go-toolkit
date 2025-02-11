package http

import (
	"net/http"

	"github.com/bruno303/go-toolkit/pkg/log"
)

type LoggingMiddleware struct {
	BaseMiddleware
	logger log.Logger
}

func (m *LoggingMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.logger.Info(r.Context(), "Request received: %s %s", r.Method, r.URL.Path)
	m.Next(rw, r)
	m.logger.Info(r.Context(), "Request finished: %s %s", r.Method, r.URL.Path)
}

func LogMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: log.NewLogger("log.LoggingMiddleware"),
	}
}
