package dyndns

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
)

type Service struct {
	api     *cloudflare.API
	entry   string
	domains []string
	proxied bool
}

type Config struct {
	Token   string
	Entry   string
	Domains []string
	Proxied bool
}

func Flags(fs *flag.FlagSet, prefix string) Config {
	var config Config

	flags.New("Token", "Cloudflare token").Prefix(prefix).DocPrefix("dyndns").StringVar(fs, &config.Token, "", nil)
	flags.New("Domain", "Domain to configure").Prefix(prefix).DocPrefix("dyndns").StringSliceVar(fs, &config.Domains, nil, nil)
	flags.New("Entry", "DNS Entry CNAME").Prefix(prefix).DocPrefix("dyndns").StringVar(fs, &config.Entry, "dyndns", nil)
	flags.New("Proxied", "Proxied").Prefix(prefix).DocPrefix("dyndns").BoolVar(fs, &config.Proxied, false, nil)

	return config
}

func New(config Config) (Service, error) {
	api, err := cloudflare.NewWithAPIToken(config.Token)
	if err != nil {
		return Service{}, fmt.Errorf("create API client: %w", err)
	}

	return Service{
		domains: config.Domains,
		entry:   config.Entry,
		proxied: config.Proxied,

		api: api,
	}, nil
}

func (s Service) Do(ctx context.Context, ip string) error {
	for _, domain := range s.domains {
		if err := s.do(ctx, ip, domain); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) do(ctx context.Context, ip, domain string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	zoneID, err := s.api.ZoneIDByName(domain)
	if err != nil {
		return fmt.Errorf("get zone by name: %w", err)
	}

	dnsType := "A"
	if len(ip) == net.IPv6len {
		dnsType = "AAAA"
	}

	dnsName := fmt.Sprintf("%s.%s", s.entry, domain)

	return s.upsertEntry(ctx, cloudflare.ZoneIdentifier(zoneID), dnsType, dnsName, ip)
}

func (s Service) upsertEntry(ctx context.Context, zoneIdentifier *cloudflare.ResourceContainer, dnsType, dnsName, content string) error {
	records, results, err := s.api.ListDNSRecords(ctx, zoneIdentifier, cloudflare.ListDNSRecordsParams{
		Type: dnsType,
		Name: dnsName,
	})
	if err != nil {
		return fmt.Errorf("list dns records: %w", err)
	}

	if results.Count == 0 {
		slog.Info("Creating record", "type", dnsType, "name", dnsName, "content", content)
		_, err := s.api.CreateDNSRecord(ctx, zoneIdentifier, cloudflare.CreateDNSRecordParams{
			Type:    dnsType,
			Name:    dnsName,
			Content: content,
			Proxied: &s.proxied,
		})
		if err != nil {
			return fmt.Errorf("create dns record: %w", err)
		}

		return nil
	}

	slog.Info("Updating record", "type", dnsType, "name", dnsName, "content", content)
	_, err = s.api.UpdateDNSRecord(ctx, zoneIdentifier, cloudflare.UpdateDNSRecordParams{
		ID:      records[0].ID,
		Type:    dnsType,
		Name:    dnsName,
		Content: content,
		Proxied: &s.proxied,
	})

	if err != nil {
		return fmt.Errorf("update dns record: %w", err)
	}

	return nil
}
