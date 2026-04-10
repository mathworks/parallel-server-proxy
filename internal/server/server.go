// Copyright 2024 The MathWorks, Inc.
package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
)

const wildcardHost = "*"

// Build the client connect url
func ConstructConnectURL(protocol string, hostname string, port int) string {
	if hostname == "" {
		hostname = wildcardHost
	}
	u := url.URL{
		Scheme: protocol,
		Host:   net.JoinHostPort(hostname, fmt.Sprint(port)),
	}
	return u.String()
}

// Return a server socket on the specified host and port. If a json certificate file is provided, it will be used to create a TLS listener.
func CreateListener(host string, port int, certificateFile string) (net.Listener, error) {
	proxyUrl := net.JoinHostPort(host, fmt.Sprint(port))
	if certificateFile != "" {
		tlsConfig, err := loadCertificateFile(certificateFile)
		if err != nil {
			return nil, err
		}
		ln, err := tls.Listen("tcp", proxyUrl, tlsConfig)
		return ln, err
	}
	ln, err := net.Listen("tcp", proxyUrl)
	return ln, err
}

// Read the given file as raw bytes
func readFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()
	return io.ReadAll(file)
}

// Data extracted from the json certificate file
type certificateFileData struct {
	CertPEM string `json:"clientCert"`
	KeyPEM  string `json:"clientKey"`
	CAPEM   string `json:"serverCert"`
}

// Load the json certificate file and return a corresponding TLS configuration for creating a secure server socket
func loadCertificateFile(filename string) (*tls.Config, error) {
	certificateBytes, err := readFromFile(filename)
	if err != nil {
		return nil, err
	}

	var certificate certificateFileData
	err = json.Unmarshal(certificateBytes, &certificate)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling certificate file: %v", err)
	}

	tlsCert, err := tls.X509KeyPair([]byte(certificate.CertPEM), []byte(certificate.KeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to create tls key pair: %s", err.Error())
	}

	// Create a pool with the self-signed CA certificate
	caCert, err := decodeCert(certificate.CAPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to create ca cert: %s", err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// Decode the given PEM encoded certificate string
func decodeCert(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return caCert, nil
}
