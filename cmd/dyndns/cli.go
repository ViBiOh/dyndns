package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/ViBiOh/dyndns/pkg/dyndns"
	"github.com/ViBiOh/dyndns/pkg/ip"
	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func main() {
	fs := flag.NewFlagSet("dyndns", flag.ExitOnError)
	fs.Usage = flags.Usage(fs)

	url := flags.New("URL", "URL for getting IPv4 or v6").DocPrefix("ip").String(fs, "https://api64.ipify.org", nil)
	network := flags.New("Network", "Network").DocPrefix("ip").String(fs, "tcp4", nil)
	loggerConfig := logger.Flags(fs, "logger")
	dyndnsConfig := dyndns.Flags(fs, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	logger.Init(loggerConfig)

	ctx := context.Background()

	currentIP, err := ip.Get(ctx, *url, *network)
	if err != nil {
		slog.ErrorContext(ctx, "get ip", "err", err)
		os.Exit(1)
	}

	slog.Info("Current IP", "ip", currentIP)

	dyndnsApp, err := dyndns.New(dyndnsConfig)
	if err != nil {
		slog.ErrorContext(ctx, "create dyndns", "err", err)
		os.Exit(1)
	}

	if err := dyndnsApp.Do(ctx, currentIP); err != nil {
		slog.ErrorContext(ctx, "execute dyndns", "err", err)
		os.Exit(1)
	}
}
