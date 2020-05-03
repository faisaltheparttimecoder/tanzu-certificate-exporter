package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "vmware_tanzu_cert_exporter"

var (
	// ErrorTotal is a prometheus counter that indicates the total number of unexpected errors encountered by the program
	ErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "error_total",
			Help:      fmt.Sprintf("Errors generated from the program %s", programName),
		},
		[]string{"env"},
	)

	// CertExpirySeconds is a prometheus gauge that indicates the number of seconds until certificates on disk expires.
	CertExpirySeconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cert_expires_in_seconds",
			Help:      "Number of seconds til the cert expires.",
		},
		[]string{"env", "is_ca", "configurable", "property_reference", "product_guid",
			"location", "variable_path", "issuer", "valid_from", "valid_until"},
	)

	// CertExpiryUnixTimestamp is a prometheus gauge that indicates the time in unix timestamp format.
	CertExpiryUnixTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cert_expiry_in_unix_timestamp",
			Help:      "Unix timestamp of when the certificate expires.",
		},
		[]string{"env", "is_ca", "configurable", "property_reference", "product_guid",
			"location", "variable_path", "issuer", "valid_from", "valid_until"},
	)
)

func init() {
	prometheus.MustRegister(ErrorTotal)
	prometheus.MustRegister(CertExpirySeconds)
	prometheus.MustRegister(CertExpiryUnixTimestamp)
}
