package dyndns

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

type Service struct {
	api     *cloudflare.Client
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

func Flags(fs *flag.FlagSet, prefix string) *Config {
	var config Config

	flags.New("Token", "Cloudflare API Token").Prefix(prefix).DocPrefix("dyndns").StringVar(fs, &config.Token, "", nil)
	flags.New("Domain", "Domain to configure").Prefix(prefix).DocPrefix("dyndns").StringSliceVar(fs, &config.Domains, nil, nil)
	flags.New("Entry", "DNS Entry CNAME").Prefix(prefix).DocPrefix("dyndns").StringVar(fs, &config.Entry, "dyndns", nil)
	flags.New("Proxied", "Proxied").Prefix(prefix).DocPrefix("dyndns").BoolVar(fs, &config.Proxied, false, nil)

	return &config
}

func New(config *Config) (Service, error) {
	api := cloudflare.NewClient(option.WithAPIToken(config.Token))

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

	zones, err := s.api.Zones.List(ctx, zones.ZoneListParams{Name: cloudflare.F(domain)})
	if err != nil {
		return fmt.Errorf("list zones: %w", err)
	}

	if len(zones.Result) == 0 {
		return errors.New("domain not found")
	}

	dnsType := dns.RecordListParamsTypeA
	if len(ip) == net.IPv6len {
		dnsType = dns.RecordListParamsTypeAAAA
	}

	dnsName := fmt.Sprintf("%s.%s", s.entry, domain)

	return s.upsertEntry(ctx, zones.Result[0].ID, dnsName, ip, dnsType)
}

func (s Service) upsertEntry(ctx context.Context, zoneID, dnsName, content string, dnsType dns.RecordListParamsType) error {
	records, err := s.api.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(dnsName),
		}),
		Type: cloudflare.F(dnsType),
	})
	if err != nil {
		return fmt.Errorf("list dns records: %w", err)
	}

	if len(records.Result) == 0 {
		slog.LogAttrs(ctx, slog.LevelInfo, "Creating record", slog.String("type", string(dnsType)), slog.String("name", dnsName), slog.String("content", content))
		_, err := s.api.DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(zoneID),
			Body:   s.getNewBody(dnsType, dnsName, content),
		})
		if err != nil {
			return fmt.Errorf("create dns record: %w", err)
		}

		return nil
	}

	records.Result[0].Content = content

	slog.LogAttrs(ctx, slog.LevelInfo, "Updating record", slog.String("type", string(dnsType)), slog.String("name", dnsName), slog.String("content", content))
	_, err = s.api.DNS.Records.Update(ctx, records.Result[0].ID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(zoneID),
		Body:   s.getUpdateBody(dnsType, dnsName, content),
	})
	if err != nil {
		return fmt.Errorf("update dns record: %w", err)
	}

	return nil
}

func (s Service) getNewBody(dnsType dns.RecordListParamsType, dnsName, content string) dns.RecordNewParamsBodyUnion {
	switch dnsType {
	case dns.RecordListParamsTypeA:
		return dns.ARecordParam{
			Name:    cloudflare.F(dnsName),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(content),
			Proxied: cloudflare.F(s.proxied),
		}

	case dns.RecordListParamsTypeAAAA:
		return dns.AAAARecordParam{
			Name:    cloudflare.F(dnsName),
			Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			Content: cloudflare.F(content),
			Proxied: cloudflare.F(s.proxied),
		}

	default:
		return nil
	}
}

func (s Service) getUpdateBody(dnsType dns.RecordListParamsType, dnsName, content string) dns.RecordUpdateParamsBodyUnion {
	switch dnsType {
	case dns.RecordListParamsTypeA:
		return dns.ARecordParam{
			Name:    cloudflare.F(dnsName),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(content),
			Proxied: cloudflare.F(s.proxied),
		}

	case dns.RecordListParamsTypeAAAA:
		return dns.AAAARecordParam{
			Name:    cloudflare.F(dnsName),
			Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			Content: cloudflare.F(content),
			Proxied: cloudflare.F(s.proxied),
		}

	default:
		return nil
	}
}
