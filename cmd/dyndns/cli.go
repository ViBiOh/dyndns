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

	url := flags.New("URL", "URL for getting IPv4 or v6").DocPrefix("ip").String(fs, "https://ifconfig.co/ip", nil)
	network := flags.New("Network", "Network").DocPrefix("ip").String(fs, "tcp4", nil)
	loggerConfig := logger.Flags(fs, "logger")
	dyndnsConfig := dyndns.Flags(fs, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	logger.Init(ctx, loggerConfig)

	currentIP, err := ip.Get(ctx, *url, *network)
	logger.FatalfOnErr(ctx, err, "get ip")

	slog.LogAttrs(ctx, slog.LevelInfo, "Current IP", slog.String("ip", currentIP))

	dyndnsApp, err := dyndns.New(dyndnsConfig)
	logger.FatalfOnErr(ctx, err, "create dyndns")

	if err := dyndnsApp.Do(ctx, currentIP); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "execute dyndns", slog.Any("error", err))
		os.Exit(1)
	}
}
