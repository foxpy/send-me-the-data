package handler

import (
	"log/slog"
	"net/http"
	"time"
)

type loggerMiddleware struct {
	name    string
	handler http.Handler
}

type responseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (l *loggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &responseWriter{http.StatusOK, w}

	begin := time.Now()
	l.handler.ServeHTTP(lw, r)
	duration := time.Since(begin)

	slog.Info(
		"handle HTTP request",
		"server", l.name,
		"method", r.Method,
		"url", r.URL,
		"code", lw.statusCode,
		"duration", duration,
	)
}

func WithLogger(handler http.Handler, serverName string) http.Handler {
	return &loggerMiddleware{serverName, handler}
}
