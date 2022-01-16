package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/dyndns/pkg/dyndns"
	"github.com/ViBiOh/dyndns/pkg/ip"
	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func main() {
	fs := flag.NewFlagSet("dyndns", flag.ExitOnError)

	url := flags.New("", "ip", "URL").Default("https://api64.ipify.org", nil).Label("URL for getting IPv4 or v6").ToString(fs)
	network := flags.New("", "ip", "Network").Default("tcp4", nil).Label("Network").ToString(fs)
	dyndnsConfig := dyndns.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	currentIP, err := ip.Get(*url, *network)
	logger.Fatal(err)

	logger.Info("Current IP is: %s", currentIP)

	dyndnsApp, err := dyndns.New(dyndnsConfig)
	logger.Fatal(err)

	logger.Fatal(dyndnsApp.Do(currentIP))
}
