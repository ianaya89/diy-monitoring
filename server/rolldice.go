package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	tracer = otel.Tracer("rolldice")
)

func rolldice(w http.ResponseWriter, r *http.Request) {
	performRoll(w, r, "roll", 0)
}

func rolldiceSlow(w http.ResponseWriter, r *http.Request) {
	performRoll(w, r, "roll", 3*time.Second)
}

func rolldiceError(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	message := "Internal Server Error"
	rollValueAttr := attribute.String("error", message)
	span.SetAttributes(rollValueAttr)

	http.Error(w, message, http.StatusInternalServerError)
}

func performRoll(w http.ResponseWriter, r *http.Request, spanName string, delay time.Duration) {
	_, span := tracer.Start(r.Context(), spanName)
	defer span.End()

	roll := 1 + rand.Intn(6)
	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)

	if delay > 0 {
		time.Sleep(delay)
	}

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
