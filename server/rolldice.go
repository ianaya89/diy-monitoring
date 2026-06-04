package main

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer = otel.Tracer("rolldice")
	meter  = otel.Meter("rolldice")

	rollCounter metric.Int64Counter
)

// setupInstruments must run after the meter provider is registered,
// otherwise instruments bind to the no-op provider.
func setupInstruments() error {
	var err error
	rollCounter, err = meter.Int64Counter(
		"dice.rolls",
		metric.WithDescription("Number of dice rolls by outcome"),
		metric.WithUnit("{roll}"),
	)
	return err
}

func rolldice(w http.ResponseWriter, r *http.Request) {
	performRoll(w, r, "roll", 0)
}

func rolldiceSlow(w http.ResponseWriter, r *http.Request) {
	performRoll(w, r, "roll.slow", 3*time.Second)
}

func rolldiceError(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "roll.error")
	defer span.End()

	message := "Internal Server Error"
	err := errors.New(message)
	span.RecordError(err)
	span.SetStatus(codes.Error, message)

	http.Error(w, message, http.StatusInternalServerError)
}

func performRoll(w http.ResponseWriter, r *http.Request, spanName string, delay time.Duration) {
	_, span := tracer.Start(r.Context(), spanName)
	defer span.End()

	roll := 1 + rand.Intn(6)
	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	rollCounter.Add(r.Context(), 1, metric.WithAttributes(rollValueAttr))

	if delay > 0 {
		time.Sleep(delay)
	}

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
