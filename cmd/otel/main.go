package main

import (
	"context"
	"fmt"
	"go-metrics-examples/internal/metadata"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/metrics"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register the prometheus Collector to receive metrics from the go runtime package metrics
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	addMetricsToOTEL(provider)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	// Wait for signal to close the app
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

// serveMetrics function to start prometheus server
func serveMetrics() {
	log.Printf("serving metrics at localhost%s%s", metadata.MetricsEndpointPort, metadata.MetricsPath)
	http.Handle(metadata.MetricsPath, promhttp.Handler())
	err := http.ListenAndServe(metadata.MetricsEndpointPort, nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

// addMetricsToPrometheusRegistry function to add metrics to prometheus registry
func addMetricsToOTEL(provider *metric.MeterProvider) {

	meter := provider.Meter(metadata.OtelScope)

	// Get descriptions for all supported metrics.
	metricsMeta := metrics.All()

	// Register metrics and retrieve the values in prometheus client
	for i := range metricsMeta {
		// Get metric options
		meta := metricsMeta[i]
		opt := getMetricsOptions(metricsMeta[i])
		name := normalizeOtelName(meta.Name)

		// Register metrics per type of metric
		if meta.Cumulative {
			// Register as a counter
			counter, err := meter.Float64ObservableCounter(name, api.WithDescription(meta.Description))
			if err != nil {
				log.Fatal(err)
			}
			_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
				o.ObserveFloat64(counter, metadata.GetSingleMetricFloat(meta.Name), opt)
				return nil
			}, counter)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Register as a gauge
			gauge, err := meter.Float64ObservableGauge(name, api.WithDescription(meta.Description))
			if err != nil {
				log.Fatal(err)
			}
			_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
				o.ObserveFloat64(gauge, metadata.GetSingleMetricFloat(meta.Name), opt)
				return nil
			}, gauge)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// getMetricsOptions function to get metric labels
func getMetricsOptions(metric metrics.Description) api.MeasurementOption {
	tokens := strings.Split(metric.Name, "/")
	if len(tokens) < 2 {
		return nil
	}

	nameTokens := strings.Split(tokens[len(tokens)-1], ":")
	subsystem := metadata.GetMetricSubsystemName(metric)

	// create a unique name for metric, that will be its primary key on the registry
	opt := api.WithAttributes(
		attribute.Key("Namespace").String(tokens[1]),
		attribute.Key("Subsystem").String(subsystem),
		attribute.Key("Units").String(nameTokens[1]),
	)
	return opt
}

// normalizePrometheusName function to normalize prometheus metric name
func normalizeOtelName(name string) string {
	normalizedName := strings.Replace(name, "/", "", 1)
	normalizedName = strings.Replace(normalizedName, ":", "_", -1)
	normalizedName = strings.TrimSpace(strings.ReplaceAll(normalizedName, "/", "_"))
	return normalizedName
}
