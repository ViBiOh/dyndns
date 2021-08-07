package dyndns

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
)

// App of package
type App struct {
	api *cloudflare.API

	domain string
	entry  string

	proxied bool
}

// Config of package
type Config struct {
	token   *string
	domain  *string
	entry   *string
	proxied *bool
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		token:   flags.New(prefix, "dyndns").Name("Token").Default("").Label("Cloudflare token").ToString(fs),
		domain:  flags.New(prefix, "dyndns").Name("Domain").Default("").Label("Domain to configure").ToString(fs),
		entry:   flags.New(prefix, "dyndns").Name("Entry").Default("dyndns").Label("DNS Entry CNAME").ToString(fs),
		proxied: flags.New(prefix, "dyndns").Name("Proxied").Default(false).Label("Proxied").ToBool(fs),
	}
}

// New creates new App from Config
func New(config Config) (App, error) {
	api, err := cloudflare.NewWithAPIToken(strings.TrimSpace(*config.token))
	if err != nil {
		return App{}, fmt.Errorf("unable to create API client: %s", err)
	}

	return App{
		domain:  strings.TrimSpace(*config.domain),
		entry:   strings.TrimSpace(*config.entry),
		proxied: *config.proxied,

		api: api,
	}, nil
}

// Do update dyndns on cloudflare
func (a App) Do(ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	zoneID, err := a.api.ZoneIDByName(a.domain)
	if err != nil {
		return fmt.Errorf("unable to found zone by name: %s", err)
	}

	dnsType := "A"
	if len(ip) == net.IPv6len {
		dnsType = "AAAA"
	}

	dnsRecord := cloudflare.DNSRecord{
		Type: dnsType,
		Name: fmt.Sprintf("%s.%s", a.entry, a.domain),
	}
	records, err := a.api.DNSRecords(ctx, zoneID, dnsRecord)
	if err != nil {
		return fmt.Errorf("unable to list dns records: %s", err)
	}

	dnsRecord.Content = ip
	dnsRecord.Proxied = &a.proxied

	if len(records) == 0 {
		logger.Info("Creating %s %s -> %s record", dnsRecord.Type, dnsRecord.Name, dnsRecord.Content)
		_, err := a.api.CreateDNSRecord(ctx, zoneID, dnsRecord)
		if err != nil {
			return fmt.Errorf("unable to create dns record: %s", err)
		}

		return nil
	}

	logger.Info("Updating %s %s -> %s record", dnsRecord.Type, dnsRecord.Name, dnsRecord.Content)
	return a.api.UpdateDNSRecord(ctx, zoneID, records[0].ID, dnsRecord)
}
