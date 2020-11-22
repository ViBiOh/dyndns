package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/goweb/pkg/dyndns"
	"github.com/ViBiOh/goweb/pkg/ip"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

func main() {
	fs := flag.NewFlagSet("dyndns", flag.ExitOnError)

	network := flags.New("", "ip").Name("Network").Default("tcp4").Label("Network").ToString(fs)
	dyndnsConfig := dyndns.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	currentIP, err := ip.Get(*network)
	logger.Fatal(err)

	logger.Info("Current IP is: %s", currentIP)

	dyndnsApp, err := dyndns.New(dyndnsConfig)
	logger.Fatal(err)

	logger.Fatal(dyndnsApp.Do(currentIP))
}
