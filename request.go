package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// The below functions fetch data from the URL
func fetch(method string, url string, headers map[string]string) ([]byte, error) {

	var contents []byte

	// Create new http request
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return contents, fmt.Errorf("encountered error when sending new request: %v", err)
	}

	// copy headers
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	var tlsConfig *tls.Config
	if cmdOptions.SkipSsl { // Skip SSL stuffs
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else { // add CA if provided
		cert, err := ioutil.ReadFile("server.crt")
		if err != nil {
			Fatalf("Couldn't load file: %v", err)
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(cert)
		tlsConfig = &tls.Config{
			RootCAs: certPool,
		}
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// perform request
	client := &http.Client{
		Transport: transport,
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return contents, fmt.Errorf("encountered error when requesting the data from http: %v", err)
	}

	// read response
	defer resp.Body.Close()
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return contents, fmt.Errorf("encountered error when reading the data from http: %v", err)
	}

	// Check if the response is 200
	if resp.StatusCode != http.StatusOK {
		return contents, fmt.Errorf("encountered invalid status code from http: %v", resp.StatusCode)
	}

	return contents, nil
}

// the below function does the get from the URL
func get(url string, headers map[string]string) ([]byte, error) {
	return fetch("GET", url, headers)
}

// Extract the access token from ops man
func (m *MetricsCollector) authenticate() error {
	Debugf("Extracting the access token for the ops manager: %s", cmdOptions.OpsManHostname)

	// Parse the authentication URL
	URL, err := url.Parse("https://" + cmdOptions.OpsManHostname + "/uaa/oauth/token")
	if err != nil {
		return IncrementErrorCounter(fmt.Sprintf("error in parsing authentication url: %v", err))
	}

	// Add Parameter values for URL for authentication
	parameters := url.Values{}
	parameters.Add("grant_type", "password")
	parameters.Add("username", cmdOptions.OpsManUsername)
	parameters.Add("password", cmdOptions.OpsManPassword)
	URL.RawQuery = parameters.Encode()

	// Headers.
	headers := make(map[string]string)
	headers["Accept"] = "application/json;charset=utf-8"
	headers["Authorization"] = "Basic b3BzbWFuOg=="

	// Fetch the authentication data
	content, err := get(URL.String(), headers)
	if err != nil {
		return IncrementErrorCounter(fmt.Sprintf("error during token extraction: %v", err))
	}

	// Extract the auth token from the content received
	err = json.Unmarshal(content, &m)
	if err != nil {
		return IncrementErrorCounter(fmt.Sprintf("error in json unmarshal of uaa response: %v", err))

	}

	return nil
}

// Build the header and send request to the ops manager for the data
func (m *MetricsCollector) opsmanRequestHandler() ([]byte, error) {
	Debugf("Extracting the certificate list from the URL: %s", m.OpsManCertificateListUrl)

	// Headers.
	headers := make(map[string]string)
	headers["Accept"] = "application/json;charset=utf-8"
	headers["Authorization"] = "Bearer " + m.AccessToken

	// Fetch the data
	return get(m.OpsManCertificateListUrl, headers)
}
