package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type requestLog struct {
	RemoteAddr string `json:"remote_addr"`
	Method     string `json:"method"`
	URI        string `json:"uri"`
	Duration   string `json:"duration"`
	TraceID    string `json:"trace_id"`
	SpanID     string `json:"span_id"`
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		spanCtx := span.SpanContext()

		next.ServeHTTP(w, r)

		entry := requestLog{
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			URI:        r.URL.Path,
			Duration:   time.Since(start).String(),
			TraceID:    spanCtx.TraceID().String(),
			SpanID:     spanCtx.SpanID().String(),
		}

		line, err := json.Marshal(entry)
		if err != nil {
			log.Printf("log marshal failed: %v\n", err)
			return
		}
		log.Print(string(line))
	})
}
