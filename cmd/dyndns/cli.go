package main

import (
	"context"
	"flag"
	"os"

	"github.com/ViBiOh/dyndns/pkg/dyndns"
	"github.com/ViBiOh/dyndns/pkg/ip"
	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func main() {
	fs := flag.NewFlagSet("dyndns", flag.ExitOnError)

	url := flags.String(fs, "", "ip", "URL", "URL for getting IPv4 or v6", "https://api64.ipify.org", nil)
	network := flags.String(fs, "", "ip", "Network", "Network", "tcp4", nil)
	loggerConfig := logger.Flags(fs, "logger")
	dyndnsConfig := dyndns.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	ctx := context.Background()

	logger.Global(logger.New(loggerConfig))

	currentIP, err := ip.Get(ctx, *url, *network)
	logger.Fatal(err)

	logger.Info("Current IP is: %s", currentIP)

	dyndnsApp, err := dyndns.New(dyndnsConfig)
	logger.Fatal(err)

	logger.Fatal(dyndnsApp.Do(ctx, currentIP))
}
