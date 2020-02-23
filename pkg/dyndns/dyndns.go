package dyndns

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/ViBiOh/httputils/v3/pkg/flags"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

// App of package
type App interface {
	Do(ip string) error
}

// Config of package
type Config struct {
	token   *string
	domain  *string
	entry   *string
	proxied *bool
}

type app struct {
	domain  string
	entry   string
	proxied bool

	api *cloudflare.API
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
		return nil, fmt.Errorf("unable to create API client: %s", err)
	}

	return &app{
		domain: strings.TrimSpace(*config.domain),
		entry:  strings.TrimSpace(*config.entry),

		api: api,
	}, nil
}

func (a app) Do(ip string) error {
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
	records, err := a.api.DNSRecords(zoneID, dnsRecord)
	if err != nil {
		return fmt.Errorf("unable to list dns records: %s", err)
	}

	dnsRecord.Content = ip
	dnsRecord.Proxied = a.proxied

	if len(records) == 0 {
		_, err := a.api.CreateDNSRecord(zoneID, dnsRecord)
		return fmt.Errorf("unable to create dns record: %s", err)
	}
	return a.api.UpdateDNSRecord(zoneID, records[0].ID, dnsRecord)
}
