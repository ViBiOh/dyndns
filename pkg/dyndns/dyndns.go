package dyndns

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
)

type App struct {
	api     *cloudflare.API
	entry   string
	domains []string
	proxied bool
}

type Config struct {
	token   *string
	domains *[]string
	entry   *string
	proxied *bool
}

func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		token:   flags.New("Token", "Cloudflare token").Prefix(prefix).DocPrefix("dyndns").String(fs, "", nil),
		domains: flags.New("Domain", "Domain to configure").Prefix(prefix).DocPrefix("dyndns").StringSlice(fs, nil, nil),
		entry:   flags.New("Entry", "DNS Entry CNAME").Prefix(prefix).DocPrefix("dyndns").String(fs, "dyndns", nil),
		proxied: flags.New("Proxied", "Proxied").Prefix(prefix).DocPrefix("dyndns").Bool(fs, false, nil),
	}
}

func New(config Config) (App, error) {
	api, err := cloudflare.NewWithAPIToken(strings.TrimSpace(*config.token))
	if err != nil {
		return App{}, fmt.Errorf("create API client: %w", err)
	}

	return App{
		domains: *config.domains,
		entry:   strings.TrimSpace(*config.entry),
		proxied: *config.proxied,

		api: api,
	}, nil
}

func (a App) Do(ctx context.Context, ip string) error {
	for _, domain := range a.domains {
		if err := a.do(ctx, ip, domain); err != nil {
			return err
		}
	}

	return nil
}

func (a App) do(ctx context.Context, ip, domain string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	zoneID, err := a.api.ZoneIDByName(domain)
	if err != nil {
		return fmt.Errorf("get zone by name: %w", err)
	}

	dnsType := "A"
	if len(ip) == net.IPv6len {
		dnsType = "AAAA"
	}

	dnsName := fmt.Sprintf("%s.%s", a.entry, domain)

	return a.upsertEntry(ctx, cloudflare.ZoneIdentifier(zoneID), dnsType, dnsName, ip)
}

func (a App) upsertEntry(ctx context.Context, zoneIdentifier *cloudflare.ResourceContainer, dnsType, dnsName, content string) error {
	records, results, err := a.api.ListDNSRecords(ctx, zoneIdentifier, cloudflare.ListDNSRecordsParams{
		Type: dnsType,
		Name: dnsName,
	})
	if err != nil {
		return fmt.Errorf("list dns records: %w", err)
	}

	if results.Count == 0 {
		slog.Info("Creating record", "type", dnsType, "name", dnsName, "content", content)
		_, err := a.api.CreateDNSRecord(ctx, zoneIdentifier, cloudflare.CreateDNSRecordParams{
			Type:    dnsType,
			Name:    dnsName,
			Content: content,
			Proxied: &a.proxied,
		})
		if err != nil {
			return fmt.Errorf("create dns record: %w", err)
		}

		return nil
	}

	slog.Info("Updating record", "type", dnsType, "name", dnsName, "content", content)
	_, err = a.api.UpdateDNSRecord(ctx, zoneIdentifier, cloudflare.UpdateDNSRecordParams{
		ID:      records[0].ID,
		Type:    dnsType,
		Name:    dnsName,
		Content: content,
		Proxied: &a.proxied,
	})

	if err != nil {
		return fmt.Errorf("update dns record: %w", err)
	}

	return nil
}
