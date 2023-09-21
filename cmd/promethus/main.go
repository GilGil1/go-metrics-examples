package main

import (
	"fmt"
	"log"
	"runtime/metrics"
)

func main() {

	PrintfRuntimeMetricsPackageStyle()
}

// PrintfRuntimeMetricsPackageStyle... This function demonstrates the usage of runtime/metrics package.
// It demonstrates how to get all metrics parameters in few calls, no need to adjust if a new metric is added
func PrintfRuntimeMetricsPackageStyle() {
	metricsMetadata := metrics.All()

	// retrieve metrics metadata including names, type and description and values
	for idx := range metricsMetadata {

		// get metric metadata
		metadata := metricsMetadata[idx]
		fmt.Printf("\nMetric Name=%-50s | Kind=%v | Comulative=%t | ",
			metadata.Name,
			metadata.Kind,
			metadata.Cumulative)

		// retrieve current metric value. In our case we retrieve a single metric at a time
		// however, we can retrieve multiple metrics in one call, for efficiency
		samples := make([]metrics.Sample, 1)
		samples[0].Name = metadata.Name // set the metric name to retrieve
		metrics.Read(samples)           // Reading samples

		// Printing each sample according to its type defined in the metric metadata
		switch samples[0].Value.Kind() {
		case metrics.KindUint64: // handle integer value
			fmt.Printf("value=%d\n", samples[0].Value.Uint64())

		case metrics.KindFloat64: // handle float value
			fmt.Printf("value=%f\n", samples[0].Value.Float64())

		case metrics.KindFloat64Histogram: // handle histograms
			histogram := samples[0].Value.Float64Histogram()

			// There are less buckets than count, first count is from bucket 0 to 1
			for idx := range histogram.Counts {
				fmt.Printf("\n* Bucket %10.1f - %10.1f : Counts %d", histogram.Buckets[idx], histogram.Buckets[idx+1], histogram.Counts[idx])
			}
			fmt.Printf("\n")

		case metrics.KindBad:
			log.Fatalf("unsupported metric type")
		}
		fmt.Printf("Description=%s\n", metadata.Description)
	}
}