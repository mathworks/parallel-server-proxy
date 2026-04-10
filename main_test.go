// Copyright 2024 The MathWorks, Inc.
package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArgsCertificateArgument(t *testing.T) {
	testArgs := []string{programName, "-" + hostFlag, "example.com", "-" + portFlag, "1234", "-" + certificateFlag, "cert.json", "-" + verboseFlag}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	expectedArgs := proxyInputs{
		Port:              1234,
		Host:              "example.com",
		CertificateFile:   "cert.json",
		Verbose:           true,
		DisableMututalTLS: false,
		Quiet:             false,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedArgs, *args)
}

func TestParseArgsDisableMututalTLSArgument(t *testing.T) {
	testArgs := []string{programName, "-" + hostFlag, "example.com", "-" + portFlag, "1234", "-" + disableMutualTLSFlag, "-" + quietFlag}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	expectedArgs := proxyInputs{
		Port:              1234,
		Host:              "example.com",
		CertificateFile:   "",
		Verbose:           false,
		DisableMututalTLS: true,
		Quiet:             true,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedArgs, *args)
}

func TestParseArgsHelp(t *testing.T) {
	testArgs := []string{programName, "-h"}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	assert.NoError(t, err)
	require.Nil(t, args, "command args should be nil when help flag specified")
}

func TestParseArgsNonNumericPort(t *testing.T) {
	testArgs := []string{programName, "-" + portFlag, "text", "-" + disableMutualTLSFlag}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	require.Error(t, err, "expected error when non-numeric port specified")
	require.Nil(t, args, "command args should be nil when non-numeric port specified")
}

func TestParseArgsInvalidPort(t *testing.T) {
	testArgs := []string{programName, "-" + portFlag, "222123212", "-" + disableMutualTLSFlag}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	require.Error(t, err, "expected error when invalid port specified")
	require.Nil(t, args, "command args should be nil when invalid port specified")
}

func TestParseArgsMissingCert(t *testing.T) {
	testArgs := []string{programName}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	require.Error(t, err, "expected error when no certificate or disable mutual TLS specified")
	require.Nil(t, args, "command args should be nil when no certificate or disable mutual TLS specified")
}

func TestParseArgsVerboseAndQuiet(t *testing.T) {
	testArgs := []string{programName, "-" + verboseFlag, "-" + quietFlag}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	require.Error(t, err, "expected error when both verbose and quiet output requested")
	require.Nil(t, args, "command args should be nil when both verbose and quiet output requested")
}

func TestParseArgsmTLSAndTCP(t *testing.T) {
	testArgs := []string{programName, "-" + disableMutualTLSFlag, "-" + certificateFlag, "cert.json"}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	args, err := parseArgs(flagSet, testArgs)

	require.Error(t, err, "expected error when both mTLS and TCP requested")
	require.Nil(t, args, "command args should be nil when both mTLS and TCP requested")
}
