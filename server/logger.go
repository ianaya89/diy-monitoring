package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type LogEntry struct {
	Timestamp      string `json:"Timestamp"`
	SeverityText   string `json:"SeverityText"`
	Body           string `json:"Body"`
	ServiceName    string `json:"ServiceName"`
	SeverityNumber int    `json:"SeverityNumber"`
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		spanCtx := span.SpanContext()
		traceID := spanCtx.TraceID().String()
		spanID := spanCtx.SpanID().String()

		next.ServeHTTP(w, r)

		logMessage := fmt.Sprintf(`{"remote_addr": "%s", "method": "%s", "uri": "%s", "duration": "%v", "trace_id": "%s", "span_id": "%s"}`,
			r.RemoteAddr, r.Method, r.URL.Path, time.Since(start), traceID, spanID)

		log.Print(logMessage)
	})
}
