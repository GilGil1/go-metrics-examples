# go-metrics-examples
Examples of instrumenting go program metrics

In this repository we shall present code examples to instrument and expose metrics to an external observability tool.
We collect the follwing metrics:
* Golang internal metrics in the runtime package
* Golang runtime/metrics package. Theis package exposes process metrics in a generic method, with metadata on each metric.


The metrics are exposed in Promethus format - Open Metrics (https://github.com/OpenObservability/OpenMetrics)

The instruimentation is done in the follwing libraries:
Prometheus Client library - https://github.com/prometheus/client_golang
Open Telemetry Golang instrumentation - https://opentelemetry.io/docs/instrumentation/go/

## Project structure
Those examples are limited to counters and gauges. This is the list of examples under the following folders:
* folder /cmd/basic - example to expose runtime package metric using Prometheus client
* folder /cmd/otel - example to expose runtime/metrics using Open telemetry instrumentation
* folder /cmd/prometheus - example to expose runtime/metrics using Prometheus client

### Note
* Each library enfirces different metrics naming conventions and location of lables


# Integrating the project with observability tool
The examples above exposes the metrics in Prometheus format. To collect and display them, you can 
* Set a Prometheus server on the endpoint for collecting (scraping) the metrics. The metrics will be stored in the Prometheus server
https://prometheus.io/docs/prometheus/latest/getting_started/
* Display the metrics with Grafana and set allert on them. In that case you need to set the Promethus server as a data source for Grafana.
https://grafana.com/docs/grafana/latest/datasources/prometheus/


## Metrics are exposed in this format
```
# HELP cgo_go_to_c_calls_calls_total Count of calls made from Go to C by the current process.
# TYPE cgo_go_to_c_calls_calls_total counter
cgo_go_to_c_calls_calls_total{Namespace="cgo",Subsystem="",Units="calls",otel_scope_name="github.com/open-telemetry/opentelemetry-go/example/prometheus",otel_scope_version=""} 1
# HELP cpu_classes_gc_mark_assist_cpu_seconds_total Estimated total CPU time goroutines spent performing GC tasks to assist the GC and prevent it from falling behind the application. This metric is an overestimate, and not directly comparable to system CPU time measurements. Compare only with other /cpu/classes metrics.
# TYPE cpu_classes_gc_mark_assist_cpu_seconds_total counter
```
