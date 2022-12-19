package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	httpc "github.com/milgradesec/go-libs/http"

	"github.com/rs/zerolog/log"

	"github.com/milgradesec/ddns/internal/config"
	"github.com/milgradesec/ddns/pkg/ip"
)

// CloudflareDNS implements provider.DNSProvider interface.
type CloudflareDNS struct {
	config *config.Configuration
	api    *cloudflare.API
	zoneID string
}

// New creates a Cloudflare DNS provider.
func New(config *config.Configuration) (*CloudflareDNS, error) {
	var (
		email    string
		apiKey   string
		apiToken string
		found    bool
	)

	// Authenticate using an API Token
	tokenFile, found := os.LookupEnv("CLOUDFLARE_API_TOKEN_FILE")
	if found {
		buf, err := os.ReadFile(tokenFile)
		if err != nil {
			return nil, err
		}
		apiToken = strings.TrimSpace(string(buf))

		return newWithAPIToken(apiToken, config)
	}

	apiToken, found = os.LookupEnv("CLOUDFLARE_API_TOKEN")
	if found {
		log.Debug().
			Str("name", "CLOUDFLARE_API_TOKEN").
			Str("value", "***REDACTED***").
			Msg("Found environment variable")
		return newWithAPIToken(apiToken, config)
	}

	// Authenticate using Email + API Key
	email, found = os.LookupEnv("CLOUDFLARE_EMAIL")
	if !found {
		return nil, errors.New("no Cloudflare API credentials found")
	}

	keyFile, found := os.LookupEnv("CLOUDFLARE_API_KEY_FILE")
	if found {
		buf, err := os.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		apiKey = strings.TrimSpace(string(buf))

		return newWithAPIKey(apiKey, email, config)
	}

	apiKey, found = os.LookupEnv("CLOUDFLARE_API_KEY")
	if found {
		log.Debug().
			Str("name", "CLOUDFLARE_API_KEY").
			Str("value", "***REDACTED***").
			Msg("Found environment variable")
		return newWithAPIKey(apiKey, email, config)
	}

	return nil, errors.New("unable to find Cloudflare API credentials")
}

// Creates a Cloudflare Provider using API Token.
func newWithAPIToken(apiToken string, config *config.Configuration) (*CloudflareDNS, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken, cloudflare.HTTPClient(httpc.NewHTTPClient()))
	if err != nil {
		return nil, err
	}

	cf := &CloudflareDNS{
		api:    api,
		config: config,
	}
	return cf, nil
}

// Creates a Clouflare Provider using API Key.
func newWithAPIKey(apiKey, email string, config *config.Configuration) (*CloudflareDNS, error) {
	api, err := cloudflare.New(apiKey, email, cloudflare.HTTPClient(httpc.NewHTTPClient()))
	if err != nil {
		return nil, err
	}

	cf := &CloudflareDNS{
		api:    api,
		config: config,
	}
	return cf, nil
}

// Name implements the provider.DNSProvider interface.
func (cf *CloudflareDNS) Name() string {
	return "Cloudflare"
}

// GetZoneName implements the provider.DNSProvider interface.
func (cf *CloudflareDNS) GetZoneName() string {
	return cf.config.Zone
}

// UpdateZone implements provider.DNSProvider interface.
func (cf *CloudflareDNS) UpdateZone(ctx context.Context) error {
	if cf.zoneID == "" {
		log.Debug().
			Str("provider", cf.Name()).
			Str("zone", cf.GetZoneName()).
			Msg("Using Cloudflare API to retrieve Zone ID")

		id, err := cf.api.ZoneIDByName(cf.GetZoneName())
		if err != nil {
			return fmt.Errorf("cloudflare api error: failed to retrieve zone id: %w", err)
		}
		cf.zoneID = id
	}

	log.Debug().
		Str("provider", cf.Name()).
		Str("zone", cf.GetZoneName()).
		Msg("Using Cloudflare API to fetch DNS records")

	// Fetch only A type records
	a := cloudflare.DNSRecord{
		Type: "A",
	}
	records, err := cf.api.DNSRecords(ctx, cf.zoneID, a)
	if err != nil {
		return fmt.Errorf("cloudflare api error: failed to list dns records: %w", err)
	}
	if len(records) == 0 {
		return errors.New("no A records found")
	}

	publicIP, err := ip.GetIP(ctx)
	if err != nil {
		return err
	}
	log.Debug().
		Str("provider", ip.DefaultProviderName).
		Str("address", publicIP).
		Msg("Detected current public IP")

	for _, r := range records {
		log.Debug().
			Str("provider", cf.Name()).
			Str("zone", cf.GetZoneName()).
			Str("type", r.Type).
			Str("value", r.Content).
			Msgf("Found record '%s'", r.Name)

		if cf.config.IsExcluded(r.Name) {
			log.Info().
				Str("provider", cf.Name()).
				Str("zone", cf.GetZoneName()).
				Msgf("Changes for record '%s' are excluded by configuration", r.Name)
			continue
		}

		if r.Content == publicIP {
			log.Debug().
				Str("provider", cf.Name()).
				Str("zone", cf.GetZoneName()).
				Msgf("No changes needed for '%s'", r.Name)
			continue
		}

		log.Debug().
			Str("provider", cf.Name()).
			Str("zone", cf.GetZoneName()).
			Msgf("Record '%s' needs update", r.Name)
		err := cf.updateDNSRecord(ctx, publicIP, r)
		if err != nil {
			return err
		}

		log.Info().
			Str("provider", cf.Name()).
			Str("zone", cf.GetZoneName()).
			Msgf("Updated record '%s' from %s to %s", r.Name, r.Content, publicIP)
	}
	return nil
}

func (cf *CloudflareDNS) updateDNSRecord(ctx context.Context, value string, r cloudflare.DNSRecord) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Debug().
		Str("provider", cf.Name()).
		Str("zone", cf.GetZoneName()).
		Msgf("Using Cloudflare API to update DNS record '%s'", r.Name)

	rr := cloudflare.DNSRecord{
		Type:    r.Type,
		Name:    r.Name,
		Content: value,
		Proxied: r.Proxied,
	}

	if err := cf.api.UpdateDNSRecord(ctx, cf.zoneID, r.ID, rr); err != nil {
		return fmt.Errorf("error updating record %s: %w", r.Name, err)
	}
	return nil
}
