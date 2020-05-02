package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type MetricsCollector struct {
	AccessToken              string `json:"access_token"`
	OpsManCertificateListUrl string
	Certificates             []struct {
		Configurable      bool      `json:"configurable"`
		IsCa              bool      `json:"is_ca"`
		PropertyReference string    `json:"property_reference"`
		PropertyType      string    `json:"property_type"`
		ProductGUID       string    `json:"product_guid"`
		Location          string    `json:"location"`
		VariablePath      string    `json:"variable_path"`
		Issuer            string    `json:"issuer"`
		ValidFrom         time.Time `json:"valid_from"`
		ValidUntil        time.Time `json:"valid_until"`
	} `json:"certificates"`
}

// When you encounter error, increment the counter and display the error on stdout
func IncrementErrorCounter(error string) error {
	Errorf(error)
	ErrorTotal.WithLabelValues(cmdOptions.Environment).Inc()
	return fmt.Errorf(error)
}

//
func (m *MetricsCollector) collector() {
	// Authenticate the request
	err := m.authenticate()
	if err != nil {
		return
	}

	// Get the data of list of certs
	m.OpsManCertificateListUrl = fmt.Sprintf("https://%s/api/v0/deployed/certificates", cmdOptions.OpsManHostname)
	c, err := m.opsmanRequestHandler()
	if err != nil {
		IncrementErrorCounter(fmt.Sprintf("error during extraction of certs from opsman: %v", err))
		return
	}

	// Load the information to the cert struct
	err = json.Unmarshal(c, &m)
	if err != nil {
		IncrementErrorCounter(fmt.Sprintf("error in unmarshal of certificates: %v", err))
		return
	}

	// Convert to Prometheus metrics
	m.metric()
}

func (m *MetricsCollector) metric() {
	for _, c := range m.Certificates {
		CertExpirySeconds.WithLabelValues(
			cmdOptions.Environment,
			strconv.FormatBool(c.IsCa),
			strconv.FormatBool(c.Configurable),
			c.PropertyReference,
			c.ProductGUID,
			c.Location,
			c.VariablePath,
			c.Issuer,
			fmt.Sprintf("%v", c.ValidFrom),
			fmt.Sprintf("%v", c.ValidUntil)).Set(time.Until(c.ValidUntil).Seconds())

		CertExpiryUnixTimestamp.WithLabelValues(
			cmdOptions.Environment,
			strconv.FormatBool(c.IsCa),
			strconv.FormatBool(c.Configurable),
			c.PropertyReference,
			c.ProductGUID,
			c.Location,
			c.VariablePath,
			c.Issuer,
			fmt.Sprintf("%v", c.ValidFrom),
			fmt.Sprintf("%v", c.ValidUntil)).Set(float64(c.ValidUntil.Unix()))
	}
}
