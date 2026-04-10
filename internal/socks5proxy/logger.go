// Copyright 2024 The MathWorks, Inc.
package socks5proxy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
)

// socks5 logger adaptor for slog logger
type Socks5Logger struct {
	Logger *slog.Logger
}

func NewSocks5Logger(logger *slog.Logger) socks5.Logger {
	return &Socks5Logger{Logger: logger}
}

func (sl Socks5Logger) Errorf(format string, args ...interface{}) {
	sl.Logger.Error(fmt.Sprintf(format, args...))
}

// socks5 ruleset which logs all commands
type LoggingRuleSet struct {
	wrappedRuleSet socks5.RuleSet
	logger         *slog.Logger
}

func NewLoggingRuleSet(ruleset socks5.RuleSet, logger *slog.Logger) socks5.RuleSet {
	return &LoggingRuleSet{wrappedRuleSet: ruleset, logger: logger}
}

func (p *LoggingRuleSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	switch req.Command {
	case statute.CommandConnect:
		p.logger.Debug("CONNECT request for " + req.DestAddr.String())
	case statute.CommandBind:
		p.logger.Debug("BIND request for " + req.DestAddr.String())
	case statute.CommandAssociate:
		p.logger.Debug("ASSOCIATE request for " + req.DestAddr.String())
	}
	return (p.wrappedRuleSet).Allow(ctx, req)
}
