package main

import (
	"fmt"
	"go-metrics-examples/internal/metadata"
	"net/http"
	"runtime/metrics"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main function to register metrics and start prometheus server
func main() {

	// Get descriptions for all supported metrics.
	metricsMeta := metrics.All()

	// Register metrics and retreivinbg the values in prometgheus client
	addMetricsToPrometheusRegistry(metricsMeta)

	// Setup the prometheus metrics endpoint on port 2112,
	// to use it run http://ipadress:2112/metrics or set prometheus to that address
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

// addMetricsToPrometheusRegistry function to add metrics to prometheus registry
func addMetricsToPrometheusRegistry(metricsMeta []metrics.Description) {
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
	subsystem := getSubsystemName(metric)

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

// getSubsystemName function to get subsystem name from metric name
func getSubsystemName(metric metrics.Description) string {
	tokens := strings.Split(metric.Name, "/")
	if len(tokens) < 2 {
		return ""
	}
	if len(tokens) > 3 {
		subsystemTokens := tokens[2 : len(tokens)-1]
		subsystem := strings.Join(subsystemTokens, "_")
		subsystem = normalizePrometheusName(subsystem)
		return subsystem
	}
	return ""
}

// normalizePrometheusName function to normalize prometheus name
func normalizePrometheusName(name string) string {
	return strings.TrimSpace(strings.ReplaceAll(name, "-", "_"))
}
