// Copyright 2024 The MathWorks, Inc.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/mathworks/parallelserverproxy/internal/server"
	"github.com/mathworks/parallelserverproxy/internal/socks5proxy"
)

// Define all arguments names and default values
const (
	programName          = "parallelserverproxy"
	certificateFlag      = "certificate"
	disableMutualTLSFlag = "disableMutualTLS"
	hostFlag             = "host"
	portFlag             = "port"
	quietFlag            = "quiet"
	verboseFlag          = "verbose"
	defaultPort          = 1080
)

type proxyInputs struct {
	CertificateFile   string
	DisableMututalTLS bool
	Host              string
	Port              int
	Quiet             bool
	Verbose           bool
}

func main() {
	// Process input command-line arguments
	flagSet := flag.NewFlagSet(programName, flag.ExitOnError)
	args, err := parseArgs(flagSet, os.Args)
	if err != nil {
		fmt.Println(err)
		flagSet.Usage()
		os.Exit(1)
	}
	if args == nil {
		os.Exit(0)
	}

	// Create logger
	var programLevel = new(slog.LevelVar) // Info by default
	if args.Quiet {
		handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: programLevel})
		slog.SetDefault(slog.New(handler))
		log.SetFlags(0)
	} else {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
		slog.SetDefault(slog.New(handler))

		if args.Verbose {
			programLevel.Set(slog.LevelDebug)
		}
	}
	logger := slog.Default()

	// Bind to server socket
	if args.DisableMututalTLS && !args.Quiet {
		fmt.Printf("Warning: Client verification disabled with -%s option.\n", disableMutualTLSFlag)
	}
	ln, err := server.CreateListener(args.Host, args.Port, args.CertificateFile)
	if err != nil {
		errorAndExit(err)
	}
	defer ln.Close()

	// Start SOCKS5 server
	if err := proxyConnections(ln, logger, args); err != nil {
		errorAndExit(err)
	}
}

// Print given error and exit with code 1
func errorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// parseArgs validates and parses command-line arguments
func parseArgs(flagSet *flag.FlagSet, args []string) (*proxyInputs, error) {
	inputs := proxyInputs{}

	flagSet.StringVar(&inputs.CertificateFile, certificateFlag, "", "Location of the certificate json file to authenticate clients with")
	flagSet.BoolVar(&inputs.DisableMututalTLS, disableMutualTLSFlag, false, "Run in insecure mode without encryption or client verification")
	flagSet.StringVar(&inputs.Host, hostFlag, "", "Proxy hostname (Default all interfaces)")
	flagSet.IntVar(&inputs.Port, portFlag, defaultPort, "Proxy port")
	flagSet.BoolVar(&inputs.Quiet, quietFlag, false, "No output")
	flagSet.BoolVar(&inputs.Verbose, verboseFlag, false, "Verbose output")
	if showHelpIfNeeded(args[1:], flagSet) {
		return nil, nil
	}

	err := flagSet.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	// Check for extra positional arguments
	if flagSet.NArg() > 0 {
		return nil, fmt.Errorf("unexpected arguments: %v", flagSet.Args())
	}

	// Ensure valid port range
	if inputs.Port < 0 || inputs.Port > 65535 {
		return nil, fmt.Errorf("invalid port value: %d. Must be in the range 0 to 65535", inputs.Port)
	}

	// Check for required input arguments
	if inputs.CertificateFile == "" && !inputs.DisableMututalTLS {
		return nil, fmt.Errorf("must either provide a certificate json file (-%s) or specify insecure mode (-%s)", certificateFlag, disableMutualTLSFlag)
	}

	// Error on ambiguous combinations of flags
	if inputs.CertificateFile != "" && inputs.DisableMututalTLS {
		return nil, fmt.Errorf("cannot both specify a certificate json file (-%s) and specify insecure mode (-%s)", certificateFlag, disableMutualTLSFlag)
	}

	if inputs.Quiet && inputs.Verbose {
		return nil, fmt.Errorf("cannot set both -%s and -%s flags", quietFlag, verboseFlag)
	}

	return &inputs, nil
}

// Start SOCKS5 proxy server and print connect url for clients to use
func proxyConnections(socket net.Listener, logger *slog.Logger, args *proxyInputs) error {
	socksPrefix := socks5proxy.Socks5SecureUrlPrefix
	if args.DisableMututalTLS {
		socksPrefix = socks5proxy.Socks5UrlPrefix
	}
	displayedProxyUrl := server.ConstructConnectURL(socksPrefix, args.Host, args.Port)
	if !args.Quiet {
		fmt.Printf("SOCKS5 proxy ready to accept connections at: %s\n", displayedProxyUrl)
	}
	err := socks5proxy.StartSOCKS5Proxy(socket, logger)
	return err
}

var helpFlags = []string{"-h", "-help", "--help"}

func isHelpFlag(arg string) bool {
	for _, h := range helpFlags {
		if arg == h {
			return true
		}
	}
	return false
}

// If any command-line argument is a help flag, show help text
func showHelpIfNeeded(args []string, flags *flag.FlagSet) bool {
	for _, arg := range args {
		if isHelpFlag(arg) {
			printHelp(flags)
			return true
		}
	}
	return false
}

var longHelpText = `Start a proxy server to proxy traffic between MATLAB clients and a MATLAB Parallel Server cluster.
This can be used to create a single access point for the cluster on the host where the proxy server is run and via the port specified.

The proxy server uses the SOCKS5 protocol to allow connecting clients to specify the destination to connect to. Authentication of 
connecting clients is provided by mutual TLS (mTLS) using client certificates. Certificate files can be generated using the mjssetup tool.`

// Print help text
func printHelp(flags *flag.FlagSet) {
	fmt.Printf("%s\n\n", longHelpText)
	fmt.Printf("%s", getAllUsageText())
	fmt.Printf("Input arguments:\n")
	flags.PrintDefaults()
}

var usageExamples = []struct {
	description string
	command     string
}{
	{"Start a proxy server using mutual TLS",
		fmt.Sprintf("%s -%s <certificate-file> [-%s <proxy-host>] [-%s <proxy-port>]", programName, certificateFlag, hostFlag, portFlag)},
	{"Start an insecure proxy server using plain TCP without encryption or authentication of client connections",
		fmt.Sprintf("%s -%s [-%s <proxy-host>] [-%s <proxy-port>]", programName, disableMutualTLSFlag, hostFlag, portFlag)},
	{"Start a proxy server with no command window output",
		fmt.Sprintf("%s -%s <certificate-file> -%s", programName, certificateFlag, quietFlag)},
	{"Start a proxy server with verbose command window output",
		fmt.Sprintf("%s -%s <certificate-file> -%s", programName, certificateFlag, verboseFlag)},
}

// Create a string containing usage text
func getAllUsageText() string {
	txt := fmt.Sprintf("Usage: %s [<args>]", programName)
	txt = txt + "\n\n"
	for _, usageExample := range usageExamples {
		txt = txt + getUsageText(usageExample.description, usageExample.command) + "\n"
	}
	return txt
}

func getUsageText(description string, command string) string {
	return fmt.Sprintf("%s:\n%s\n", description, command)
}
