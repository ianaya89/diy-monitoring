package main

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// logger ships records through the OTel log bridge to the collector.
// The bridge auto-injects trace_id/span_id from the request context.
var logger = otelslog.NewLogger("rolldice")

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.InfoContext(r.Context(), "http_request",
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("uri", r.URL.Path),
			slog.String("duration", time.Since(start).String()),
		)
	})
}
