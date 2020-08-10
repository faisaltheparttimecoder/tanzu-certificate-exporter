package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

const (
	programVersion = "0.1"
	programName    = "vmware-tanzu-cert-exporter"
	metricURL      = "/metrics"
)

// Start the background thread to monitor the certificate expiry date
func startCertificateMonitoring() {
	m := new(MetricsCollector)
	go func() {
		for {
			// Get the certificate information
			m.collector()
			// Sleep for the next interval to scrap again
			time.Sleep(time.Duration(cmdOptions.Interval) * time.Second)
			// Reset the cache
			ResetMetrics()
		}
	}()
}

// Start the http server and send in the certificate info
func startHttpServer() {
	// Start the process of getting certificate information in the background
	startCertificateMonitoring()

	// Listen and serve prometheus request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Hello Welcome, <br/> You can view the metrics <a href=\"%s\">here</a>", metricURL)
	})
	http.Handle(metricURL, promhttp.Handler())

	Infof("Starting server on port %d, serving data at path %s", cmdOptions.Port, metricURL)
	http.ListenAndServe(fmt.Sprintf(":%d", cmdOptions.Port), nil)
}

func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
