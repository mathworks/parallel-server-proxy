// Copyright 2024 The MathWorks, Inc.
package server

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCertificate(t *testing.T) {
	// Set up a test certificate file
	testCertFile := "testdata/cert.json"
	tlsConfig, err := loadCertificateFile(testCertFile)
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)

	// Verify the TLS configuration
	assert.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
	assert.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
	assert.Len(t, tlsConfig.Certificates, 1)
	assert.NotNil(t, tlsConfig.ClientCAs)
}

func TestListener(t *testing.T) {
	// Test creating a listener without a certificate file
	ln, err := CreateListener("localhost", 30808, "")
	require.NoError(t, err)
	assert.NotNil(t, ln)
	ln.Close()
}

func TestSecureListener(t *testing.T) {
	// Test creating a listener with a certificate file
	testCertFile := "testdata/cert.json"
	ln, err := CreateListener("localhost", 30809, testCertFile)
	require.NoError(t, err)
	assert.NotNil(t, ln)
	ln.Close()
}
