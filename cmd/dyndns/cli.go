package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/goweb/pkg/dyndns"
	"github.com/ViBiOh/goweb/pkg/ip"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

func main() {
	fs := flag.NewFlagSet("dyndns", flag.ExitOnError)

	dyndnsConfig := dyndns.Flags(fs, "")

	logger.Fatal(fs.Parse(os.Args[1:]))

	dyndnsApp, err := dyndns.New(dyndnsConfig)
	logger.Fatal(err)

	ip, err := ip.Get()
	logger.Fatal(err)

	logger.Info("Current IP is : %s", ip)

	logger.Fatal(dyndnsApp.Do(ip))
}
