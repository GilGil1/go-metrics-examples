package main

import (
	"fmt"
	"runtime"
)

// This is an example of how to use the runtime package to get metrics
func main() {
	PrintOldStyleMetrics()
}

// PrintOldStyleMetrics... Printing the metrics based on runtime package. Each of the metric is retrieved in a different method.
// Each metric will have a different function
func PrintOldStyleMetrics() {
	fmt.Printf("Operating system is: %s\n", runtime.GOOS)
	// Get number of go routines, max os threads allocated to process and host number of cpus
	numGoroutines := runtime.NumGoroutine()
	fmt.Printf("The pograms is using %d go routines\n", numGoroutines)

	maxThreads := runtime.GOMAXPROCS(0)
	fmt.Printf("The program is configured to %d max threads", maxThreads)

	numCPUs := runtime.NumCPU()
	fmt.Printf("the host has %d cpus\n", numCPUs)
}
