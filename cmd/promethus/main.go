package main

import (
	"fmt"
	"go-metrics-examples/internal/metadata"
	"log"
	"net/http"
	"runtime/metrics"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main function to register metrics and start prometheus server
func main() {

	// Register metrics and retrieve the values in prometheus client
	addMetricsToPrometheusRegistry()

	// Setup the prometheus handler
	log.Printf("serving metrics at localhost%s%s", metadata.MetricsEndpointPort, metadata.MetricsPath)
	http.Handle(metadata.MetricsPath, promhttp.Handler())
	http.ListenAndServe(metadata.MetricsEndpointPort, nil)
}

// addMetricsToPrometheusRegistry function to add metrics to prometheus registry
func addMetricsToPrometheusRegistry() {
	// Get descriptions for all supported metrics.
	metricsMeta := metrics.All()
	// Register metrics and retrieve the values in prometheus client
	for i := range metricsMeta {
		meta := metricsMeta[i]
		opts := getMetricsOptions(metricsMeta[i])
		if meta.Cumulative {
			// Register as a counter
			funcCounter := prometheus.NewCounterFunc(prometheus.CounterOpts(opts), func() float64 {
				return metadata.GetSingleMetricFloat(meta.Name)
			})
			prometheus.MustRegister(funcCounter)
		} else {
			// Register as a gauge
			funcGauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts(opts), func() float64 {
				return metadata.GetSingleMetricFloat(meta.Name)
			})
			prometheus.MustRegister(funcGauge)
		}
	}
}

// getMetricsOptions function to get prometheus options for a metric
func getMetricsOptions(metric metrics.Description) prometheus.Opts {
	tokens := strings.Split(metric.Name, "/")
	fmt.Printf("%s\n", metric.Name)
	if len(tokens) < 2 {
		return prometheus.Opts{}
	}
	nameTokens := strings.Split(tokens[len(tokens)-1], ":")
	// create a unique name for metric, that will be its primary key on the registry
	validName := normalizePrometheusName(strings.Join(nameTokens[:2], "_"))
	subsystem := metadata.GetMetricSubsystemName(metric)

	units := nameTokens[1]
	help := fmt.Sprintf("Units:%s, %s", units, metric.Description)
	opts := prometheus.Opts{
		Namespace: tokens[1],
		Subsystem: subsystem,
		Name:      validName,
		Help:      help,
	}
	return opts
}

// normalizePrometheusName function to normalize prometheus name
func normalizePrometheusName(name string) string {
	return strings.TrimSpace(strings.ReplaceAll(name, "-", "_"))
}
