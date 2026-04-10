// Copyright 2024 The MathWorks, Inc.
package socks5proxy

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/bufferpool"
)

const Socks5UrlPrefix = "socks5"
const Socks5SecureUrlPrefix = "socks5s"
const BUFFER_POOL_KB_ENV = "PARALLEL_PROXY_BUFFER_POOL_KB"

func StartSOCKS5Proxy(socket net.Listener, logger *slog.Logger) error {
	// Create a SOCKS5 server
	ruleset := socks5.NewPermitConnAndAss()
	ruleset = NewLoggingRuleSet(ruleset, logger)

	// Size of buffer used to copy bytes for proxied connections
	bufferPoolKB := 32
	if bufferPoolEnv := os.Getenv(BUFFER_POOL_KB_ENV); bufferPoolEnv != "" {
		envBufferPoolKB, err := strconv.Atoi(bufferPoolEnv)
		if err == nil {
			bufferPoolKB = envBufferPoolKB
		}
	}
	logger.Debug(fmt.Sprintf("SOCKS5 proxy buffer pool = %d KB", bufferPoolKB))

	socksLogger := NewSocks5Logger(logger)
	server := socks5.NewServer(
		socks5.WithLogger(socksLogger),
		socks5.WithRule(ruleset),
		socks5.WithBufferPool(bufferpool.NewPool(bufferPoolKB*1024)),
	)

	err := server.Serve(socket)
	return err
}
