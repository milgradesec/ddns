package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
	Zone string

	zoneID string
	api    *cloudflare.API
	config *config.Configuration
}

// New creates a Cloudflare DNS provider.
func New(config *config.Configuration) (*CloudflareDNS, error) { //nolint
	var (
		zone     string
		email    string
		apiKey   string
		apiToken string
		found    bool
	)

	zone, found = os.LookupEnv("CLOUDFLARE_ZONE")
	if !found {
		return nil, errors.New("CLOUDFLARE_ZONE not set")
	}
	config.Zone = zone
	log.Info().Msgf("CLOUDFLARE_ZONE => %s", zone)

	// Authenticate using an API Token
	tokenFile, found := os.LookupEnv("CLOUDFLARE_API_TOKEN_FILE")
	if found {
		log.Info().Msg("CLOUDFLARE_API_TOKEN_FILE found")

		buf, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			return nil, err
		}
		apiToken = strings.TrimSpace(string(buf))

		return newWithAPIToken(apiToken, config)
	}

	apiToken, found = os.LookupEnv("CLOUDFLARE_API_TOKEN")
	if found {
		log.Info().Msg("CLOUDFLARE_API_TOKEN found")
		return newWithAPIToken(apiToken, config)
	}

	// Authenticate using Email + API Key
	email, found = os.LookupEnv("CLOUDFLARE_EMAIL")
	if !found {
		return nil, errors.New("no Cloudflare API credentials found")
	}
	log.Info().Msg("CLOUDFLARE_EMAIL found")

	keyFile, found := os.LookupEnv("CLOUDFLARE_API_KEY_FILE")
	if found {
		log.Info().Msg("CLOUDFLARE_API_KEY_FILE found")

		buf, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		apiKey = strings.TrimSpace(string(buf))

		return newWithAPIKey(apiKey, email, config)
	}

	apiKey, found = os.LookupEnv("CLOUDFLARE_API_KEY")
	if found {
		log.Info().Msg("CLOUDFLARE_API_KEY found")
		return newWithAPIKey(apiKey, email, config)
	}

	return nil, errors.New("unable to find Cloudflare API credentials")
}

// func newWithAPIToken(token string) (*cloudflare.API, error) {
// 	return cloudflare.NewWithAPIToken(token, cloudflare.HTTPClient(httpc.NewHTTPClient()))
// }

// func newWithAPIKey(key string, email string) (*cloudflare.API, error) {
// 	return cloudflare.New(key, email, cloudflare.HTTPClient(httpc.NewHTTPClient()))
// }

// Name implements the provider.DNSProvider interface.
func (cf *CloudflareDNS) Name() string {
	return "Cloudflare"
}

// GetZoneName implements the provider.DNSProvider interface.
func (cf *CloudflareDNS) GetZoneName() string {
	return cf.Zone
}

// Creates a Cloudflare Provider using API Token.
func newWithAPIToken(apiToken string, config *config.Configuration) (*CloudflareDNS, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken, cloudflare.HTTPClient(httpc.NewHTTPClient()))
	if err != nil {
		return nil, err
	}

	cf := &CloudflareDNS{
		Zone:   config.Zone,
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
		Zone:   config.Zone,
		api:    api,
		config: config,
	}
	return cf, nil
}

// UpdateZone implements provider.DNSProvider interface.
func (cf *CloudflareDNS) UpdateZone(ctx context.Context) error {
	if cf.zoneID == "" {
		id, err := cf.api.ZoneIDByName(cf.Zone)
		if err != nil {
			return fmt.Errorf("cloudflare api error: failed to retrieve zone id: %w", err)
		}
		cf.zoneID = id
	}

	publicIP, err := ip.GetIP(ctx)
	if err != nil {
		return err
	}

	records, err := cf.api.DNSRecords(ctx, cf.zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return fmt.Errorf("cloudflare api error: failed to list dns records: %w", err)
	}

	for i := range records {
		r := records[i]

		if cf.config.IsExcluded(r.Name) {
			continue
		}

		if r.Type == "A" {
			if r.Content != publicIP {
				rr := cloudflare.DNSRecord{
					Type:    r.Type,
					Name:    r.Name,
					Content: publicIP,
					Proxied: r.Proxied,
				}

				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()

				if err := cf.api.UpdateDNSRecord(ctx, cf.zoneID, r.ID, rr); err != nil {
					return fmt.Errorf("error updating %s: %w", r.Name, err)
				}
				log.Info().Msgf("updated %s from %s to %s", r.Name, r.Content, publicIP)
			}
		}
	}
	return nil
}
